package rad

import (
	"fmt"
	"strings"
	"sync/atomic"
	"unsafe"
)

// Node stores all leaf data
type Node struct {
	edges     unsafe.Pointer
	edgecount int32 // would be nice if this could be uint8
	prefix    []byte
	value     interface{}
}

func (n *Node) next(b byte) *Node {
	edges := (*[256]unsafe.Pointer)(atomic.LoadPointer(&n.edges))
	if edges == nil {
		return nil
	}

	return (*Node)(atomic.LoadPointer(&edges[b]))
}

func (n *Node) setNext(b byte, node *Node) {
	if n.edges == nil {
		n.edges = unsafe.Pointer(&[256]unsafe.Pointer{})
	}

	edges := (*[256]unsafe.Pointer)(n.edges)

	edges[int(b)] = unsafe.Pointer(node)

	n.edgecount++
}

func (n *Node) swapNext(b byte, existing, next *Node) bool {
	edges := (*[256]unsafe.Pointer)(atomic.LoadPointer(&n.edges))

	if edges == nil {
		n.setupEdges()
		edges = (*[256]unsafe.Pointer)(atomic.LoadPointer(&n.edges))
	}

	old := unsafe.Pointer(existing)
	new := unsafe.Pointer(next)
	successful := atomic.CompareAndSwapPointer(&edges[b], old, new)

	if existing == nil && successful {
		atomic.AddInt32(&n.edgecount, 1)
	}

	return successful
}

func (n *Node) setupEdges() {
	// swap edges and ignore if it fails
	old := unsafe.Pointer(nil)
	new := unsafe.Pointer(&[256]unsafe.Pointer{})
	_ = atomic.CompareAndSwapPointer(&n.edges, old, new)
}

func (n *Node) print() {
	output := []string{"{"}

	output = append(output, fmt.Sprintf("	Prefix Length: %d", len(n.prefix)))
	output = append(output, fmt.Sprintf("	Prefix: %s", string(n.prefix)))
	output = append(output, fmt.Sprintf("	Value: %d", n.value))
	output = append(output, fmt.Sprintf("	Edge Count: %d", n.edgecount))

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
