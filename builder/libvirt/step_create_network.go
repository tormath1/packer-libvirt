package libvirt

import (
	"context"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type stepCreateNetwork struct {
	Debug bool
}

func (s *stepCreateNetwork) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)
	config := state.Get("config").(*Config)

	net, err := driver.GetNetwork(config.NetworkName)
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	if net != nil {
		return multistep.ActionContinue
	}
	newNet := &networkLibvirt{
		NetworkName:       config.NetworkName,
		NetworkMode:       config.NetworkMode,
		NetworkBridgeName: config.NetworkBridgeName,
	}
	if _, err := driver.CreateNetwork(newNet); err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	ui.Say("network created")
	return multistep.ActionContinue
}

func (s *stepCreateNetwork) Cleanup(state multistep.StateBag) {}
