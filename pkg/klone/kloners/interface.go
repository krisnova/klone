package kloners

import "github.com/kris-nova/klone/pkg/kloneprovider"

type Kloner interface {
	Clone(repo kloneprovider.Repo) (string, error)
	Pull(remote string) error
	AddRemote(name, url string) error
	DeleteRemote(name string) error
	GetCloneDirectory(repo kloneprovider.Repo) string
}
