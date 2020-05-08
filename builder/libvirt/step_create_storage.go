package libvirt

import (
	"context"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type stepCreateStorage struct {
	Debug bool
}

func (s *stepCreateStorage) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)
	config := state.Get("config").(*Config)

	pool, err := driver.GetPool(config.PoolName)
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	if pool != nil {
		return multistep.ActionContinue
	}
	newPool := &poolLibvirt{
		PoolName:       config.PoolName,
		PoolTargetPath: config.PoolTargetPath,
		PoolType:       config.PoolType,
	}
	if _, err := driver.CreatePool(newPool); err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	ui.Say("pool created")
	vol, err := driver.GetVolume(config.PoolName, config.VolumeName)
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	if vol != nil {
		return multistep.ActionContinue
	}
	newVol := &volLibvirt{
		VolumeName:             config.VolumeName,
		VolumeAllocation:       config.VolumeAllocation,
		VolumeCapacityUnit:     config.VolumeCapacityUnit,
		VolumeTargetPath:       config.VolumeTargetPath,
		VolumeTargetFormatType: config.VolumeTargetFormatType,
		VolumeCapacity:         config.VolumeCapacity,
		VolumeType:             config.VolumeType,
	}
	if _, err := driver.CreateVolume(config.PoolName, newVol); err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	ui.Say("volume created")
	return multistep.ActionContinue
}

func (s *stepCreateStorage) Cleanup(state multistep.StateBag) {
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
	ui.Say("volume deleted")

	pool, err := driver.GetPool(config.PoolName)
	if err != nil {
		ui.Error(err.Error())
		return
	}
	if pool == nil {
		return
	}
	name, err = pool.GetName()
	if err != nil {
		ui.Error(err.Error())
		return
	}
	if err := driver.DeletePool(name); err != nil {
		ui.Error(err.Error())
		return
	}
	ui.Say("pool deleted")
	return
}
