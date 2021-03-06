package libvirt

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/pkg/errors"

	"github.com/hashicorp/packer/packer"
	"libvirt.org/libvirt-go"

	libvirtXML "github.com/tormath1/packer-builder-libvirt/builder/libvirt/internal/template"
)

// driverLibvirt is the implementation
// of Driver interface
type driverLibvirt struct {
	conn *libvirt.Connect
	ui   packer.Ui
}

type poolLibvirt struct {
	PoolName       string
	PoolType       string
	PoolTargetPath string
}

type volLibvirt struct {
	VolumeName             string
	VolumeAllocation       int
	VolumeCapacityUnit     string
	VolumeCapacity         int
	VolumeTargetPath       string
	VolumeTargetFormatType string
	VolumeType             string
}

type networkLibvirt struct {
	NetworkName       string
	NetworkMode       string
	NetworkBridgeName string
}

type domainLibvirt struct {
	DomainName       string
	DomainType       string
	DomainMemoryUnit string
	DomainMemory     int
	DomainVCPU       int
	DomainDiskType   string
	PoolName         string
	VolumeName       string
	ISOProtoURL      string
	ISOPathURL       string
	ISOHostURL       string
	ISOPortURL       string
	IP               string
}

func (dom *domainLibvirt) GetName() (string, error) {
	return dom.DomainName, nil
}

func (dom *domainLibvirt) GetXML() (string, error) {
	var domXML bytes.Buffer
	tmpl, err := template.
		New("domain").
		Parse(libvirtXML.DomainXML)
	if err != nil {
		return "", errors.Wrap(err, "unable to parse template")
	}
	if err := tmpl.Execute(&domXML, dom); err != nil {
		return "", errors.Wrap(err, "unable to execute template")
	}
	return domXML.String(), nil
}

func (dom *domainLibvirt) GetIP() (string, error) {
	return dom.IP, nil
}

func (nl *networkLibvirt) GetName() (string, error) {
	return nl.NetworkName, nil
}

func (vl *volLibvirt) GetName() (string, error) {
	return vl.VolumeName, nil
}

func (nl *networkLibvirt) GetXML() (string, error) {
	var netXML bytes.Buffer
	tmpl, err := template.
		New("net").
		Parse(libvirtXML.NetworkXML)
	if err != nil {
		return "", errors.Wrap(err, "unable to parse template")
	}
	if err := tmpl.Execute(&netXML, nl); err != nil {
		return "", errors.Wrap(err, "unable to execute template")
	}
	return netXML.String(), nil
}

func (vl *volLibvirt) GetXML() (string, error) {
	var volXML bytes.Buffer
	tmpl, err := template.
		New("vol").
		Parse(libvirtXML.VolumeXML)
	if err != nil {
		return "", errors.Wrap(err, "unable to parse template")
	}
	if err := tmpl.Execute(&volXML, vl); err != nil {
		return "", errors.Wrap(err, "unable to execute template")
	}
	return volXML.String(), nil
}

func (dl *driverLibvirt) GetPool(name string) (Pool, error) {
	pool, err := dl.conn.LookupStoragePoolByName(name)
	if err != nil {
		return nil, nil
	}
	defer pool.Free()
	poolName, err := pool.GetName()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get name of pool: %s", name)
	}
	return &poolLibvirt{
		PoolName: poolName,
	}, nil
}

func (dl *driverLibvirt) CreatePool(pool Pool) (Pool, error) {
	poolXML, err := pool.GetXML()
	if err != nil {
		return nil, errors.Wrap(err, "unable to create pool XML")
	}
	p, err := dl.conn.StoragePoolDefineXML(poolXML, 0)
	if err != nil {
		return nil, errors.Wrap(err, "unable to define pool XML")
	}
	if err := p.Create(libvirt.STORAGE_POOL_CREATE_WITH_BUILD); err != nil {
		return nil, errors.Wrap(err, "unable to start pool")
	}
	return pool, nil
}

