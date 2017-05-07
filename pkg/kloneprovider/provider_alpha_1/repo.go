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
// repo.go is an implementation of a git repository according to klone

package provider_alpha_1

import (
	"github.com/kris-nova/klone/pkg/kloneprovider"
	"github.com/google/go-github/github"
)

type Repo struct {
	impl *github.Repository
}

func (r *Repo) SetImplementation(impl interface{}) {
	gh := impl.(*github.Repository)
	r.impl = gh
}

func (r *Repo) GitCloneUrl() (string, error) {
	return "", nil
}
func (r *Repo) HttpsCloneUrl() (string, error) {
	return "", nil
}
func (r *Repo) Language() (string, error) {
	return "", nil
}
func (r *Repo) Owner() (string, error) {
	return "", nil
}
func (r *Repo) Name() (string, error) {
	return *r.impl.Name, nil
}
func (r *Repo) Description() (string, error) {
	return "", nil
}
func (r *Repo) ForkedFrom() (kloneprovider.Repo, error) {
	return &Repo{}, nil
}
func (r *Repo) GetRepoController() (kloneprovider.RepoController, error) {
	return &RepoController{}, nil
}
func (r *Repo) GetKlonefile() ([]byte, error) {
	return []byte(""), nil
}

type RepoController struct {
	//
}

func (ctl *RepoController) SetRepo(repo kloneprovider.Repo) (error) {
	return nil
}
func (ctl *RepoController) SetRemote(string, string) (error) {
	return nil
}
func (ctl *RepoController) SetInitCommand(kloneprovider.Command) {
}
func (ctl *RepoController) Init() (error) {
	return nil
}
func (ctl *RepoController) SetCloneCommand(kloneprovider.Command) {

}
func (ctl *RepoController) Clone() (error) {
	return nil
}
func (ctl *RepoController) Rsync() (error) {
	return nil
}
