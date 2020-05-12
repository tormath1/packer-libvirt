package libvirt

import (
	"context"
	"net/url"
	"strings"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type stepCreateDomain struct {
	Debug bool
}

func (s *stepCreateDomain) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)
	config := state.Get("config").(*Config)

	URL, _ := url.Parse(config.RawSingleISOUrl)
	// TODO: handling error

	dom := &domainLibvirt{
		DomainName:       "packer-domain-01",
		DomainType:       config.DomainType,
		DomainMemoryUnit: config.DomainMemoryUnit,
		DomainMemory:     config.DomainMemory,
		DomainVCPU:       config.DomainVCPU,
		DomainDiskType:   config.DomainDiskType,
		PoolName:         config.PoolName,
		VolumeName:       config.VolumeName,
		ISOProtoURL:      URL.Scheme,
		ISOPathURL:       URL.Path,
	}

	hostPort := strings.Split(URL.Host, ":")
	if len(hostPort) == 2 {
		dom.ISOHostURL = hostPort[0]
		dom.ISOPortURL = hostPort[1]
	} else {
		dom.ISOHostURL = URL.Host
		dom.ISOPortURL = "80"
	}
	if _, err := driver.CreateDomain(dom); err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	ui.Say("domain created")
	return multistep.ActionContinue
}

func (s *stepCreateDomain) Cleanup(state multistep.StateBag) {}
