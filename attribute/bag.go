package attribute

// Bag is a generic mechanism to access a set of attributes.
type Bag interface {
	// Get returns an attribute value.
	Get(name string) (value interface{}, found bool)

	// Names returns the names of all the attributes known to this bag.
	Names() []string

	// Done indicates the bag can be reclaimed.
	Done()

	// DebugString provides a dump of an attribute Bag that avoids affecting the
	// calculation of referenced attributes.
	DebugString() string
}
