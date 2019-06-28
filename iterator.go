package rad

import "fmt"

// Iterate over every key from a given point
func (r *Radix) Iterate(from []byte, fn func(key []byte, value interface{})) {
	var node *Node

	if len(from) > 0 {
		_, node, _, _ = r.find(from)
	} else {
		node = r.root
	}

	r.iterate(from, node, fn)
}

func (r *Radix) iterate(key []byte, node *Node, fn func(key []byte, value interface{})) {
	if node == nil {
		fmt.Println("NIL")
		return
	}

	if node.edges == nil {
		return
	}

	for i := 0; i < 256; i++ {
		next := node.next(byte(i))
		ckey := make([]byte, len(key))
		copy(ckey, key)

		if next == nil {
			continue
		}

		node = next

		ckey = append(ckey, byte(i))

		if len(node.prefix) > 0 {
			ckey = append(ckey, node.prefix...)
		}

		if node.value != nil {
			fn(ckey, node.value)
		}

		r.iterate(ckey, node, fn)
	}
}
