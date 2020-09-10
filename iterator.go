package rad

// Iterate over every key from a given point
func (r *Radix) Iterate(from []byte, fn func(key []byte, value Comparable) error) error {
	var node *Node

	if len(from) > 0 {
		_, node, _, _ = r.find(from)
	} else {
		node = r.root
	}

	return r.iterate(from, node, fn)
}

func (r *Radix) iterate(key []byte, node *Node, fn func(key []byte, value Comparable) error) error {
	if node.edges == nil {
		return nil
	}

	for i := 0; i < 256; i++ {
		next := node.next(byte(i))
		if next == nil {
			continue
		}

		ckey := make([]byte, len(key))
		copy(ckey, key)

		ckey = append(ckey, byte(i))

		if len(next.prefix) > 0 {
			ckey = append(ckey, next.prefix...)
		}

		if next.value != nil {
			err := fn(ckey, next.value)
			if err != nil {
				return err
			}
		}

		err := r.iterate(ckey, next, fn)
		if err != nil {
			return err
		}
	}

	return nil
}
