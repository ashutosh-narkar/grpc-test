package attribute

import (
	"bytes"
	"fmt"
	"sort"
	"sync"
)

// MutableBag is a generic mechanism to read and write a set of attributes.
//
// Bags can be chained together in a parent/child relationship. A child bag
// represents a delta over a parent. By default a child looks identical to
// the parent. But as mutations occur to the child, the two start to diverge.
// Resetting a child makes it look identical to its parent again.
type MutableBag struct {
	parent Bag
	values map[string]interface{}
}

var mutableBags = sync.Pool{
	New: func() interface{} {
		return &MutableBag{
			values: make(map[string]interface{}),
		}
	},
}

// GetMutableBag returns an initialized bag.
//
// Bags can be chained in a parent/child relationship. You can pass nil if the
// bag has no parent.
//
// When you are done using the mutable bag, call the Done method to recycle it.
func GetMutableBag(parent Bag) *MutableBag {
	mb := mutableBags.Get().(*MutableBag)

	if parent == nil {
		mb.parent = empty
	} else {
		mb.parent = parent
	}

	return mb
}

// Done indicates the bag can be reclaimed.
func (mb *MutableBag) Done() {
	// prevent use of a bag that's in the pool
	if mb.parent == nil {
		panic(fmt.Errorf("attempt to use a bag after its Done method has been called"))
	}

	mb.parent = nil
	mb.Reset()
	mutableBags.Put(mb)
}

// Reset removes all local state.
func (mb *MutableBag) Reset() {
	// my kingdom for a clear method on maps!
	for k := range mb.values {
		delete(mb.values, k)
	}
}

// Get returns an attribute value.
func (mb *MutableBag) Get(name string) (interface{}, bool) {
	// prevent use of a bag that's in the pool
	if mb.parent == nil {
		panic(fmt.Errorf("attempt to use a bag after its Done method has been called"))
	}

	var r interface{}
	var b bool
	if r, b = mb.values[name]; !b {
		r, b = mb.parent.Get(name)
	}
	return r, b
}

// DebugString prints out the attributes from the parent bag, then
// walks through the local changes and prints them as well.
func (mb *MutableBag) DebugString() string {
	if len(mb.values) == 0 {
		return mb.parent.DebugString()
	}

	var buf bytes.Buffer
	buf.WriteString(mb.parent.DebugString())
	buf.WriteString("---\n")

	keys := make([]string, 0, len(mb.values))
	for key := range mb.values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		buf.WriteString(fmt.Sprintf("%-30s: %v\n", key, mb.values[key]))
	}
	return buf.String()
}
