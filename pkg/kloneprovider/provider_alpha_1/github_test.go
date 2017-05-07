package provider_alpha_1

import (
	"testing"
)

func TestConnectToGithub(t *testing.T) {
	server := GitServer{}
	creds, err := server.GetCredentials()
	if err != nil {
		t.Errorf("unable to get credentials: %v", err)
	}
	err = server.Authenticate(creds)
	if err != nil {
		t.Errorf("unable to auth: %v", err)
	}
	repos, err := server.GetRepos()
	if err != nil {
		t.Errorf("unable to get repos: %v", err)
	}
	if len(repos) < 1 {
		t.Error("unable to look up repos")
	}
}

// Todo (@kris-nova) we need to come up with a better way of handling this.
// Right now this only works when we use a user/pass and MFA code
//func TestNoAccessToken(t *testing.T) {
//	os.Remove(Cache)
//	server := GitServer{}
//	creds, err := server.GetCredentials()
//	if err != nil {
//		t.Errorf("unable to get credentials: %v", err)
//	}
//	err = server.Authenticate(creds)
//	if err != nil {
//		t.Errorf("unable to auth: %v", err)
//	}
//	newCache := local.SGetContent(Cache)
//	if newCache != "" {
//		t.Error("Unable to write new access token cache")
//	}
//}
