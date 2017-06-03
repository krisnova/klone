package container

import (
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/command/container"
	cliflags "github.com/docker/cli/cli/flags"
	"github.com/kris-nova/klone/pkg/local"
	"os"
	"strings"
	"fmt"
)

type Options struct {
	Query   string
	Image   string
	name    string
	Command []string
}

func Run(o *Options) error {
	o.init()
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


	o.Command = append([]string{"./BOOTSTRAP.sh"}, o.)
	err := cobra.RunE(cobra, append([]string{o.name}, o.Command...))
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
