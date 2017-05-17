package golang

import (
	"github.com/kris-nova/klone/pkg/klone/kloners"
	"github.com/kris-nova/klone/pkg/kloneprovider"
	"gopkg.in/src-d/go-git.v4"
	"os"
	"github.com/kris-nova/klone/pkg/local"
	"fmt"
	"gopkg.in/src-d/go-git.v4/config"
	"strings"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

type Kloner struct {
	gitServer kloneprovider.GitServer
}

// This is the logic that defins a Clone() for a Go repository
// Of course we need to check out into $GOPATH
func (k *Kloner) Clone(repo kloneprovider.Repo) (string, error) {
	o := &git.CloneOptions{
		URL:               repo.GitCloneUrl(),
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}
	o.Auth = &ssh.PublicKeysCallback{}

	path := k.GetCloneDirectory(repo)
	local.Printf("Cloning into $GOPATH [%s]", path)
	r, err := git.PlainClone(path, false, o)
	if err != nil {
		if strings.Contains(err.Error(), "repository already exists") {
			local.Printf("Clone: %s", err.Error())
			return path, nil
		} else if strings.Contains(err.Error(), "unknown capability") {
			// Todo (@kris-nova) handle capability errors better https://github.com/kris-nova/klone/issues/5
			local.RecoverableErrorf("bypassing capability error: %v ", err)
		} else {
			return "", fmt.Errorf("unable to clone repository: %v", err)
		}
	}
	local.Printf("Checking out HEAD")
	ref, err := r.Head()
	if err != nil {
		return "", fmt.Errorf("unable to checkout HEAD: %v", err)
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return "", fmt.Errorf("unable to checkout latest commit: %v", err)
	}
	local.Printf("HEAD checked out HEAD at [%s]", commit.Hash)
	return path, nil
}

func (k *Kloner) DeleteRemote(name string, repo kloneprovider.Repo) error {
	path := k.GetCloneDirectory(repo)
	grepo, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("unable to open repository: %v", err)
	}
	err = grepo.DeleteRemote(name)
	if err != nil {
		return err
	}
	return nil
}

// Add remote will add a new remote, and fetch from the remote branch
func (k *Kloner) AddRemote(name string, remote kloneprovider.Repo, base kloneprovider.Repo) error {
	path := k.GetCloneDirectory(base)
	grepo, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("unable to open repository: %v", err)
	}
	c := &config.RemoteConfig{
		Name: name,
		URL:  remote.GitCloneUrl(),
	}
	local.Printf("Adding remote [%s][%s]", name, remote.GitCloneUrl())
	r, err := grepo.CreateRemote(c)
	if err != nil {
		if strings.Contains(err.Error(), "remote already exists") {
			local.Printf("Remote: %s", err.Error())
			return nil
		} else {
			return fmt.Errorf("unable create remote: %v", err)
		}
	}
	local.Printf("Fetching remote [%s]", remote.Name())
	f := &git.FetchOptions{
		RemoteName: remote.Name(),
	}

	// This is required for the git@github.com origin pattern
	f.Auth = &ssh.PublicKeysCallback{}

	err = r.Fetch(f)
	if err != nil {
		if strings.Contains(err.Error(), "already up-to-date") {
			local.Printf("Fetch: %s", err.Error())
			return nil
		}
		return fmt.Errorf("unable to fetch remote: %v", err)
	}
	return nil
}

func NewKloner(srv kloneprovider.GitServer) (kloners.Kloner) {
	return &Kloner{
		gitServer: srv,
	}
}

type customPathFunc func(repo kloneprovider.Repo) string

// forkedFromCustomPath will check if a repository was forked from
// a certain parent, if so use a custom path setting
var forkedFromCustomPath = map[string]customPathFunc{
	"kubernetes": repoToKubernetesPath,
}

// repoToCloneDirectory will take a repository and reason about
// where to check out the repository on your local filesystem
func (k *Kloner) GetCloneDirectory(repo kloneprovider.Repo) string {
	var path string
	// Default path
	path = fmt.Sprintf("%s/src/%s/%s/%s", Gopath(), k.gitServer.GetServerString(), repo.Owner(), repo.Name())

	// Check for custom path overrides
	if repo.ForkedFrom() != nil {
		for forkedFromOwner, customFunc := range forkedFromCustomPath {
			if repo.ForkedFrom().Owner() == forkedFromOwner {
				path = customFunc(repo)
				break
			}
		}
	}
	return path
}

// Logic for getting $GOPATH
func Gopath() string {
	epath := os.Getenv("GOPATH")
	if epath == "" {
		// It's now safe to assume $HOME/go
		// thanks to Dave Cheney and the folks
		// who work on the standard library
		// https://github.com/golang/go/issues/17262
		path := fmt.Sprintf("%s/go", local.Home())
		return path
	} else if strings.Contains(epath, ":") {
		espl := strings.Split(epath, ":")
		if len(espl) <= 1 {
			return epath
		}
		// Here we will take the first gopath defined and use that
		return espl[0]
	}
	return epath
}
