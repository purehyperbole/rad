package rad

import (
	"fmt"
	"strings"
	"sync/atomic"
	"unsafe"
)

// Node stores all leaf data
type Node struct {
	edges  *unsafe.Pointer
	prefix []byte
	value  interface{}
}

func (n *Node) next(b byte) *Node {
	edges := (*[256]unsafe.Pointer)(atomic.LoadPointer(n.edges))
	if edges == nil {
		return nil
	}

	return (*Node)(atomic.LoadPointer(&edges[b]))
}

func (n *Node) setNext(b byte, node *Node) {
	edges := (*[256]unsafe.Pointer)(*n.edges)

	if edges == nil {
		edges = &[256]unsafe.Pointer{}
		n.edges = upptr(unsafe.Pointer(edges))
	}

	edges[b] = unsafe.Pointer(node)
}

func (n *Node) swapNext(b byte, existing, next *Node) bool {
	edges := (*[256]unsafe.Pointer)(atomic.LoadPointer(n.edges))

	if edges == nil {
		edges = n.setupEdges()
	}

	old := unsafe.Pointer(existing)
	new := unsafe.Pointer(next)
	return atomic.CompareAndSwapPointer(&edges[b], old, new)
}

func (n *Node) setupEdges() *[256]unsafe.Pointer {
	edges := &[256]unsafe.Pointer{}

	// swap edges and ignore if it fails
	old := unsafe.Pointer(nil)
	new := unsafe.Pointer(edges)

	if atomic.CompareAndSwapPointer(n.edges, old, new) {
		return edges
	}

	return (*[256]unsafe.Pointer)(atomic.LoadPointer(n.edges))
}

func (n *Node) print() {
	output := []string{"{"}

	output = append(output, fmt.Sprintf("	Prefix Length: %d", len(n.prefix)))
	output = append(output, fmt.Sprintf("	Prefix: %s", string(n.prefix)))
	output = append(output, fmt.Sprintf("	Value: %d", n.value))

	output = append(output, "	Edges: [")

	if n.edges != nil {
		for i := 0; i < 256; i++ {
			edge := n.next(byte(i))
			if edge != nil {
				output = append(output, fmt.Sprintf("		%s: %s", string(byte(i)), edge.prefix))
			}
		}
	}

	output = append(output, "	]")
	output = append(output, "}")

	fmt.Println(strings.Join(output, "\n"))
}

func upptr(p unsafe.Pointer) *unsafe.Pointer {
	return &p
}
