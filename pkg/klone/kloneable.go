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
// kloneable.go represents a repository that has been reasoned about, and is ready to klone

package klone

import (
	"github.com/kris-nova/klone/pkg/kloneprovider"
	"github.com/kris-nova/klone/pkg/klone/kloners"
	"github.com/kris-nova/klone/pkg/klone/kloners/simple"
	"github.com/kris-nova/klone/pkg/klone/kloners/golang"
	"strings"
	"github.com/kris-nova/klone/pkg/local"
	"errors"
)

const (
	StyleOwner         Style = 1 // The user is the owner, and the repository is not a fork
	StyleAlreadyForked Style = 2 // The user is the owner, and the repository was forked from somewhere
	StyleNeedsFork     Style = 3 // The user is NOT the owner, and the user does NOT have a fork already
	StyleTryingFork    Style = 4 // The user is NOT the owner, and the repository is already forked
)

// NewKlonerFunc defines the type of function we expect for new kloners
type NewKlonerFunc func(server kloneprovider.GitServer) (kloners.Kloner)

// LanguageToKloner maps languages to kloners
// All language keys should be lower case, and they are cast as such before assertion
var LanguageToKloner = map[string]NewKlonerFunc{
	"":   simple.NewKloner, // Empty lang can use a simple kloner
	"go": golang.NewKloner, // Go gets a special kloner
}

// Kloneable is a data structure that holds all relevant data to klone a repository
type Kloneable struct {
	gitServer kloneprovider.GitServer
	repo      kloneprovider.Repo
	style     Style
	kloner    kloners.Kloner
}

// Klone is the only exported method, and is the only way to take action on a Kloneable data structure
func (k *Kloneable) Klone() error {
	k.findKloner() // First things first, we will need a kloner
	switch k.style {
	case StyleOwner:
		return k.kloneOwner()
	case StyleAlreadyForked:
		return k.kloneAlreadyForked()
	case StyleNeedsFork:
		return k.kloneNeedsFork()
	case StyleTryingFork:
		return k.kloneTryingFork()
	}
	return nil
}

// findKloner is the logic that selects a kloner to use on a repository.
// Todo (@kris-nova) let's support .Klonefile's!
func (k *Kloneable) findKloner() error {
	if k.gitServer == nil {
		return errors.New("nil getServer")
	}
	lowerlang := strings.ToLower(k.repo.Language())
	if newKlonerFunc, ok := LanguageToKloner[lowerlang]; ok {
		kloner := newKlonerFunc(k.gitServer)
		local.Printf("Found Kloner [%s]", k.repo.Language())
		k.kloner = kloner
	} else {
		local.Printf("Unsupported language [%s], using Kloner [simple]")
		k.kloner = simple.NewKloner(k.gitServer)
	}
	return nil
}
