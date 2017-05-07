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
// klone.go is the primary logic for a "klone" operation, here we have a function that
// ONLY accepts a string (a repository name) and should be able to reason about what
// needs to be done to `git clone` the repo onto your machine.

package klone

import (
	"github.com/kris-nova/klone/pkg/kloneprovider"
	"github.com/kris-nova/klone/pkg/klone/kloners/simple"
	"github.com/kris-nova/klone/pkg/klone/kloners"
	"github.com/kris-nova/klone/pkg/local"
	"fmt"
	"strings"
)

type Style int

const (
	StyleOwner         Style = 1
	StyleAlreadyForked Style = 2
	StyleNeedsFork     Style = 3
	StyleTryingFork    Style = 4
)

// Klone is the main entry point for a klone routine. This
// is the procedural logic for "kloning" a git repository.
func Klone(name string) error {
	local.Printf("Kloning [%s]", name)
	provider, err := NewProviderAlpha1()
	if err != nil {
		return err
	}
	local.Printf("Loading git configuration")
	cfg, err := provider.GetGitConfig()
	if err != nil {
		return err
	}
	local.Printf("Loading server")
	srv, err := provider.GetGitServer()
	if err != nil {
		return err
	}
	local.Printf("Parsing credentials")
	crds, err := srv.GetCredentials()
	if err != nil {
		return err
	}
	local.Printf("Authenticating")
	err = srv.Authenticate(crds)
	if err != nil {
		return err
	}
	local.Printf("Reticulating splines")

	var repo kloneprovider.Repo
	// Logic for "kops" and "kris-nova/kops" queries
	if strings.Contains(name, "/") {
		spl := strings.Split(name, "/")
		if len(spl) != 2 {
			return fmt.Errorf("Invalid repository name: %s", name)
		}
		owner := spl[0]
		rname := spl[1]
		repo, err = srv.GetRepoByOwner(owner, rname)
		if err != nil {
			return err
		}
	} else {
		repo, err = srv.GetRepo(name)
		if err != nil {
			return err
		}
	}

	if repo == nil {
		local.Printf("Unable to lookup repo: %s", name)
		return fmt.Errorf("Invalid repository name: %s", name)
	}
	local.PrintExclaimf("Found repository [%s/%s]!", repo.Owner(), repo.Name())

	// Style
	var s Style
	if (repo.Owner() == srv.OwnerName()) && (repo.ForkedFrom() == nil) {
		// It's ours, and we have no parent - just a normal klone
		local.Printf("[OWNER] klone [%s/%s]", repo.Owner, repo.Name())
		s = StyleOwner
	} else if (repo.Owner() == srv.OwnerName()) && (repo.ForkedFrom() != nil) {
		// It's ours, and we have a parent - so we are kloning a fork
		local.Printf("[ALREADY-FORKED] klone [%s/%s] forked from [%s/%s]", repo.Owner(), repo.Name(), repo.ForkedFrom().Owner(), repo.ForkedFrom().Name())
		s = StyleAlreadyForked
	} else if (repo.Owner() != srv.OwnerName()) && (repo.ForkedFrom() == nil) {
		// It's not ours, and we have no parent. We are totally going to fork this repo.
		local.Printf("[NEEDS-FORK] klone [%s/%s] forked from [%s/%s]", srv.OwnerName(), repo.Name(), repo.Owner(), repo.Name())
		s = StyleNeedsFork
	} else if (repo.Owner() != srv.OwnerName()) && (repo.ForkedFrom() != nil) {
		fmt.Println(repo.ForkedFrom())
		// It's not ours (but maybe we have access) and we have a parent
		local.Printf("[TRYING-FORK] klone [%s/%s] forked from [%s/%s]", srv.OwnerName(), repo.Name(), repo.ForkedFrom().Owner(), repo.ForkedFrom().Name())
		s = StyleTryingFork
	} else {
		// We should never get here.. but still erroring just in case
		local.PrintFatal("Unable to parse kloning style! Major error!")
	}

	if err := kloneRepo(repo, cfg, s); err != nil {
		return err
	}

	return nil
}

// KloneRepo is the pattern that holds the "klone" operation together
// 1. Find our Kloner
// 2. Add our Context
// 3. Klone :)
func kloneRepo(repo kloneprovider.Repo, cfg kloneprovider.GitConfig, style Style) error {

	kloner, err := getKloner(repo)
	if err != nil {
		return err
	}
	ctx, err := getContext()
	if err != nil {
		return err
	}
	kloner.SetContext(ctx)
	if err = kloner.Klone(); err != nil {
		return err
	}
	return nil
}

// getKloner will "find" the kloner implementation we should use
// for this repo. The parsing logic here favors .Klonefile's and
// if a .Klonefile is detected it will always override other
// configuration.
// Todo (@kris-nova) Do we want to flip this logic? Hrmm..
func getKloner(repo kloneprovider.Repo) (kloners.Kloner, error) {
	var kloner kloners.Kloner
	klonefile := repo.GetKlonefile()
	if len(klonefile) > 1 {
		// Somebody is using a .Klonefile
		k, err := getKlonerFromKlonefile(klonefile)
		if err != nil {
			return nil, err
		}
		kloner = k
	} else {
		k, err := getKlonerFromRepo(repo)
		if err != nil {
			return nil, err
		}
		kloner = k
	}
	return kloner, nil
}

// getContext is a function that will handle finding context information
// at runtime.
// Todo (@kris-nova) we should have a few ways of getting context.
// Also should this be a proper context?
// 1. Detecting it in memory
// 2. Passed in a context file on local disk
// 3. URL
// 4. Command line flags to populate elements of context
func getContext() (kloners.KlonerContext, error) {
	// Todo (@kris-nova) we are defaulting to simple to get this compiling
	simpleCtx := &simple.KlonerContext{}
	return simpleCtx, nil
}

// getKlonerFromRepo will select a Kloner based on metrics from the repo.
// (We look at things like primary language, and file extensions)
func getKlonerFromRepo(repo kloneprovider.Repo) (kloners.Kloner, error) {
	// Todo (@kris-nova) we are defaulting to simple to get this compiling
	simpleKloner := simple.NewKloner()
	return simpleKloner, nil
}

// getKlonerFromKlonefile will take a raw slice of bytes (from a .Klonefile in the repo)
// and attempt to parse the .Klonefile for metrics about which Kloner to use
func getKlonerFromKlonefile(klonefile []byte) (kloners.Kloner, error) {
	// Todo (@kris-nova) we are defaulting to simple to get this compiling
	simpleKloner := simple.NewKloner()
	return simpleKloner, nil
}
