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
	root *node
	size int
	nnum int
	mu   sync.RWMutex
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
	trie.mu.Lock()

	n := trie.root
	for _, r := range key {
		c := n.children[r]
		if c == nil {
			c = &node{
				symbol:   r,
				parent:   n,
				children: make(map[rune]*node),
			}
			n.children[r] = c
			trie.nnum++
		}
		n = c
	}

	n.data = data

	trie.size++

	trie.mu.Unlock()
}

// Search returns the data stored at the given key.
func (trie *Trie) Search(key string) interface{} {
	trie.mu.RLock()
	defer trie.mu.RUnlock()

	n := trie.root.findNode(key)
	if n == nil {
		return nil
	}
	return n.data
}

// WithPrefix returns the map of all the keys and
// their corresponding data for the given key prefix.
func (trie *Trie) WithPrefix(prefix string) map[string]interface{} {
	results := make(map[string]interface{})

	trie.mu.RLock()
	defer trie.mu.RUnlock()

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
		for r, c := range n.children {
			childPrefix := prefix + string(r)
			if c.data != nil {
				results[childPrefix] = c.data
			}
			findResults(c, childPrefix)
		}
	}
	findResults(n, prefix)

	return results
}

// Delete removes the data stored at the given key and
// returns true on success and false if the key wasn't
// previously set.
func (trie *Trie) Delete(key string) bool {
	trie.mu.Lock()
	defer trie.mu.Unlock()

	n := trie.root.findNode(key)
	if n == nil || n.data == nil {
		return false
	}

	n.data = nil

	for n.data == nil && len(n.children) == 0 && n.parent != nil {
		parent := n.parent
		delete(parent.children, n.symbol)
		n.parent = nil
		n = parent
		trie.nnum--
	}

	trie.size--

	return true
}

// Len returns the total number of keys stored in the trie.
func (trie *Trie) Len() int {
	trie.mu.RLock()
	defer trie.mu.RUnlock()
	return trie.size
}

// NodeNum returns the total number of internal nodes
// in the trie, which can be useful for debugging.
func (trie *Trie) NodeNum() int {
	trie.mu.RLock()
	defer trie.mu.RUnlock()
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
