package rad

// Radix tree
type Radix struct {
	root *Node
}

// New Radix Tree
func New() *Radix {
	return &Radix{
		root: &Node{},
	}
}

// Insert a keypair
func (r *Radix) Insert(key []byte, value interface{}) {
	n, i, dv := r.lookup(key)

	switch {
	// found keys prefix differs : split
	case len(n.prefix) > 0 && len(n.prefix) > dv:
		n = r.splitnode(n, dv, key[i:])
	// key matches : update node!
	case i == len(key):
		break
	// found key and its prefix is a sub key : create
	case len(n.prefix) > 0 && dv == len(n.prefix) || i < len(key):
		n = r.createnode(n, key[i:])
	}

	n.value = value
}

func (r *Radix) Lookup(key []byte) interface{} {
	n, pos, _ := r.lookup(key)
	if n == nil || len(key) > pos {
		return nil
	}

	return n.value
}

func (r *Radix) lookup(key []byte) (*Node, int, int) {
	var i, dv int

	n := r.root

	for n.next(key[i]) != nil {
		n = n.next(key[i])
		if n == nil {
			return nil, 0, 0
		}
		i++

		if len(n.prefix) > 0 {
			dv = divergence(n.prefix, key[i:])

			if len(n.prefix) > dv {
				return n, i, dv
			}

			i = i + dv
		}

		// if we've found the key, break the loop
		if i == len(key) {
			break
		}
	}

	return n, i, dv
}

func (r *Radix) createnode(parent *Node, prefix []byte) *Node {
	return nil
}

func (r *Radix) splitnode(parent *Node, dv int, prefix []byte) *Node {
	return nil
}
