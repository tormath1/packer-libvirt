# Packer Libvirt builder
Use your libvirt provider to build images with Packer

![Go](https://github.com/tormath1/packer-libvirt/workflows/Go/badge.svg)
[![#packer-libvirt](https://img.shields.io/badge/IRC-%23packer--libvirt-blue)](https://webchat.freenode.net/)

:warning: it's under active development, it's not functional yet. Take a look at the roadmad :relaxed: :warning:

## Roadmap

This roadmap is temporary and more information will be added.

- [x] pool creation
- [x] volume creation
- [x] network creation
- [x] domain creation
- [ ] domain provisioning (initial provisionning with VNC)

## Install

You need to create the directory of packer custom plugins:

```
$ mkdir ${HOME}/.packer.d/plugins
```

Buil and install the plugin:

```shell
$ make install
```

## Parameters

The list of dedicated Libvirt parameters supported by the plugin.

| name                      | description                                                        | required | default                              |
|---------------------------|--------------------------------------------------------------------|----------|--------------------------------------|
| libvirt_URI               | URI of the Libvirt server. https://libvirt.org/uri.html            | true     |                                      |
| pool_type                 | Type of the pool to create. Only "dir" is supported for now        | false    | dir                                  |
| pool_name                 | Name of the pool to create / use in order to store the volume      | false    | pool-packer                          |
| pool_target_path          | The target path of the pool                                        | false    | /var/lib/libvirt/{{pool_name}}       |
| volume_name               | Name of the volume to create in the pool                           | false    | {{pool_name}}-volume                 |
| volume_type               | The type of the volume to create. Only "file" is supported for now | false    | file                                 |
| volume_allocation         | The allocation of the volume to create                             | false    | 0                                    |
| volume_capacity_unit      | The unit of the volume capacity                                    | false    | G (for GB)                           |
| volume_capacity           | The capacity of the volume in {{volume_capacity_unit}}             | false    | 8                                    |
| volume_target_format_type | The type of the target format. Only "qcow2" is supported for now   | false    | qcow2                                |
| volume_target_path        | The target path                                                    | false    | {{pool_target_path}}/{{volume_name}} |
| network_name              | The name of the network to use / create                            | true     |                                      |
| network_mode              | The mode of the forwarding method. Now, only 'nat' is supported    | false    | nat                                  |
| network_bridge_name       | The name of the bridge interface to use                            | false    | virbr0                               |
| domain_type               | The type of the domain. Only KVM is supported for now              | false    | kvm                                  |
| domain_memory_unit        | The unit of the memory for the domain                              | false    | MiB                                  |
| domain_memory             | The memory of the domain                                           | false    | 1024                                 |
| domain_vcpu               | The number of VCPU                                                 | false    | 1                                    |
| domain_disk_type          | The type of the disk to use. Only "qcow2" is supported.            | false    | qcow2                                |
