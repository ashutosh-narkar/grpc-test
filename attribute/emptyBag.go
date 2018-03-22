package attribute

// emptyBag is an attribute bag that is always empty. It is primarily used as a backstop in a
// chain of bags
type emptyBag struct{}

var empty = &emptyBag{}
var emptySlice = []string{}

func (eb *emptyBag) Get(name string) (interface{}, bool) { return nil, false }
func (eb *emptyBag) Names() []string                     { return emptySlice }
func (eb *emptyBag) Done()                               {}
func (eb *emptyBag) DebugString() string                 { return "" }
