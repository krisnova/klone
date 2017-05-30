package container

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/kris-nova/klone/pkg/local"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type Options struct {
	Image   string
	Command []string
}

func Run(o *Options) error {
	local.Printf("Running klone in container with image [%s]", o.Image)
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	cl, err := cli.ImagePull(ctx, o.Image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(cl)
	bLen := len(bytes)
	if bLen > 208 {
		local.Printf("Successfully pulled [%d] bytes for image [%s]", bLen, o.Image)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: o.Image,
		Cmd:   o.Command,
	}, nil, nil, "")
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	tch := make(chan bool)
	go func() {
		time.Sleep(time.Second * 256)
		tch <- true
	}()

	attachOptions := types.ContainerAttachOptions{}
	hj, err := cli.ContainerAttach(ctx, resp.ID, attachOptions)
	if err != nil {
		return err
	}
	bb, err := ioutil.ReadAll(hj.Reader)
	if err != nil {
		return err
	}

	fmt.Println(string(bb))

	okChan, errChan := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case <-okChan:
		local.Printf("Container exited")
		break
	case err := <-errChan:
		return err
	case <-tch:
		return errors.New("Timeout while waiting for klone")
	}
	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, out)
	return nil
}
