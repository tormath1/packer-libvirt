package main

import (
	"github.com/hashicorp/packer/packer/plugin"
	"github.com/tormath1/packer-builder-libvirt/builder/libvirt"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterBuilder(new(libvirt.Builder))
	server.Serve()
}
