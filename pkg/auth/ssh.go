package auth

import (
	"github.com/kris-nova/klone/pkg/local"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

// Private key to use, defined in /cmd
var OptPrivateKey = "~/.ssh/id_rsa"
var PrivateKeyBytesOverride []byte

func GetTransport() (transport.AuthMethod, error) {
	bytes := local.BGetContent(OptPrivateKey)
	if len(PrivateKeyBytesOverride) > 1 {
		bytes = PrivateKeyBytesOverride
	}
	pk, err := ssh.NewPublicKeys("git", bytes, "")
	if err != nil {
		return nil, err
	}
	return pk, nil
}
