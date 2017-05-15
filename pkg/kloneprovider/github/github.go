// Copyright © 2017 Kris Nova <kris@nivenly.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//  _  ___
// | |/ / | ___  _ __   ___
// | ' /| |/ _ \| '_ \ / _ \
// | . \| | (_) | | | |  __/
// |_|\_\_|\___/|_| |_|\___|
//
// server.go is a representation of GitHub.com as a git server

package github

import (
	"github.com/kris-nova/klone/pkg/kloneprovider"
	"context"
	//"golang.org/x/oauth2"
	"github.com/google/go-github/github"
	"os"
	"bufio"
	"github.com/kris-nova/klone/pkg/local"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
	"golang.org/x/oauth2"
	"strings"
	"fmt"
)

var Cache = fmt.Sprintf("%s/.klone/auth", local.Home())

const (
	AccessTokenNote = "Access token automatically managed my Klone. More information: https://github.com/kris-nova/klone."
)

// GitServer is a representation of GitHub.com, by design we never store credentials here in memory
// Todo (@kris-nova) encrypt these creds please!
type GitServer struct {
	username string
	client   *github.Client
	ctx      context.Context
	cursor   kloneprovider.Repo
	repos    map[string]kloneprovider.Repo
	usr      *github.User
}

func (s *GitServer) Fork(parent kloneprovider.Repo, newOwner string) (kloneprovider.Repo, error) {
	c := &github.RepositoryCreateForkOptions{}
	// Override c.Orginzation here if we ever need one!
	repo, _, err := s.client.Repositories.CreateFork(s.ctx, parent.Owner(), parent.Name(), c)
	if err != nil {
		return nil, fmt.Errorf("unable to fork repository [%s]: %v", parent.Name(), err)
	}
	r := &Repo{
		impl: repo,
	}
	return r, nil
}

// GetServerString returns a server string is the string we would want to use in things like $GOPATH
// In this case we know we are dealing with GitHub.com so we can safely return it.
func (s *GitServer) GetServerString() string {
	return "github.com"
}

func (s *GitServer) OwnerName() string {
	return *s.usr.Login
}
func (s *GitServer) OwnerEmail() string {
	return *s.usr.Email
}

// Authenticate will parse configuration with the following hierarchy.
// 1. Access token from local cache
// 2. Access token from env var
// 3. Username/Password from env var
// Authenticate will then attempt to log in (prompting for MFA if necessary)
// Authenticate will then attempt to ensure a unique access token created by klone for future access
// To ensure a new auth token, simply set the env var and klone will re-cache the new token
func (s *GitServer) Authenticate(credentials kloneprovider.GitServerCredentials) error {
	r := bufio.NewReader(os.Stdin)
	token := credentials.(*GitServerCredentials).Token
	s.ctx = context.Background()
	var client *github.Client
	var tp github.BasicAuthTransport
	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(s.ctx, ts)
		client = github.NewClient(tc)
	} else {
		username := credentials.(*GitServerCredentials).User
		password := credentials.(*GitServerCredentials).Pass
		local.Printf("Connecting to GitHub: [%s]", username)
		tp = github.BasicAuthTransport{
			Username: strings.TrimSpace(username),
			Password: strings.TrimSpace(password),
		}
		client = github.NewClient(tp.Client())
	}
	s.client = client
	user, _, err := client.Users.Get(s.ctx, "")
	s.usr = user
	// Check if we need 2 factor
	if _, ok := err.(*github.TwoFactorAuthError); ok {
		local.PrintPrompt("GitHub Two Factor Auth Code: ")
		mfa, _ := r.ReadString('\n')
		tp.OTP = strings.TrimSpace(mfa)
		user, _, err := client.Users.Get(s.ctx, "")
		if err != nil {
			return err
		}
		s.usr = user
	} else if err != nil {
		return err
	}
	name := *s.usr.Login
	s.username = name
	local.Printf("Successfully authenticated [%s]", name)
	s.ensureLocalAuthToken(token)
	return nil
}

// GetRepoByOwner is the most effecient way to look up a repository exactly by it's name and owner
func (s *GitServer) GetRepoByOwner(owner, name string) (kloneprovider.Repo, error) {
	r := &Repo{}
	repo, _, err := s.client.Repositories.Get(s.ctx, owner, name)
	if err != nil {
		local.Printf("Unable to find repo [%s/%s]", owner, name)
		return r, err
	}
	if repo == nil {
		return r, nil
	}
	r.impl = repo
	r.assumedOwner = owner
	if *repo.Fork && repo.Parent.Owner != nil {
		r.forkedFrom = &Repo{impl: repo.Parent}
	}
	return r, nil

}

func (s *GitServer) NewRepo(name, desc string) (kloneprovider.Repo, error) {
	t := true
	gitRepo := &github.Repository{}
	gitRepo.Name = &name
	gitRepo.Description = &desc
	gitRepo.AutoInit = &t
	repo, _, err := s.client.Repositories.Create(s.ctx, "", gitRepo)
	if err != nil {
		return nil, err
	}
	r := &Repo{}
	r.impl = repo

	return r, nil
}

func (s *GitServer) DeleteRepoByOwner(name, owner string) (bool, error) {
	_, err := s.client.Repositories.Delete(s.ctx, owner, name)
	if err != nil {
		return false, err
	}
	return true, nil

}

func (s *GitServer) DeleteRepo(name string) (bool, error) {
	_, err := s.client.Repositories.Delete(s.ctx, s.username, name)
	if err != nil {
		return false, err
	}
	return true, nil

}

