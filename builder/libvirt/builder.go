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
	return nil, nil
}

func (b *Builder) ConfigSpec() hcldec.ObjectSpec {
	return nil
}
