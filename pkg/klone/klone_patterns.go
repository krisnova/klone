package klone

import (
	"github.com/kris-nova/klone/pkg/local"
	"strings"
	"fmt"
)

const secondsToWaitForGithubClone = 20

// The user is the owner, and the repository is not a fork
func (k *Kloneable) kloneOwner() (string, error) {
	local.Printf("Attempting git clone")
	path, err := k.kloner.Clone(k.repo)
	if err != nil {
		return "", err
	}
	return path, nil
}

// The user is the owner, and the repository was forked from somewhere
func (k *Kloneable) kloneAlreadyForked() (string, error) {
	local.Printf("Attempting git clone")
	path, err := k.kloner.Clone(k.repo)
	if err != nil {
		return "", err
	}
	local.Printf("Register remote [upstream]")
	err = k.kloner.AddRemote("upstream", k.repo.ForkedFrom(), k.repo)
	if err != nil {
		return path, err
	}
	return path, nil
}

// The user is NOT the owner, and the repository is already forked
func (k *Kloneable) kloneTryingFork() (string, error) {
	return k.kloneNeedsFork()
}

// The user is NOT the owner, and the user does NOT have a fork already
func (k *Kloneable) kloneNeedsFork() (string, error) {
	local.Printf("Forking [%s/%s] to [%s/%s]", k.repo.Owner(), k.repo.Name(), k.gitServer.OwnerName(), k.repo.Name())
	// GitHub fork
	newRepo, err := k.gitServer.Fork(k.repo, k.gitServer.OwnerName())
	if err != nil {
		if strings.Contains(err.Error(), "job scheduled on GitHub side") {
			// Forking takes a while in GitHub so let's wait for it
			for i := 1; i <= secondsToWaitForGithubClone; i++ {
				repo, err := k.gitServer.GetRepo(k.repo.Name())
				if err == nil {
					local.Printf("Succesfully detected new repository [%s/%s]", repo.Owner(), repo.Name())
					break
				}
				if i == secondsToWaitForGithubClone {
					return "", fmt.Errorf("unable to detect forked repository after waiting %d seconds", secondsToWaitForGithubClone)
				}
			}

		} else {
			return "", err
		}
	}
	local.Printf("Attempting git clone")
	path, err := k.kloner.Clone(k.repo)
	if err != nil {
		return "", err
	}
	local.Printf("Register remote [upstream]")
	err = k.kloner.AddRemote("upstream", newRepo, k.repo)
	if err != nil {
		return path, err
	}
	return path, nil
}
