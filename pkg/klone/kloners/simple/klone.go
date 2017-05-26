package simple

import (
	"fmt"
	"github.com/kris-nova/klone/pkg/klone/kloners"
	"github.com/kris-nova/klone/pkg/kloneprovider"
	"github.com/kris-nova/klone/pkg/local"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"strings"
)

type Kloner struct {
	gitServer kloneprovider.GitServer
}

func (k *Kloner) Clone(repo kloneprovider.Repo) (string, error) {
	o := &git.CloneOptions{
		URL:               repo.GitCloneUrl(),
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}
	o.Auth = &ssh.PublicKeysCallback{}

	path := k.GetCloneDirectory(repo)
	local.Printf("Cloning into [%s]", path)
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

func (k *Kloner) AddRemote(name, url string, base kloneprovider.Repo) error {
	path := k.GetCloneDirectory(base)
	grepo, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("unable to open repository: %v", err)
	}
	c := &config.RemoteConfig{
		Name: name,
		URL:  url,
	}

	local.Printf("Adding remote [%s][%s]", name, url)
	_, err = grepo.CreateRemote(c)
	if err != nil {
		if strings.Contains(err.Error(), "remote already exists") {
			local.Printf("Remote: %s", err.Error())
			return nil
		} else {
			return fmt.Errorf("unable create remote: %v", err)
		}
	}
	//local.Printf("Fetching remote [%s]", url)
	//f := &git.FetchOptions{
	//	RemoteName: remote.Name(),
	//}
	//
	//// This is required for the git@github.com origin pattern
	//f.Auth = &ssh.PublicKeysCallback{}
	//
	//err = r.Fetch(f)
	//if err != nil {
	//	if strings.Contains(err.Error(), "already up-to-date") {
	//		local.Printf("Fetch: %s", err.Error())
	//		return nil
	//	}
	//	return fmt.Errorf("unable to fetch remote: %v", err)
	//}
	return nil
}
func (k *Kloner) Pull(name string, remote kloneprovider.Repo) error {
	return nil
}

func (k *Kloner) GetCloneDirectory(repo kloneprovider.Repo) string {
	return fmt.Sprintf("%s/%s", local.Home(), repo.Name())
}

func NewKloner(srv kloneprovider.GitServer) kloners.Kloner {
	return &Kloner{
		gitServer: srv,
	}
}
