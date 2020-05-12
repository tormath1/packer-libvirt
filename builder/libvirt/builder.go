package libvirt

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/common/bootcommand"
	"github.com/hashicorp/packer/common/shutdowncommand"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
)

type Builder struct {
	config Config
	runner multistep.Runner
}

type Config struct {
	common.PackerConfig            `mapstructure:",squash"`
	common.HTTPConfig              `mapstructure:",squash"`
	common.ISOConfig               `mapstructure:",squash"`
	bootcommand.VNCConfig          `mapstructure:",squash"`
	shutdowncommand.ShutdownConfig `mapstructure:",squash"`

	// Libvirt URI to target. https://libvirt.org/uri.html
	LibvirtURI string `mapstructure:"libvirt_URI" required: "true"`

	// PoolType the type of the pool to create. Defaut goes to `dir`
	PoolType string `mapstructure:"pool_type" required: "false"`
	// PoolName the name of the pool to store the created domain volume
	// If the pool does not exist a new one will be created with PoolName, PoolType and PoolTargetPath
	PoolName string `mapstructure:"pool_name" required: "true"`
	// PoolTargetPath the target path of the pool. Where the volume will be stored
	PoolTargetPath string `mapstructure:"pool_target_path" required: "false"`

	// VolumeName the name of the volume to create
	VolumeName string `mapstructure:"volume_name" required: "false"`
	// VolumeType the type of the volume to create (default: "file")
	VolumeType string `mapstructure:"volume_type" required: "false"`
	// VolumeName the allocation of the volume to create
	VolumeAllocation int `mapstructure:"volume_allocation" required: "false"`
	// VolumeCapacityUnit the unit of the volume capacity
	VolumeCapacityUnit string `mapstructure:"volume_capacity_unit" required: "false"`
	// VolumeCapacity the volume capacity
	VolumeCapacity int `mapstructure:"volume_capacity" required: "false"`
	// VolumeTargetFormatType the type of the target format (default: "qcow2")
	VolumeTargetFormatType string `mapstructure:"volume_target_format_type" required: "false"`
	// VolumeTargetPath is the target path (within the PoolTargetPath)
	VolumeTargetPath string `mapstructure:"volume_target_path" required: "false"`

	// NetworkName is the network to use or to be created
	NetworkName string `mapstructure:"network_name" required: "true"`
	// NetworkMode the mode of the network. Default 'nat'
	NetworkMode string `mapstructure:"network_mode" required: "false"`
	// NetworkBridgeName the name of the bridge interface to use. Default 'virbr0'
	NetworkBridgeName string `mapstructure:"network_bridge_name" required: "false"`

	DomainType       string `mapstructure:"domain_type" required: "false"`
	DomainMemoryUnit string `mapstructure:"domain_memory_unit" required: "false"`
	DomainMemory     int    `mapstructure:"domain_memory" required: "false"`
	DomainVCPU       int    `mapstructure:"domain_vcpu" required: "false"`
	DomainDiskType   string `mapstructure:"domain_disk_type" required: "false"`

	ctx interpolate.Context
}

func (b *Builder) Prepare(raws ...interface{}) ([]string, []string, error) {
	err := config.Decode(&b.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &b.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{
				"boot_command",
				"qemuargs",
			},
		},
	}, raws...)
	if err != nil {
		return nil, nil, err
	}
	if len(b.config.PoolType) == 0 {
		b.config.PoolType = "dir"
	}
	if len(b.config.PoolTargetPath) == 0 {
		b.config.PoolTargetPath = fmt.Sprintf("/var/lib/libvirt/%s", b.config.PoolName)
	}
	if len(b.config.VolumeType) == 0 {
		b.config.VolumeType = "file"
	}
	if len(b.config.VolumeCapacityUnit) == 0 {
		b.config.VolumeCapacityUnit = "G"
	}
	if b.config.VolumeCapacity == 0 {
		b.config.VolumeCapacity = 8
	}
	if len(b.config.VolumeTargetPath) == 0 {
		b.config.VolumeTargetPath = fmt.Sprintf("%s/%s", b.config.PoolTargetPath, b.config.VolumeName)
	}
	if len(b.config.VolumeTargetFormatType) == 0 {
		b.config.VolumeTargetFormatType = "qcow2"
	}
	if len(b.config.VolumeName) == 0 {
		b.config.VolumeName = fmt.Sprintf("%s-volume.%s", b.config.PoolName, b.config.VolumeTargetFormatType)
	}
	if len(b.config.NetworkName) == 0 {
		b.config.NetworkName = "packer-network"
	}
	if len(b.config.NetworkMode) == 0 {
		b.config.NetworkMode = "nat"
	}
	if len(b.config.NetworkBridgeName) == 0 {
		b.config.NetworkBridgeName = "virbr0"
	}
	if len(b.config.DomainType) == 0 {
		b.config.DomainType = "kvm"
	}
	if len(b.config.DomainMemoryUnit) == 0 {
		b.config.DomainMemoryUnit = "MiB"
	}
	if b.config.DomainMemory == 0 {
		b.config.DomainMemory = 1024
	}
	if b.config.DomainVCPU == 0 {
		b.config.DomainVCPU = 1
	}
	if len(b.config.DomainDiskType) == 0 {
		b.config.DomainDiskType = "qcow2"
	}
	return nil, nil, nil
}

func (b *Builder) Run(ctx context.Context, ui packer.Ui, hook packer.Hook) (packer.Artifact, error) {
	driver, err := NewDriverLibvirt(b.config.LibvirtURI, ui)
	if err != nil {
		return nil, err
	}
	state := new(multistep.BasicStateBag)
	state.Put("config", &b.config)
	state.Put("driver", driver)
	state.Put("ui", ui)

	steps := []multistep.Step{
		new(stepCreateStorage),
		new(stepCreateNetwork),
		new(stepCreateDomain),
	}
	b.runner = common.NewRunner(steps, b.config.PackerConfig, ui)
	b.runner.Run(ctx, state)
	return nil, nil
}

func (b *Builder) ConfigSpec() hcldec.ObjectSpec {
	return nil
}
