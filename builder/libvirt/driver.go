package libvirt

// Driver is the interface
// that has to be implemented to communicate
// in order to communicate with Libvirt host
type Driver interface {
	// GetPool returns a pool matching the string
	// this prevents duplication error
	GetPool(string) (Pool, error)

	// CreatePool creates a pool on the libvirt
	// host
	CreatePool(Pool) (Pool, error)

	// DeletePool delete a pool
	DeletePool(string) error

	// GetVolume returns a volume matching the string
	// this prevents duplication error
	GetVolume(string, string) (Volume, error)

	// CreateVolume creates a volume on the libvirt
	// host
	CreateVolume(string, Volume) (Volume, error)

	// DeleteVolume delete a volume
	DeleteVolume(string, string) error

	// GetNetwork returns a volume matching the string
	// this prevents duplication error
	GetNetwork(string) (Network, error)

	// CreateNetwork creates a volume on the libvirt
	// host
	CreateNetwork(Network) (Network, error)

	// DeleteNetwork delete a volume
	DeleteNetwork(string) error

	// Create a domain
	CreateDomain(Domain) (Domain, error)
}

type Pool interface {
	// GetName returns the name of
	// the pool
	GetName() (string, error)

	// GetXML returns the XML template
	// rendered with the actual values
	GetXML() (string, error)
}

type Volume interface {
	// GetName returns the name of
	// the pool
	GetName() (string, error)

	// GetXML returns the XML template
	// rendered with the actual values
	GetXML() (string, error)
}

type Network interface {
	// GetName returns the name of
	// the network
	GetName() (string, error)

	// GetXML returns the XML template
	// rendered with the actual values
	GetXML() (string, error)
}

type Domain interface {
	// GetIP returns the IP of the domain
	GetIP() (string, error)

	// GetName returns the name of the domain
	GetName() (string, error)

	// GetXML returns the XML template
	// of the domain
	GetXML() (string, error)
}
