# Packer Libvirt builder
Use your libvirt provider to build images with Packer

![Go](https://github.com/tormath1/packer-libvirt/workflows/Go/badge.svg)

:warning: it's under active development, it's not functional yet. Take a look at the roadmad :relaxed: :warning:

## Roadmap

This roadmap is temporary and more information will be added.

- [x] pool creation
- [ ] volume creation
- [ ] network creation
- [ ] domain creation
- [ ] domain provisioning

## Install

Build the plugin:

```shell
$ make
```

Install it to your plugin default location (`${HOME}/.packer.d/plugins`): 

```shell
$ mv packer-builder-libvirt ${HOME}/.packer.d/plugins
```

