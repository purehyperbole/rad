package rad

import "github.com/purehyperbole/lunar/node"

// returns shared and divergent characters respectively
func divergence(prefix, key []byte) int {
	var i int

	for i < len(key) && i < len(prefix) {
		if key[i] != prefix[i] {
			break
		}
		i++
	}

	return i
}

func splitprefix(prefix []byte) [][]byte {
	var p []byte

	pfxs := make([][]byte, 0, len(prefix)/node.MaxPrefix+1)

	for len(prefix) >= node.MaxPrefix {
		p, prefix = prefix[:node.MaxPrefix], prefix[node.MaxPrefix:]
		pfxs = append(pfxs, p)
	}

	if len(prefix) > 0 {
		pfxs = append(pfxs, prefix[:len(prefix)])
	}

	return pfxs
}
