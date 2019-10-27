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
	root *node
	size int
	nnum int
}

type node struct {
	symbol   rune
	parent   *node
	children map[rune]*node
	data     interface{}
}

// NewTrie creates a new empty trie.
func NewTrie() *Trie {
	return &Trie{
		root: &node{
			children: make(map[rune]*node),
		},
		nnum: 1,
	}
}

// Insert adds or replaces the data stored at the given key.
func (trie *Trie) Insert(key string, data interface{}) {
	trie.Lock()

	n := trie.root
	for _, r := range key {
		childNode := n.children[r]
		if childNode == nil {
			childNode = &node{
				symbol:   r,
				parent:   n,
				children: make(map[rune]*node),
			}
			n.children[r] = childNode
			trie.nnum++
		}
		n = childNode
	}

	n.data = data

	trie.size++

	trie.Unlock()
}

// Search returns the data stored at the given key.
func (trie *Trie) Search(key string) interface{} {
	trie.RLock()
	n := trie.root.findNode(key)
	trie.RUnlock()

	if n == nil {
		return nil
	}
	return n.data
}

// HasPrefix returns the map of all the keys and
// their corresponding data for the given key prefix.
func (trie *Trie) HasPrefix(prefix string) map[string]interface{} {
	results := make(map[string]interface{})

	trie.RLock()
	defer trie.RUnlock()

	n := trie.root.findNode(prefix)
	if n == nil {
		return results
	}

	if n.data != nil {
		results[prefix] = n.data
	}

	// Explicit declaration is needed for recursion to work
	var findResults func(*node, string)
	findResults = func(n *node, prefix string) {
		for r, childNode := range n.children {
			childPrefix := prefix + string(r)
			if childNode.data != nil {
				results[childPrefix] = childNode.data
			}
			findResults(childNode, childPrefix)
		}
	}
	findResults(n, prefix)

	return results
}

// Delete removes the data stored at the given key and
// returns true on success and false if the key wasn't
// previously set.
func (trie *Trie) Delete(key string) bool {
	trie.Lock()
	defer trie.Unlock()

	n := trie.root.findNode(key)
	if n == nil || n.data == nil {
		return false
	}

	n.data = nil

	for n.data == nil && len(n.children) == 0 && n.parent != nil {
		delete(n.parent.children, n.symbol)
		parent := n.parent
		n.parent = nil
		n = parent
		trie.nnum--
	}

	trie.size--

	return true
}

// Len returns the total number of keys stored in the trie.
func (trie *Trie) Len() int {
	trie.RLock()
	defer trie.RUnlock()
	return trie.size
}

// NodeNum returns the total number of internal nodes
// in the trie, which can be useful for debugging.
func (trie *Trie) NodeNum() int {
	trie.RLock()
	defer trie.RUnlock()
	return trie.nnum
}

// Ensure it is called inside the mutex lock
func (n *node) findNode(key string) *node {
	for _, r := range key {
		n = n.children[r]
		if n == nil {
			return nil
		}
	}
	return n
}
