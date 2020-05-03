package libvirt

import (
	"context"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/common/bootcommand"
	"github.com/hashicorp/packer/common/shutdowncommand"
	"github.com/hashicorp/packer/packer"
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

	// The Libvirt URI to target. https://libvirt.org/uri.html
	LibvirtURI string `mapstructure:"libvirt_URI" required: "true"`
}

func (b *Builder) Prepare(raws ...interface{}) ([]string, []string, error) {
	return nil, nil, nil
}

func (b *Builder) Run(ctx context.Context, ui packer.Ui, hook packer.Hook) (packer.Artifact, error) {
	return nil, nil
}

func (b *Builder) ConfigSpec() hcldec.ObjectSpec {
	return nil
}
