package trie

type Node struct {
	Symbol   rune
	Parent   *Node
	Children map[rune]*Node
	Data     interface{}
}

func NewTrie() *Node {
	return &Node{
		Children: make(map[rune]*Node),
	}
}

func (trie *Node) Insert(key string, data interface{}) {
	for _, r := range key {
		node := trie.Children[r]
		if node == nil {
			node = &Node{
				Symbol:   r,
				Parent:   trie,
				Children: make(map[rune]*Node),
			}
			trie.Children[r] = node
		}
		trie = node
	}

	trie.Data = data
}

func (trie *Node) Search(key string) interface{} {
	trie = trie.findNode(key)
	if trie != nil {
		return trie.Data
	}
	return nil
}

func (trie *Node) HasPrefix(key string) map[string]interface{} {
	var results = make(map[string]interface{})

	pnode := trie.findNode(key)
	if pnode == nil {
		return results
	}

	if pnode.Data != nil {
		results[key] = pnode.Data
	}

	var findNodes func(*Node, string)
	findNodes = func(node *Node, prefix string) {
		for r, node := range node.Children {
			key := prefix + string(r)
			if node.Data != nil {
				results[key] = node.Data
			}
			findNodes(node, key)
		}
	}
	findNodes(pnode, key)

	return results
}

func (trie *Node) Delete(key string) bool {
	trie = trie.findNode(key)
	if trie == nil || trie.Data == nil {
		return false
	}

	trie.Data = nil

	for trie.Data == nil && len(trie.Children) == 0 && trie.Parent != nil {
		delete(trie.Parent.Children, trie.Symbol)
		parent := trie.Parent
		trie.Parent = nil
		trie = parent
	}
	return true
}

func (trie *Node) findNode(key string) *Node {
	for _, r := range key {
		trie = trie.Children[r]
		if trie == nil {
			return nil
		}
	}
	return trie
}
