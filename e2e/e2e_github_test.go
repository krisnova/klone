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
)

var GitServer kloneprovider.GitServer

// TestMain will setup the e2e testing suite by creating a new (and concurrent) connection
// to the Git provider
func TestMain(m *testing.M) {
	provider, err := klone.NewGithubProvider()
	if err != nil {
		fmt.Printf("Unable to get klone provider: %v\n", err)
		os.Exit(-1)
	}
	gitServer, err := provider.NewGitServer()
	if err != nil {
		fmt.Printf("Unable to get GitHub server: %v\n", err)
		os.Exit(-1)
	}
	crds, err := gitServer.GetCredentials()
	if err != nil {
		fmt.Printf("Unable to get GitHub credentials: %v\n", err)
		os.Exit(-1)
	}
	err = gitServer.Authenticate(crds)
	if err != nil {
		fmt.Printf("Unable to get authenticate against GitHub: %v\n", err)
		os.Exit(-1)
	}
	GitServer = gitServer
	os.Exit(m.Run())
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
