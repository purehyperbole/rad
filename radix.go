package rad

import (
	"fmt"
	"strings"
)

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

// Insert or update a keypair
func (r *Radix) Insert(key []byte, value interface{}) {
	parent, pos, dv := r.findInsertionPoint(key)

	parent.print()
	fmt.Println(pos)
	fmt.Println(dv)

	switch {
	case pos == len(key):
		// update!
		fmt.Println("update")
	case pos < len(key) && dv == 0:
		// insert simple
		fmt.Println("insert")
		parent.setNext(key[pos], &Node{
			prefix: key[pos+1:],
			value:  key,
		})
	default:
		fmt.Println("UGH")
	}
}

// Lookup a value by key
func (r *Radix) Lookup(key []byte) interface{} {
	var i, dv int

	n := r.root

	for n.next(key[i]) != nil {
		n = n.next(key[i])
		if n == nil {
			break
		}
		i++

		if len(n.prefix) > 0 {
			dv = divergence(n.prefix, key[i:])

			if len(n.prefix) > dv {
				break
			}

			i = i + dv
		}

		// if we've found the key, break the loop
		if i == len(key) {
			break
		}
	}

	if n == nil || len(key) > i {
		return nil
	}

	return n.value
}

// Graphviz returns a graphviz formatted string of all the nodes in the tree
// this should only be run on trees with relatively few nodes
func (r *Radix) Graphviz() string {
	var gvzc int

	gvoutput := []string{"digraph G {"}

	r.graphviz(&gvoutput, &gvzc, "[-1] ROOT", r.root)

	gvoutput = append(gvoutput, "}")

	return fmt.Sprint(strings.Join(gvoutput, "\n"))
}

func (r *Radix) graphviz(gvoutput *[]string, gvzc *int, previous string, n *Node) {
	for i, next := range n.edges {
		if next != nil {
			(*gvzc)++

			n = next

			(*gvoutput) = append((*gvoutput), fmt.Sprintf("  \"%s\" -> \"[%d] %s\" [label=\"%s\"]", previous, *gvzc, string(n.prefix), string(byte(i))))

			r.graphviz(gvoutput, gvzc, fmt.Sprintf("[%d] %s", *gvzc, string(n.prefix)), n)
		}
	}
}

func (r *Radix) findInsertionPoint(key []byte) (*Node, int, int) {
	var i, dv int
	var parent *Node

	n := r.root

	for n.next(key[i]) != nil {
		parent = n

		n = n.next(key[i])
		i++

		if len(n.prefix) > 0 {
			dv = divergence(n.prefix, key[i:])

			if len(n.prefix) > dv {
				return parent, i, dv
			}

			i = i + dv
		}

		// if we've found the key, return its parent node so it can be pointed to the new node
		if i == len(key) {
			return parent, i, dv
		}
	}

	return n, i, dv
}
