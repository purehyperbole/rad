package rad

// Node stores all leaf data
type Node struct {
	edge   [256]*Node
	prefix []byte
	value  interface{}
}

func (n *Node) next(b byte) *Node {
	return n.edge[b]
}
