// Copyright 2019 Pavel Knoblokh. All rights reserved.
// Use of this source code is governed by MIT License
// that can be found in the LICENSE file.

// Package trie implements a thread-safe trie, also known as
// digital tree or prefix tree. It can be used as a drop-in
// replacement for usual Go maps with string keys.
package trie

import "sync"

// A Trie is an ordered tree data structure.
type Trie struct {
	sync.RWMutex
	root *Node
}

type Node struct {
	symbol   rune
	parent   *Node
	children map[rune]*Node
	data     interface{}
}

// NewTrie creates a new empty trie.
func NewTrie() *Trie {
	return &Trie{
		root: &Node{
			children: make(map[rune]*Node),
		},
	}
}

// Insert adds or replaces the data stored at the given key.
func (trie *Trie) Insert(key string, data interface{}) {
	trie.Lock()

	node := trie.root
	for _, r := range key {
		childNode := node.children[r]
		if childNode == nil {
			childNode = &Node{
				symbol:   r,
				parent:   node,
				children: make(map[rune]*Node),
			}
			node.children[r] = childNode
		}
		node = childNode
	}

	node.data = data

	trie.Unlock()
}

// Search returns the data stored at the given key.
func (trie *Trie) Search(key string) interface{} {
	trie.RLock()
	node := trie.root.findNode(key)
	trie.RUnlock()

	if node == nil {
		return nil
	}
	return node.data
}

// HasPrefix returns the map of all the keys and
// their corresponding data for the given key prefix.
func (trie *Trie) HasPrefix(prefix string) map[string]interface{} {
	var results = make(map[string]interface{})

	trie.RLock()
	defer trie.RUnlock()

	node := trie.root.findNode(prefix)
	if node == nil {
		return results
	}

	if node.data != nil {
		results[prefix] = node.data
	}

	var findResults func(*Node, string)
	findResults = func(node *Node, prefix string) {
		for r, childNode := range node.children {
			childPrefix := prefix + string(r)
			if childNode.data != nil {
				results[childPrefix] = childNode.data
			}
			findResults(childNode, childPrefix)
		}
	}
	findResults(node, prefix)

	return results
}

// Delete removes the data stored at the given key and
// returns true on success and false if the key wasn't
// previously set.
func (trie *Trie) Delete(key string) bool {
	trie.Lock()
	defer trie.Unlock()

	node := trie.root.findNode(key)
	if node == nil || node.data == nil {
		return false
	}

	node.data = nil

	for node.data == nil && len(node.children) == 0 && node.parent != nil {
		delete(node.parent.children, node.symbol)
		parent := node.parent
		node.parent = nil
		node = parent
	}

	return true
}

// Ensure it is called inside the mutex lock
func (node *Node) findNode(key string) *Node {
	for _, r := range key {
		node = node.children[r]
		if node == nil {
			return nil
		}
	}
	return node
}
