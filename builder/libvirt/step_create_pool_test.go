package libvirt

import (
	"bytes"
	"context"
	"testing"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type mockDriver struct{}

func (md *mockDriver) GetPool(name string) (Pool, error) {
	return nil, nil
}

func (md *mockDriver) CreatePool(p Pool) (Pool, error) {
	return p, nil
}

func (md *mockDriver) DeletePool(name string) error {
	return nil
}

func TestStepCreatePool(t *testing.T) {
	// fake state
	c := &Config{
		LibvirtURI: "foo@bar",
		PoolName:   "my-pool",
	}
	state := new(multistep.BasicStateBag)
	state.Put("config", c)
	state.Put("driver", &mockDriver{})
	state.Put("ui", &packer.BasicUi{
		Reader: new(bytes.Buffer),
		Writer: new(bytes.Buffer),
	})

	t.Run("Success", func(t *testing.T) {
		step := new(stepCreatePool)
		defer step.Cleanup(state)
		if action := step.Run(context.Background(), state); action != multistep.ActionContinue {
			t.Fatalf("bad action: %#v", action)
		}
	})
}
