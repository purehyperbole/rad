package rad

// Radix tree
type Radix struct {
	root *Node
}

// New Radix Tree
func New() *Radix {
	return &Radix{
		root: &Node{
			edges: &[256]*Node{},
		},
	}
}

// Insert or update a keypair
func (r *Radix) Insert(key []byte, value interface{}) bool {
	parent, node, pos, dv := r.find(key)

	switch {
	case shouldInsert(key, node, parent, pos, dv):
		return r.insertNode(key, value, parent, node, pos, dv)
	case shouldUpdate(key, node, parent, pos, dv):
		return r.updateNode(key, value, parent, node, pos, dv)
	case shouldSplitThreeWay(key, node, parent, pos, dv):
		return r.splitThreeWay(key, value, parent, node, pos, dv)
	case shouldSplitTwoWay(key, node, parent, pos, dv):
		return r.splitTwoWay(key, value, parent, node, pos, dv)
	}

	return false
}

// MustInsert attempts to insert a value until it is successful
func (r *Radix) MustInsert(key []byte, value interface{}) {
	for !r.Insert(key, value) {
	}
}

// Lookup a value by key
func (r *Radix) Lookup(key []byte) interface{} {
	_, node, pos, _ := r.find(key)

	if node == nil || len(key) > pos {
		return nil
	}

	return node.value
}

func (r *Radix) find(key []byte) (*Node, *Node, int, int) {
	var pos, dv int
	var node, parent *Node

	node = r.root

	for node.next(key[pos]) != nil {
		parent = node
		node = node.next(key[pos])
		pos++

		if len(node.prefix) > 0 {
			dv = divergence(node.prefix, key[pos:])

			if len(node.prefix) > dv {
				return parent, node, pos, dv
			}

			pos = pos + dv
		}

		// if we've found the key, return its parent node so it can be pointed to the new node
		if pos == len(key) {
			return parent, node, pos, dv
		}
	}

	return node, nil, pos, dv
}

func (r *Radix) insertNode(key []byte, value interface{}, parent, node *Node, pos, dv int) bool {
	return parent.swapNext(key[pos], nil, &Node{
		prefix: key[pos+1:],
		value:  value,
	})
}

func (r *Radix) updateNode(key []byte, value interface{}, parent, node *Node, pos, dv int) bool {
	edgePos := pos - (len(node.prefix) + 1)

	return parent.swapNext(key[edgePos], node, &Node{
		prefix: node.prefix,
		value:  value,
		edges:  node.edges,
	})
}

func (r *Radix) splitTwoWay(key []byte, value interface{}, parent, node *Node, pos, dv int) bool {
	var pfx []byte

	// fix issue where key is found, but is occupied by another node with prefix
	if len(key) > pos {
		pfx = key[pos : pos+dv]
	}

	n1 := &Node{
		prefix: pfx,
		value:  value,
	}

	n2 := &Node{
		prefix: node.prefix[dv+1:],
		value:  node.value,
	}

	n1.setNext(node.prefix[dv], n2)

	return parent.swapNext(key[pos-1], node, n1)
}

func (r *Radix) splitThreeWay(key []byte, value interface{}, parent, node *Node, pos, dv int) bool {
	n1 := &Node{
		prefix: node.prefix[:dv],
	}

	n2 := &Node{
		prefix: node.prefix[dv+1:],
		value:  node.value,
		edges:  node.edges,
	}

	n3 := &Node{
		prefix: key[pos+dv+1:],
		value:  value,
	}

	n1.setNext(node.prefix[dv], n2)
	n1.setNext(key[pos+dv], n3)

	return parent.swapNext(key[pos-1], node, n1)
}
