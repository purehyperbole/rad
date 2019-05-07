package rad

import (
	"fmt"
	"strings"
	"sync/atomic"
	"unsafe"
)

// Node stores all leaf data
type Node struct {
	edges  *[256]*Node
	prefix []byte
	value  interface{}
}

func (n *Node) next(b byte) *Node {
	if n.edges == nil {
		return nil
	}

	return n.edges[b]
}

func (n *Node) setNext(b byte, node *Node) {
	if n.edges == nil {
		n.edges = &[256]*Node{}
	}

	n.edges[int(b)] = node
}

func (n *Node) swapNext(b byte, existing, next *Node) bool {
	if n.edges == nil {
		n.edges = &[256]*Node{}
	}

	oPtr := (*unsafe.Pointer)(unsafe.Pointer(&n.edges[b]))
	old := unsafe.Pointer(existing)
	new := unsafe.Pointer(next)
	return atomic.CompareAndSwapPointer(oPtr, old, new)
}

func (n *Node) print() {
	output := []string{"{"}

	output = append(output, fmt.Sprintf("	Prefix Length: %d", len(n.prefix)))
	output = append(output, fmt.Sprintf("	Prefix: %s", string(n.prefix)))
	output = append(output, fmt.Sprintf("	Value: %d", n.value))

	output = append(output, "	Edges: [")

	for char, edge := range n.edges {
		if edge != nil {
			output = append(output, fmt.Sprintf("		%s: %s", string(byte(char)), edge.prefix))
		}
	}

	output = append(output, "	]")
	output = append(output, "}")

	fmt.Println(strings.Join(output, "\n"))
}
