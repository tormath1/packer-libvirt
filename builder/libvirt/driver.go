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
}

type Pool interface {
	// GetName returns the name of
	// the pool
	GetName() (string, error)

	// GetXML returns the XML template
	// rendered with the actual values
	GetXML() (string, error)
}
