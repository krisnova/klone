package container

import (
	"fmt"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/command/container"
	cliflags "github.com/docker/cli/cli/flags"
	"github.com/kris-nova/klone/pkg/local"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

type Options struct {
	Query   string
	Image   string
	name    string
	Command []string
}

func Run(o *Options) error {

	o.init()
	err := ensureBootstrapFileLocal()
	if err != nil {
		return fmt.Errorf("Unable to ensure local bootstrap file: %v", err)
	}

	local.Printf("Container: %s", o.name)
	cli := command.NewDockerCli(os.Stdin, os.Stdout, os.Stderr)
	opts := &cliflags.ClientOptions{
		Common: &cliflags.CommonOptions{},
	}
	cli.Initialize(opts)
	cobra := container.NewRunCommand(cli)

	// Todo (@kris-nova) I have opinions and here they are. Let people have their own opinions. (Make this configurable)
	cobra.Flags().Set("name", o.name)
	cobra.Flags().Set("rm", "1")
	cobra.Flags().Set("interactive", "1")
	cobra.Flags().Set("tty", "1")

	// Bootstrap /tmp/klone
	cobra.Flags().Set("volume", fmt.Sprintf("%s:/tmp/klone", path.Dir(bootstrapFile)))

	// Bootstrap ~/.ssh
	cobra.Flags().Set("volume", fmt.Sprintf("%s/.ssh:/root/.ssh", local.Home()))

	// Bootstrap command
	o.Command = append([]string{"bash", "/tmp/klone/BOOTSTRAP.sh", o.Query}, o.Command...)
	err = cobra.RunE(cobra, append([]string{o.name}, o.Command...))
	if err != nil {
		return err
	}
	return nil
}

func (o *Options) init() {
	if strings.Contains(o.Image, ":") {
		spl := strings.Split(o.Image, ":")
		if len(spl) > 1 {
			o.name = spl[0]
			return
		}

	}
	o.name = o.Image
}

var bootstrapFile = fmt.Sprintf("%s/.klone/BOOTSTRAP.sh", local.Home())
var remoteBootstrapFileUrl = "https://raw.githubusercontent.com/kris-nova/klone/master/hack/BOOTSTRAP.sh"

func ensureBootstrapFileLocal() error {
	wd, err := os.Getwd()
	if err != nil {
		wd = ""
	}
	if _, err := os.Stat(fmt.Sprintf("%s/hack", wd)); err == nil {
		local.PrintExclaimf("Found local hack directory")
		localBootstrapFile := fmt.Sprintf("%s/hack/BOOTSTRAP.sh", wd)
		local.SPutContent(local.Version, fmt.Sprintf("%s/hack/version", wd))
		local.SPutContent(local.SGetContent(fmt.Sprintf("%s/.klone/auth", local.Home())), fmt.Sprintf("%s/hack/auth", wd))
		bootstrapFile = localBootstrapFile
		return nil
	}

	r, err := http.Get(remoteBootstrapFileUrl)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	content := string(body)
	local.SPutContent(content, bootstrapFile)
	return nil
}
