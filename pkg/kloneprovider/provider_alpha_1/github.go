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

package provider_alpha_1

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
	repos    []kloneprovider.Repo
	usr      *github.User
}

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
		s.username = username
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
	name := s.usr.Name
	local.Printf("%s successfully connected!", *name)
	s.ensureLocalAuthToken(token)
	return nil
}

func (s *GitServer) GetRepos() ([]kloneprovider.Repo, error) {
	var providerRepos []kloneprovider.Repo
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
				providerRepos = append(providerRepos, r)
			}
			if resp.NextPage == 0 {
				break
			}
			opt.ListOptions.Page = resp.NextPage
		}
		s.repos = providerRepos
	}
	return s.repos, nil
}

// ensureLocalAuthToken is the cache for GitHub tokens. This function
// should handle caching a working token (if necessary) to use in the future.
func (s *GitServer) ensureLocalAuthToken(token string) {
	cache := fmt.Sprintf("%s/.klone/auth", local.Home())
	authStr := local.SGetContent(cache)
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
		err = local.SPutContent(str, cache)
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
		err := local.SPutContent(token, cache)
		if err != nil {
			local.RecoverableErrorf("Unable to ensure local auth token: %v", err)
			return
		}
		local.Printf("Successfully cached access token!")
	}
}

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
	cache := fmt.Sprintf("%s/.klone/auth", local.Home())
	r := bufio.NewReader(os.Stdin)

	// Ensure our cache directory before looking for creds

	//Token
	var token string

	token = local.SGetContent(cache)
	if token == "" {
		os.MkdirAll(fmt.Sprintf("%s/.klone", local.Home()), 0700)
		local.SPutContent("", cache)
	}
	t := os.Getenv("KLONE_GITHUBTOKEN")
	if t != "" {
		// Logic here is to always overwrite with env vars in case a user changes
		token = t
	}
	creds.Token = token

	//User
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

	//Pass
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
