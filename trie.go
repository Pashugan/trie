package trie

import "sync"

type Trie struct {
	sync.RWMutex
	Root *Node
}

type Node struct {
	Symbol   rune
	Parent   *Node
	Children map[rune]*Node
	Data     interface{}
}

func NewTrie() *Trie {
	return &Trie{
		Root: &Node{
			Children: make(map[rune]*Node),
		},
	}
}

func (trie *Trie) Insert(key string, data interface{}) {
	trie.Lock()

	node := trie.Root
	for _, r := range key {
		childNode := node.Children[r]
		if childNode == nil {
			childNode = &Node{
				Symbol:   r,
				Parent:   node,
				Children: make(map[rune]*Node),
			}
			node.Children[r] = childNode
		}
		node = childNode
	}

	node.Data = data

	trie.Unlock()
}

func (trie *Trie) Search(key string) interface{} {
	trie.RLock()
	node := trie.Root.findNode(key)
	trie.RUnlock()

	if node == nil {
		return nil
	}
	return node.Data
}

func (trie *Trie) HasPrefix(prefix string) map[string]interface{} {
	var results = make(map[string]interface{})

	trie.RLock()
	defer trie.RUnlock()

	node := trie.Root.findNode(prefix)
	if node == nil {
		return results
	}

	if node.Data != nil {
		results[prefix] = node.Data
	}

	var findResults func(*Node, string)
	findResults = func(node *Node, prefix string) {
		for r, childNode := range node.Children {
			childPrefix := prefix + string(r)
			if childNode.Data != nil {
				results[childPrefix] = childNode.Data
			}
			findResults(childNode, childPrefix)
		}
	}
	findResults(node, prefix)

	return results
}

func (trie *Trie) Delete(key string) bool {
	trie.Lock()
	defer trie.Unlock()

	node := trie.Root.findNode(key)
	if node == nil || node.Data == nil {
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

// Ensure it is called inside the mutex lock
func (node *Node) findNode(key string) *Node {
	for _, r := range key {
		node = node.Children[r]
		if node == nil {
			return nil
		}
	}
	return node
}