func (dl *driverLibvirt) DeletePool(name string) error {
	pool, err := dl.conn.LookupStoragePoolByName(name)
	if err != nil {
		return nil
	}
	defer pool.Free()
	if err := pool.Destroy(); err != nil {
		return errors.Wrapf(err, "unable to destroy pool: %s", name)
	}
	if err := pool.Delete(libvirt.STORAGE_POOL_DELETE_NORMAL); err != nil {
		return errors.Wrapf(err, "unable to delete pool: %s", name)
	}
	if err := pool.Undefine(); err != nil {
		return errors.Wrapf(err, "unable to undefine pool: %s", name)
	}
	return nil
}

func (dl *driverLibvirt) GetVolume(pool, name string) (Volume, error) {
	p, err := dl.conn.LookupStoragePoolByName(name)
	if err != nil {
		return nil, nil
	}
	defer p.Free()
	vol, err := p.LookupStorageVolByName(name)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get vol %s of pool %s", name, pool)
	}
	defer vol.Free()
	volName, err := vol.GetName()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get vol %s of pool %s", name, pool)
	}
	return &volLibvirt{
		VolumeName: volName,
	}, nil
}

func (dl *driverLibvirt) CreateVolume(pool string, vol Volume) (Volume, error) {
	p, err := dl.conn.LookupStoragePoolByName(pool)
	if err != nil {
		return nil, nil
	}
	defer p.Free()
	volXML, err := vol.GetXML()
	if err != nil {
		return nil, errors.Wrap(err, "unable to create volume XML")
	}
	v, err := p.StorageVolCreateXML(volXML, libvirt.STORAGE_VOL_CREATE_PREALLOC_METADATA)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create vol")
	}
	defer v.Free()
	return vol, nil
}

func (dl *driverLibvirt) DeleteVolume(pool string, vol string) error {
	p, err := dl.conn.LookupStoragePoolByName(pool)
	if err != nil {
		return nil
	}
	defer p.Free()
	v, err := p.LookupStorageVolByName(vol)
	if err != nil {
		return errors.Wrapf(err, "unable to get vol %s of pool %s", vol, pool)
	}
	defer v.Free()
	if err := v.Delete(libvirt.STORAGE_VOL_DELETE_ZEROED); err != nil {
		return errors.Wrapf(err, "unable to delete volume: %s from pool: %s", vol, pool)
	}
	return nil
}

func (pl *poolLibvirt) GetName() (string, error) {
	return pl.PoolName, nil
}

func (pl *poolLibvirt) GetXML() (string, error) {
	var poolXML bytes.Buffer
	tmpl, err := template.
		New("pool").
		Parse(libvirtXML.PoolXML)
	if err != nil {
		return "", errors.Wrap(err, "unable to parse template")
	}
	if err := tmpl.Execute(&poolXML, pl); err != nil {
		return "", errors.Wrap(err, "unable to execute template")
	}
	return poolXML.String(), nil
}

func (dl *driverLibvirt) GetNetwork(name string) (Network, error) {
	net, _ := dl.conn.LookupNetworkByName(name)
	if net == nil {
		return nil, nil
	}
	defer net.Free()
	return &networkLibvirt{
		NetworkName: name,
	}, nil
}

func (dl *driverLibvirt) CreateNetwork(net Network) (Network, error) {
	netXML, err := net.GetXML()
	if err != nil {
		return nil, errors.Wrap(err, "unable to create network XML")
	}
	if _, err := dl.conn.NetworkCreateXML(netXML); err != nil {
		return nil, errors.Wrap(err, "unable to create network")
	}
	return net, nil
}

func (dl *driverLibvirt) DeleteNetwork(name string) error {
	return nil
}

func NewDriverLibvirt(URI string, ui packer.Ui) (Driver, error) {
	c, err := libvirt.NewConnect(URI)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect: %s", URI)
	}
	ui.Say(fmt.Sprintf("connected to %s", URI))
	return &driverLibvirt{
		conn: c,
		ui:   ui,
	}, nil
}

func (dl *driverLibvirt) CreateDomain(dom Domain) (Domain, error) {
	domXML, err := dom.GetXML()
	if err != nil {
		return nil, errors.Wrap(err, "unable to create domain XML")
	}
	if _, err := dl.conn.DomainCreateXML(domXML, libvirt.DOMAIN_NONE); err != nil {
		return nil, errors.Wrap(err, "unable to create domain")
	}
	// TODO: Get the IP
	return dom, nil
}
