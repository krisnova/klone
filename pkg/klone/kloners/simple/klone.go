package simple

import (
	"github.com/kris-nova/klone/pkg/klone/kloners"
	"github.com/kris-nova/klone/pkg/kloneprovider"
)

type Kloner struct {
	gitServer kloneprovider.GitServer
}

func (k *Kloner) Clone(repo kloneprovider.Repo) (string, error) {
	return "", nil
}
func (k *Kloner) AddRemote(name string, remote kloneprovider.Repo) error {
	return nil
}
func (k *Kloner) Fork(parent kloneprovider.Repo) error {
	return nil
}
func (k *Kloner) Init() error {
	return nil
}

func NewKloner(srv kloneprovider.GitServer) (kloners.Kloner) {
	return &Kloner{
		gitServer: srv,
	}
}
