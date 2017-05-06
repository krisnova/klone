package golang

import (
	"github.com/kris-nova/klone/pkg/klone/kloners"
)

type Kloner struct {
}

type KlonerContext struct {
}

func NewKloner() (kloners.Kloner) {
	return &Kloner{}
}

func (k *Kloner) SetContext(context kloners.KlonerContext) {

}

func (k *Kloner) Klone() error {
	return nil
}
