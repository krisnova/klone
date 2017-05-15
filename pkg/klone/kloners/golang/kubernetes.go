package golang

import (
	"github.com/kris-nova/klone/pkg/kloneprovider"
	"fmt"
)

// Kubernetes does not follow the traditional path logic, so we have to hard code it
func repoToKubernetesPath(repo kloneprovider.Repo) string {
	path := fmt.Sprintf("%s/src/%s/%s", Gopath(), "k8s.io", repo.Name())
	return path
}
