package libvirt

import (
	"context"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type stepCreatePool struct {
	Debug bool
}

func (s *stepCreatePool) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
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
	return multistep.ActionContinue
}

func (s *stepCreatePool) Cleanup(state multistep.StateBag) {
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)
	config := state.Get("config").(*Config)

	pool, err := driver.GetPool(config.PoolName)
	if err != nil {
		ui.Error(err.Error())
		return
	}
	if pool == nil {
		return
	}
	name, err := pool.GetName()
	if err != nil {
		ui.Error(err.Error())
		return
	}
	if err := driver.DeletePool(name); err != nil {
		ui.Error(err.Error())
		return
	}
	return
}
