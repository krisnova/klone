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
	"github.com/kris-nova/klone/pkg/local"
	"fmt"
	"strings"
)

type Style int

// Klone is the main entry point for a klone routine. This
// is the procedural logic for "kloning" a git repository.
// This will attempt to look up relevant repository information
// and set a kloning "style" for the klone
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
	local.Printf("Found repository [%s/%s]", repo.Owner(), repo.Name())

	kloneable := &Kloneable{
		server: srv,
		config: cfg,
	}

	// Reason about our repository
	if (repo.Owner() == srv.OwnerName()) && (repo.ForkedFrom() == nil) {
		// It's ours, and we have no parent - just a normal klone
		local.Printf("[OWNER] klone [%s/%s]", repo.Owner, repo.Name())
		kloneable.style = StyleOwner
		kloneable.repo = repo
	} else if (repo.Owner() == srv.OwnerName()) && (repo.ForkedFrom() != nil) {
		// It's ours, and we have a parent - so we are kloning a fork
		local.Printf("[ALREADY-FORKED] klone [%s/%s] forked from [%s/%s]", repo.Owner(), repo.Name(), repo.ForkedFrom().Owner(), repo.ForkedFrom().Name())
		kloneable.style = StyleAlreadyForked
		kloneable.repo = repo
	} else if (repo.Owner() != srv.OwnerName()) && (repo.ForkedFrom() == nil) {
		// It's not ours, and we have no parent. We are totally going to fork this repo (as long as we haven't already)
		possible, err := srv.GetRepoByOwner(srv.OwnerName(), repo.Name())
		if err != nil {
			local.RecoverableErrorf("Unable to find [%s/%s]: %v", srv.OwnerName(), repo.Name(), err)
		}
		if possible == nil {
			local.Printf("[NEEDS-FORK] klone [%s/%s] forked from [%s/%s]", srv.OwnerName(), repo.Name(), repo.Owner(), repo.Name())
			kloneable.style = StyleNeedsFork
			kloneable.repo = repo
		} else {
			local.Printf("[ALREADY-FORKED] klone [%s/%s] forked from [%s/%s]", possible.Owner(), possible.Name(), possible.ForkedFrom().Owner(), possible.ForkedFrom().Name())
			kloneable.style = StyleAlreadyForked
			kloneable.repo = possible
		}
	} else if (repo.Owner() != srv.OwnerName()) && (repo.ForkedFrom() != nil) {
		// It's not ours (but maybe we have access) and we have a parent
		local.Printf("[TRYING-FORK] klone [%s/%s] forked from [%s/%s]", srv.OwnerName(), repo.Name(), repo.ForkedFrom().Owner(), repo.ForkedFrom().Name())
		kloneable.style = StyleTryingFork
		kloneable.repo = repo
	} else {
		// We should never get here.. but still erroring just in case
		local.PrintFatal("Unable to parse kloning style! Major error!")
	}

	// We now have something that is Klonable, let's klone it
	err = kloneable.Klone()
	if err != nil {
		// Todo (@kris-nova) Can we please make klone atomic? :)
		local.Printf("Unable to complete klone. Klone does not clean up after itself, there might be incomplete work!")
		return err
	}
	local.PrintExclaimf("Klone completed")
	return nil
}
