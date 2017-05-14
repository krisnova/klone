package kloners

import "github.com/kris-nova/klone/pkg/kloneprovider"

type Kloner interface {
	Clone(repo kloneprovider.Repo) (string, error)
	AddRemote(name string, remote kloneprovider.Repo) error
}
