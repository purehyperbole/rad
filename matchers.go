package rad

func shouldInsert(key []byte, node, parent *Node, pos, dv int) bool {
	return pos < len(key) && node == nil
}

func shouldUpdate(key []byte, node, parent *Node, pos, dv int) bool {
	return len(key) == pos && dv == len(node.prefix) || len(key) == pos && len(node.prefix) == 0
}

func shouldSplitTwoWay(key []byte, node, parent *Node, pos, dv int) bool {
	return (len(key) - (pos + dv)) == 0
}

func shouldSplitThreeWay(key []byte, node, parent *Node, pos, dv int) bool {
	return (len(key) - (pos + dv)) > 0
}