// GetRepo is the most effecient way to look up a repository exactly by it's name and assumed owner (you)
func (s *GitServer) GetRepo(name string) (kloneprovider.Repo, error) {
	r := &Repo{}
	repo, _, err := s.client.Repositories.Get(s.ctx, s.username, name)
	if err != nil {
		local.Printf("Unable to find repo [%s/%s]", s.username, name)
		return r, err
	}
	if repo == nil {
		return r, nil
	}
	r.impl = repo
	if *repo.Fork {
		r.forkedFrom = &Repo{impl: repo.Parent}
	}
	return r, nil
}

// GetRepos will return (and cache) a hash map of repositories by name for some
// convenient O(n*log(n)) look up!
func (s *GitServer) GetRepos() (map[string]kloneprovider.Repo, error) {
	providerRepos := make(map[string]kloneprovider.Repo)
	if len(s.repos) == 0 {
		opt := &github.RepositoryListOptions{
			ListOptions: github.ListOptions{PerPage: 100},
		}
		for {
			repos, resp, err := s.client.Repositories.List(s.ctx, s.username, opt)
			if err != nil {
				return providerRepos, err
			}
			for _, repo := range repos {
				r := &Repo{impl: repo}
				providerRepos[*r.impl.Name] = r
				// Here is where we look up forked repo information (if we have it)
				if *repo.Fork {
					rr, _, err := s.client.Repositories.Get(s.ctx, s.username, *r.impl.Name)
					if err != nil {
						return providerRepos, err
					}
					parent := &Repo{impl: rr.Parent}
					r.forkedFrom = parent
				}
			}
			if resp.NextPage == 0 {
				break
			}
			opt.ListOptions.Page = resp.NextPage
		}
		s.repos = providerRepos
	}
	local.Printf("Cached %d repositories in memory", len(s.repos))
	return s.repos, nil
}

// ensureLocalAuthToken is the cache for GitHub tokens. This function
// should handle caching a working token (if necessary) to use in the future.
func (s *GitServer) ensureLocalAuthToken(token string) {
	authStr := local.SGetContent(Cache)
	if authStr == "" {
		// We have no cache
		req := github.AuthorizationRequest{}
		note := AccessTokenNote
		req.Note = &note
		var auth *github.Authorization
		auth, resp, err := s.client.Authorizations.Create(s.ctx, &req)
		if resp.StatusCode == 422 && strings.Contains(err.Error(), "Code:already_exists") {
			// It already exists let's delete and create a new one
			auths, _, err := s.client.Authorizations.List(s.ctx, nil)
			if err != nil {
				local.RecoverableErrorf("Unable to delete existing auth token [%s]: %v", *auth.ID, err)
				return
			}
			for _, a := range auths {
				n := *a.Note
				if n == AccessTokenNote {
					_, err := s.client.Authorizations.Delete(s.ctx, *a.ID)
					if err != nil {
						local.RecoverableErrorf("Unable to delete existing auth token [%s]: %v", *auth.ID, err)
						return
					}
					auth, _, err = s.client.Authorizations.Create(s.ctx, &req)
					if err != nil {
						local.RecoverableErrorf("Unable to delete existing auth token [%s]: %v", *auth.ID, err)
						return
					}
					break
				}
			}
		} else if err != nil {
			local.RecoverableErrorf("Unable to ensure local auth token: %v", err)
			return
		}
		str := *auth.Token
		err = local.SPutContent(str, Cache)
		if err != nil {
			local.RecoverableErrorf("Unable to ensure local auth token: %v", err)
			return
		}
		local.Printf("Successfully cached access token!")
	} else if authStr != token && s.usr != nil {
		// We have a cache but it conflicts.
		// Check if we have conflicting tokens, but were able to auth
		// If so we probably have a new token, so let's overwrite
		local.Printf("Overwriting existing token with new token")
		err := local.SPutContent(token, Cache)
		if err != nil {
			local.RecoverableErrorf("Unable to ensure local auth token: %v", err)
			return
		}
		local.Printf("Successfully cached access token!")
	}
}

// GitServerCredentials are how we log into GitHub.com
type GitServerCredentials struct {
	User  string
	Pass  string
	Token string
}

// GetCredentials will look for GitHub access credentials with the following parsing logic:
// 1. Environmental variables $KLONE_GITHUBUSER $KLONE_GITHUBPASS
// 2. Config file ~/.klone/githubcreds
func (s *GitServer) GetCredentials() (kloneprovider.GitServerCredentials, error) {
	var creds GitServerCredentials
	r := bufio.NewReader(os.Stdin)

	// ----------------------------------------------------------------------------------------
	// Token
	var token string
	token = local.SGetContent(Cache)
	if token == "" {
		os.MkdirAll(fmt.Sprintf("%s/.klone", local.Home()), 0700)
		local.SPutContent("", Cache)
	}
	t := os.Getenv("KLONE_GITHUBTOKEN")
	if t != "" {
		// Logic here is to always overwrite with env vars in case a user changes
		token = t
	}
	creds.Token = token

	// ----------------------------------------------------------------------------------------
	// User
	var user string
	user = os.Getenv("KLONE_GITHUBUSER")
	if user == "" {
		local.PrintPrompt("GitHub Username: ")
		u, err := r.ReadString('\n')
		if err != nil {
			return creds, err
		}
		user = u
	}

	// ----------------------------------------------------------------------------------------
	// Pass
	var pass string
	pass = os.Getenv("KLONE_GITHUBPASS")
	if pass == "" {
		local.PrintPrompt("GitHub Password: ")
		bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
		p := string(bytePassword)
		pass = p
	}

	creds.User = user
	creds.Pass = pass
	return &creds, nil
}
