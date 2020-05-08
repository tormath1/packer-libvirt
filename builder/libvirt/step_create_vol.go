package libvirt

import (
	"context"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type stepCreateVolume struct {
	Debug bool
}

func (s *stepCreateVolume) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)
	config := state.Get("config").(*Config)

	vol, err := driver.GetVolume(config.PoolName, config.VolumeName)
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	if vol != nil {
		return multistep.ActionContinue
	}
	newVol := &volLibvirt{
		VolumeName:         config.VolumeName,
		VolumeAllocation:   config.VolumeAllocation,
		VolumeCapacityUnit: config.VolumeCapacityUnit,
		VolumeTargetPath:   config.VolumeTargetPath,
		VolumeCapacity:     config.VolumeCapacity,
	}
	if _, err := driver.CreateVolume(config.PoolName, newVol); err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	return multistep.ActionContinue
}

func (s *stepCreateVolume) Cleanup(state multistep.StateBag) {
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)
	config := state.Get("config").(*Config)

	vol, err := driver.GetVolume(config.PoolName, config.VolumeName)
	if err != nil {
		ui.Error(err.Error())
		return
	}
	if vol == nil {
		return
	}
	name, err := vol.GetName()
	if err != nil {
		ui.Error(err.Error())
		return
	}
	if err := driver.DeleteVolume(config.PoolName, name); err != nil {
		ui.Error(err.Error())
		return
	}
	return
}
