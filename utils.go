package rad

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
