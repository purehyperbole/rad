package rad

import (
	"fmt"
	"strings"
)

// Graphviz returns a graphviz formatted string of all the nodes in the tree
// this should only be run on trees with relatively few nodes
func Graphviz(r *Radix) string {
	var gvzc int

	gvoutput := []string{"digraph G {"}

	graphviz(r, &gvoutput, &gvzc, "[-1] ROOT", r.root)

	gvoutput = append(gvoutput, "}")

	return strings.Join(gvoutput, "\n")
}

func graphviz(r *Radix, gvoutput *[]string, gvzc *int, previous string, n *Node) {
	if n.edges == nil {
		return
	}

	for i := 0; i < 256; i++ {
		next := n.next(byte(i))
		if next == nil {
			continue
		}

		(*gvzc)++

		(*gvoutput) = append((*gvoutput), fmt.Sprintf("  \"%s\" -> \"[%d] %s\" [label=\"%s\"]", previous, *gvzc, string(next.prefix), string(byte(i))))

		graphviz(r, gvoutput, gvzc, fmt.Sprintf("[%d] %s", *gvzc, string(next.prefix)), next)
	}
}
