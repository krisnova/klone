package e2e

import (
	"testing"
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"github.com/kris-nova/klone/pkg/local"
	"strings"
	"github.com/kris-nova/klone/pkg/klone"
	"os"
	"github.com/kris-nova/klone/pkg/kloneprovider"
	"github.com/kris-nova/klone/pkg/klone/kloners/golang"
	"github.com/kris-nova/klone/pkg/options"
)

var GitServer kloneprovider.GitServer

// TestMain will setup the e2e testing suite by creating a new (and concurrent) connection
// to the Git provider
func TestMain(m *testing.M) {
	options.R.TestAuthMode = true // Enable test mode auth
	provider := klone.NewGithubProvider()
	gitServer, err := provider.NewGitServer()
	if err != nil {
		fmt.Printf("Unable to get GitHub server: %v\n", err)
		os.Exit(-1)
	}
	GitServer = gitServer
	os.Exit(m.Run())
}

// TestNewRepoOwnerKlone will ensure a throw away repository is created and then attempt to
// klone the repository.. Will ensure origin is set to the new repository
func TestNewRepoOwnerKlone(t *testing.T) {
	path := fmt.Sprintf("%s/klone-e2e-empty", local.Home())

	repo, err := GitServer.GetRepoByOwner(GitServer.OwnerName(), "klone-e2e-empty")
	if err != nil && !strings.Contains(err.Error(), "404 Not Found") {
		t.Fatalf("Unable to attempt to search for repo: %v", err)
	}
	if repo != nil && repo.Owner() == GitServer.OwnerName() {
		_, err := GitServer.DeleteRepo("klone-e2e-empty")
		if err != nil {
			t.Fatalf("Unable to delete repo: %v", err)
		}
	}
	repo, err = GitServer.NewRepo("klone-e2e-empty", "A throw-away repository created by Klone (@kris-nova)")
	if err != nil {
		t.Fatalf("Unable to create new repo: %v", err)
	}
	err = IdempotentKlone(path, "klone-e2e-empty")
	if err != nil {
		t.Fatalf("Error kloning: %v", err)
	}
	r, err := git.PlainOpen(path)
	if err != nil {
		t.Fatalf("Error opening path: %v", err)
	}
	remotes, err := r.Remotes()
	if err != nil {
		t.Fatalf("Error reading remotes: %v", err)
	}
	originOk := false
	for _, remote := range remotes {
		rspl := strings.Split(remote.String(), "\t")
		if len(rspl) < 3 {
			t.Fatalf("Invalid remote string: %s", remote.String())
		}
		name := rspl[0]
		url := rspl[1]
		if strings.Contains(name, "origin") && strings.Contains(url, fmt.Sprintf("git://github.com/%s/klone-e2e-empty.git", GitServer.OwnerName())) {
			originOk = true
		}
	}
	if originOk == false {
		t.Fatal("Error detecting remote [origin]")
	}
}

// TestGoLanguageNeedsFork will attempt to klone a repository that we KNOW the language for (Go)
// The test also handles recursively removing any local files as well as ensuring a GitHub fork
// is removed before running the test. This test (by design) will use the Golang kloner.
func TestGoLanguageNeedsFork(t *testing.T) {
	path := fmt.Sprintf("%s/src/%s/%s/klone-e2e-go", golang.Gopath(), GitServer.GetServerString(), GitServer.OwnerName())
	t.Logf("Test path: %s", path)
	repo, err := GitServer.GetRepoByOwner(GitServer.OwnerName(), "klone-e2e-go")
	if err != nil {
		t.Fatalf("Unable to attempt to search for repo: %v", err)
	}
	if repo != nil && repo.Owner() == GitServer.OwnerName() {
		_, err := GitServer.DeleteRepo("klone-e2e-go")
		if err != nil {
			t.Fatalf("Unable to delete repo: %v", err)
		}
	}
	err = IdempotentKlone(path, "Nivenly/klone-e2e-go")
	if err != nil {
		t.Fatalf("Error kloning: %v", err)
	}
	r, err := git.PlainOpen(path)
	if err != nil {
		t.Fatalf("Error opening path: %v", err)
	}
	remotes, err := r.Remotes()
	if err != nil {
		t.Fatalf("Error reading remotes: %v", err)
	}
	originOk, upstreamOk := false, false
	for _, remote := range remotes {
		rspl := strings.Split(remote.String(), "\t")
		if len(rspl) < 3 {
			t.Fatalf("Invalid remote string: %s", remote.String())
		}
		name := rspl[0]
		url := rspl[1]
		if strings.Contains(name, "origin") && strings.Contains(url, fmt.Sprintf("git://github.com/%s/klone-e2e-go.git", GitServer.OwnerName())) {
			originOk = true
		}
		if strings.Contains(name, "upstream") && strings.Contains(url, fmt.Sprintf("git://github.com/%s/klone-e2e-go.git", "Nivenly")) {
			upstreamOk = true
		}
	}
	if originOk == false {
		t.Fatal("Error detecting remote [origin]")
	}
	if upstreamOk == false {
		t.Fatal("Error detecting remote [upstream]")
	}
}

// TestUnknownLanguageNeedsFork will attempt to klone a repository that we KNOW we won't
// be able to detect a language for. The test also handles recursively removing any local files
// as well as ensuring a GitHub fork is removed before running the test. This test (by design)
// will use the simple kloner
func TestUnknownLanguageNeedsFork(t *testing.T) {
	path := fmt.Sprintf("%s/klone-e2e-unknown", local.Home())
	t.Logf("Test path: %s", path)
	repo, err := GitServer.GetRepoByOwner(GitServer.OwnerName(), "klone-e2e-unknown")
	if err != nil {
		t.Fatalf("Unable to attempt to search for repo: %v", err)
	}
	if repo != nil && repo.Owner() == GitServer.OwnerName() {
		_, err := GitServer.DeleteRepo("klone-e2e-unknown")
		if err != nil {
			t.Fatalf("Unable to delete repo: %v", err)
		}
	}
	err = IdempotentKlone(path, "Nivenly/klone-e2e-unknown")
	if err != nil {
		t.Fatalf("Error kloning: %v", err)
	}
	r, err := git.PlainOpen(path)
	if err != nil {
		t.Fatalf("Error opening path: %v", err)
	}
	remotes, err := r.Remotes()
	if err != nil {
		t.Fatalf("Error reading remotes: %v", err)
	}
	originOk, upstreamOk := false, false
	for _, remote := range remotes {
		rspl := strings.Split(remote.String(), "\t")
		if len(rspl) < 3 {
			t.Fatalf("Invalid remote string: %s", remote.String())
		}
		name := rspl[0]
		url := rspl[1]
		if strings.Contains(name, "origin") && strings.Contains(url, fmt.Sprintf("git://github.com/%s/klone-e2e-unknown.git", GitServer.OwnerName())) {
			originOk = true
		}
		if strings.Contains(name, "upstream") && strings.Contains(url, fmt.Sprintf("git://github.com/%s/klone-e2e-unknown.git", "Nivenly")) {
			upstreamOk = true
		}
	}
	if originOk == false {
		t.Fatal("Error detecting remote [origin]")
	}
	if upstreamOk == false {
		t.Fatal("Error detecting remote [upstream]")
	}
}
