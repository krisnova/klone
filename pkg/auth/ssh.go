package auth

import (
	"github.com/kris-nova/klone/pkg/local"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

func GetTransport() (transport.AuthMethod, error) {
	bytes := local.BGetContent("~/.ssh/id_rsa")
	pk, err := ssh.NewPublicKeys("git", bytes, "")
	if err != nil {
		return nil, err
	}
	return pk, nil
}
