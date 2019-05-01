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
	for i, next := range n.edges {
		if next == nil {
			continue
		}
		(*gvzc)++

		n = next

		(*gvoutput) = append((*gvoutput), fmt.Sprintf("  \"%s\" -> \"[%d] %s\" [label=\"%s\"]", previous, *gvzc, string(n.prefix), string(byte(i))))

		graphviz(r, gvoutput, gvzc, fmt.Sprintf("[%d] %s", *gvzc, string(n.prefix)), n)
	}
}
