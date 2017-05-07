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
	t.Logf("token; %s", creds.(*GitServerCredentials).AccessToken)
	err = server.Authenticate(creds)
	if err != nil {
		t.Errorf("unable to auth: %v", err)
	}
	_, err = server.GetRepos()
	if err != nil {
		t.Errorf("unable to get repos: %v", err)
	}

}
