package container

import (
	"fmt"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/command/container"
	"github.com/docker/cli/cli/command/image"
	cliflags "github.com/docker/cli/cli/flags"
	"github.com/google/uuid"
	"github.com/kris-nova/klone/pkg/local"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

type Options struct {
	Query      string
	Image      string
	Command    []string
	SaveString string

	root string
	name string
}

func Run(o *Options) error {
	o.init()
	err := ensureBootstrapFileLocal()
	if err != nil {
		return fmt.Errorf("Unable to ensure local bootstrap file: %v", err)
	}
	local.Printf("Running container [%s] with image [%s]", o.name, o.Image)
	cli := command.NewDockerCli(os.Stdin, os.Stdout, os.Stderr)
	opts := &cliflags.ClientOptions{
		Common: &cliflags.CommonOptions{},
	}
	cli.Initialize(opts)
	cobra := container.NewRunCommand(cli)

	// ----- env vars
	for _, e := range os.Environ() {
		spl := strings.Split(e, "=")
		k := spl[0]
		v := spl[1]
		if strings.HasPrefix(k, "KLONE_CONTAINER_") {
			newspl := strings.Split(k, "_")
			newk := newspl[len(newspl)-1]
			cobra.Flags().Set("env", fmt.Sprintf("%s=%s", newk, v))
			local.Printf("Passing variable to container $%s='%s'", newk, v)
		}
	}

	// Todo (@kris-nova) I have opinions and here they are. Let people have their own opinions. (Make this configurable)
	cobra.Flags().Set("name", o.name)
	cobra.Flags().Set("interactive", "1")
	cobra.Flags().Set("tty", "1")

	// Bootstrap /tmp/klone
	cobra.Flags().Set("volume", fmt.Sprintf("%s:/tmp/klone", path.Dir(bootstrapFile)))

	// Bootstrap ~/.ssh
	cobra.Flags().Set("volume", fmt.Sprintf("%s/.ssh:/root/.ssh", local.Home()))

	o.Command = append([]string{"bash", "/tmp/klone/BOOTSTRAP.sh", o.Query}, o.Command...)

	err = cobra.RunE(cobra, append([]string{o.Image}, o.Command...))
	if err != nil {
		return err
	}

	if o.SaveString != "" {
		// This means we want to build and push the image to a registry
		err := o.save(cli)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *Options) save(cli *command.DockerCli) error {
	local.Printf("Saving with string [%s]", o.SaveString)
	push := image.NewPushCommand(cli)
	err := push.RunE(push, []string{o.name, "latest"})
	if err != nil {
		return err
	}
	return nil
}

func (o *Options) init() {
	if strings.Contains(o.Image, ":") {
		spl := strings.Split(o.Image, ":")
		if len(spl) > 1 {
			o.root = spl[0]
		}

	}
	if strings.Contains(o.Image, "/") {
		spl := strings.Split(o.Image, "/")
		if len(spl) > 1 {
			o.root = spl[0]
		}

	}
	if o.root == "" {
		o.root = o.Image
	}
	o.root = strings.Replace(o.root, "/", "", -1)
	o.root = strings.Replace(o.root, ":", "", -1)
	o.root = strings.Replace(o.root, "_", "", -1)
	u, err := uuid.NewRandom()
	if err != nil {
		local.RecoverableErrorf("Unable to generate UUID: %v", err)
		o.name = o.root
		return
	}
	o.name = fmt.Sprintf("%s_%s", o.root, u.String())
}

var bootstrapFile = fmt.Sprintf("%s/.klone/BOOTSTRAP.sh", local.Home())
var remoteBootstrapFileUrl = "https://raw.githubusercontent.com/kris-nova/klone/master/hack/BOOTSTRAP.sh"

func ensureBootstrapFileLocal() error {
	wd, err := os.Getwd()
	if err != nil {
		wd = ""
	}
	if _, err := os.Stat(fmt.Sprintf("%s/hack", wd)); err == nil {
		local.PrintExclaimf("Found local hack directory for container bootstrap")
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
