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
	"strings"
	"github.com/kris-nova/klone/pkg/kloneprovider"
	"github.com/kris-nova/klone/pkg/klone/kloners/simple"
	"github.com/kris-nova/klone/pkg/klone/kloners"
	"github.com/kris-nova/klone/pkg/local"
)

// Klone is the main entry point for a klone routine. This
// is the procedural logic for "kloning" a git repository.
// The logic here is
// 1. Pick a klone provider
// 2. Connect to the git server
// 3. Parse git configuration
// 4. Attempt to find a repo that matches our input
// 5. Call "kloneRepo" on that repo
func Klone(name string) error {
	local.Printf("kloning [%s]", name)
	namelower := strings.ToLower(name)
	provider, err := NewProviderAlpha1()
	if err != nil {
		return err
	}

	cfg, err := provider.GetGitConfig()
	if err != nil {
		return err
	}

	// Server connection
	srv, err := provider.GetGitServer()
	if err != nil {
		return err
	}
	crds, err := srv.GetCredentials()
	if err != nil {
		return err
	}
	err = srv.Authenticate(crds)
	if err != nil {
		return err
	}

	repos, err := srv.GetRepos()
	if err != nil {
		return err
	}

	// Todo (@kris-nova) We can have a goroutine build a fabulous hash map
	// on repo name and pointer to repo at runtime. We can then use the hash
	// map to find our repo in O(n*log(n)).
	for _, repo := range repos {
		name, err := repo.Name()
		local.Printf("Checking repo: %s", name)
		if err != nil {
			return err
		}
		rlowername := strings.ToLower(name)
		if namelower == rlowername {
			if err = kloneRepo(repo, cfg); err != nil {
				return err
			}

		}
	}
	return nil
}

// KloneRepo is the pattern that holds the "klone" operation together
// 1. Find our Kloner
// 2. Add our Context
// 3. Klone :)
func kloneRepo(repo kloneprovider.Repo, cfg kloneprovider.GitConfig) error {
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
	klonefile, err := repo.GetKlonefile()
	if err != nil {
		return nil, err
	}
	if len(klonefile) > 1 {
		// Somebody is using a .Klonefile
		kloner, err = getKlonerFromKlonefile(klonefile)
		if err != nil {
			return nil, err
		}
	} else {
		kloner, err = getKlonerFromRepo(repo)
		if err != nil {
			return nil, err
		}
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
