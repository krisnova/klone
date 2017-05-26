package kloners

import "github.com/kris-nova/klone/pkg/kloneprovider"

type Kloner interface {
	Clone(repo kloneprovider.Repo) (string, error)
	Pull(name string, remote kloneprovider.Repo) error
	AddRemote(name, url string, base kloneprovider.Repo) error
	DeleteRemote(name string, repo kloneprovider.Repo) error
	GetCloneDirectory(repo kloneprovider.Repo) string
}
