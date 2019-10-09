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
		node, ok := trie.Children[r]
		if !ok {
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

func (trie *Node) Search(key string) (interface{}, bool) {
	node, ok := trie.findNode(key)
	if ok && node.Data != nil {
		return node.Data, true
	}
	return nil, false
}

func (trie *Node) HasPrefix(key string) map[string]interface{} {
	var results = make(map[string]interface{})

	pnode, ok := trie.findNode(key)
	if !ok {
		return results
	}

	if pnode.Data != nil {
		results[pnode.key()] = pnode.Data
	}

	var findNodes func(*Node)
	findNodes = func(node *Node) {
		for _, node := range node.Children {
			if node.Data != nil {
				results[node.key()] = node.Data
			}
			findNodes(node)
		}
	}
	findNodes(pnode)

	return results
}

func (trie *Node) Delete(key string) bool {
	node, ok := trie.findNode(key)
	if !ok || node.Data == nil {
		return false
	}

	node.Data = nil

	for node.Data == nil && len(node.Children) == 0 && node.Parent != nil {
		delete(node.Parent.Children, node.Symbol)
		parent := node.Parent
		node.Parent = nil
		node = parent
	}
	return true
}

func (trie *Node) findNode(key string) (*Node, bool) {
	var ok bool
	for _, r := range key {
		trie, ok = trie.Children[r]
		if !ok {
			return nil, false
		}
	}
	return trie, true
}

func (node *Node) key() string {
	key := make([]rune, 0)
	for node.Parent != nil {
		key = append(key, node.Symbol)
		node = node.Parent
	}
	for i := 0; i < len(key)/2; i++ {
		key[i], key[len(key)-1-i] = key[len(key)-1-i], key[i]
	}
	return string(key)
}
