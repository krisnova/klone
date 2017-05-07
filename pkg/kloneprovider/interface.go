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
// interface.go defines the kloneprovider interfaces. These represent the logic we will need
// to work with the klone (and other) routines.

package kloneprovider

// KloneProvider is the core provider for using Klone
type KloneProvider interface {
	GetGitConfig() (GitConfig, error)
	GetGitServer() (GitServer, error)
}

// Command represents a single task to perform in with git
type Command interface {
	Exec()
	GetStdErr() ([]byte, error)
	GetStdOut() ([]byte, error)
	Next() (Command)
}

// GitConfig represents how we interact with git configuration
type GitConfig interface {
	GetConfigPath() (string, error)
	GetConfigBytes() ([]byte, error)
	GetConfigString() (string, error)
	GetConfigDirective(string) (string, error)
	SetConfigDirective(string, string) error
}

// GitOperationFunc is a special type of function that can be
// called to perform a git operation
type GitOperationFunc func() error

// GitOperation represents an operation with the git library
// to perform
type GitOperation interface {
	ExecFunc(f GitOperationFunc)
	Perform() (error)
	Stdout() ([]byte, error)
	Stdin() ([]byte, error)
}

// Repo represents a git repository
type Repo interface {
	GitCloneUrl() (string, error)
	HttpsCloneUrl() (string, error)
	Language() (string, error)
	Owner() (string, error)
	Name() (string, error)
	Description() (string, error)
	ForkedFrom() (Repo, error)
	GetRepoController() (RepoController, error)
	GetKlonefile() ([]byte, error)
	SetImplementation(interface{})
}

// RepoController is how we controll a Git repository on our local filesystem
type RepoController interface {
	SetRepo(Repo) (error)
	SetRemote(string, string) (error)
	SetInitCommand(Command)
	Init() (error)
	SetCloneCommand(Command)
	Clone() (error)
	Rsync() (error)
}

// GitServer represents a git server (like github.com)
type GitServer interface {
	Authenticate(GitServerCredentials) (error)
	GetCredentials() (GitServerCredentials, error)
	GetRepos() ([]Repo, error)
}

// GitServerCredentials represents necessary information to auth with a GitServer
type GitServerCredentials interface {
}
