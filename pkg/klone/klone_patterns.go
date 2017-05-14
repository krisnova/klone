package klone

import "github.com/kris-nova/klone/pkg/local"

// The user is the owner, and the repository is not a fork
func (k *Kloneable) kloneOwner() error {
	local.Printf("Attempting git clone")
	_, err := k.kloner.Clone(k.repo)
	if err != nil {
		return err
	}
	return nil
}

// The user is the owner, and the repository was forked from somewhere
func (k *Kloneable) kloneAlreadyForked() error {
	local.Printf("Attempting git clone")
	_, err := k.kloner.Clone(k.repo)
	if err != nil {
		return err
	}
	local.Printf("Register remote [upstream]")
	err = k.kloner.AddRemote("upstream", k.repo.ForkedFrom())
	if err != nil {
		return err
	}
	return nil
}

// The user is NOT the owner, and the repository is already forked
func (k *Kloneable) kloneTryingFork() error {
	return k.kloneNeedsFork()
}

// The user is NOT the owner, and the user does NOT have a fork already
func (k *Kloneable) kloneNeedsFork() error {
	local.Printf("Forking [%s/%s] to [%s/%s]", k.repo.Owner(), k.repo.Name(), k.gitServer.OwnerName(), k.repo.Name())
	// GitHub fork
	newRepo, err := k.gitServer.Fork(k.repo, k.gitServer.OwnerName())
	if err != nil {
		return err
	}
	local.Printf("Attempting git clone")
	_, err = k.kloner.Clone(k.repo)
	if err != nil {
		return err
	}
	local.Printf("Register remote [upstream]")
	err = k.kloner.AddRemote("upstream", newRepo)
	if err != nil {
		return err
	}
	return nil
}
