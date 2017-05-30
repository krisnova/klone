package container

import "github.com/kris-nova/klone/pkg/local"

var ImageString string

type Options struct {
	Image string
}

func Run(o *Options) error {
	local.Printf("Running klone in container with image [%s]", o.Image)
	return nil
}
