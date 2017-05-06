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
// server.go is a representation of a git server to klone

package provider_alpha_1

import "github.com/kris-nova/klone/pkg/kloneprovider"

type GitServer struct {
	//
}

func (s *GitServer) Authenticate(credentials kloneprovider.GitServerCredentials) {

}
func (s *GitServer) GetUser() (string, error) {
	return "", nil
}
func (s *GitServer) GetRepos() ([]kloneprovider.Repo, error) {
	var repos []kloneprovider.Repo
	return repos, nil
}
func (s *GitServer) GetRepoCursor() (kloneprovider.RepoCursor, error) {
	return &RepoCursor{}, nil
}

type GitServerCredentials struct {
	//
}
