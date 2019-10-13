[![Build Status](https://travis-ci.com/Pashugan/trie.svg?branch=master)](https://travis-ci.com/Pashugan/trie)
[![GoDoc](https://godoc.org/github.com/Pashugan/trie?status.svg)](https://godoc.org/github.com/Pashugan/trie)
[![Go Report Card](https://goreportcard.com/badge/github.com/Pashugan/trie)](https://goreportcard.com/report/github.com/Pashugan/trie)
[![GitHub license](https://img.shields.io/github/license/Pashugan/trie)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/Pashugan/trie)](https://github.com/Pashugan/trie/stargazers)

# Trie

## Overview

Yet another thread-safe Golang [Trie](https://en.wikipedia.org/wiki/Trie) implementation
with the focus on simplicity, performance and support of concurrency. Trie is also known
as Digital Tree or Prefix Tree. It can be used as a drop-in replacement for usual Go maps
with string keys.

## Benchmarks

Benchmarking highly depends on the used dataset, and the provided results should only be
interpreted as an example. The more keys share common prefixes (e.g. as in URLs), the less
memory a trie consumes, and the faster inserts are.

```
BenchmarkHasPrefixTrie-4     536878      3528 ns/op     554 B/op       7 allocs/op
BenchmarkHasPrefixMap-4         294   4656022 ns/op     341 B/op       2 allocs/op
BenchmarkInsertTrie-4        740662      1450 ns/op     102 B/op       1 allocs/op
BenchmarkInsertMap-4        9070458       146 ns/op       1 B/op       0 allocs/op
BenchmarkSearchTrie-4       1000000      1311 ns/op       0 B/op       0 allocs/op
BenchmarkSearchMap-4       10521606       148 ns/op       0 B/op       0 allocs/op
BenchmarkDeleteTrie-4      13413570        90 ns/op       0 B/op       0 allocs/op
BenchmarkDeleteMap-4       66109558        18 ns/op       0 B/op       0 allocs/op
```

## Usage

Download and install the package (or use [Go Modules](https://blog.golang.org/using-go-modules)).
```bash
$ go get github.com/Pashugan/trie
```

```go
package main

import (
	"fmt"

	"github.com/Pashugan/trie"
)

var CityPopulation = []struct {
	City       string
	Population int
}{
	{"Brisbane", 2462637},
	{"Bridgeport", 144900},
	{"Bristol", 463400},
	{"Auckland", 1628900},
}

func main() {
	// Create an empty trie
	cityPop := trie.NewTrie()

	// Insert keys and corresponding values
	for _, item := range CityPopulation {
		cityPop.Insert(item.City, item.Population)
	}

	testKey := "Brisbane"

	// Fetch some data
	fmt.Printf("The population of %v is %v\n", testKey, cityPop.Search(testKey))
	// Output: The population of Brisbane is 2462637

	// Delete the key
	cityPop.Delete(testKey)
	fmt.Printf("The deleted key is now %v\n", cityPop.Search(testKey))
	// Output: The deleted key is now <nil>

	// Fetch keys starting with a prefix
	fmt.Printf("But the other cities on \"Bri\" are still there: %v\n", cityPop.HasPrefix("Bri"))
	// Output: But the other cities on "Bri" are still there: map[Bridgeport:144900 Bristol:463400]

	// Count the length of the trie
	fmt.Printf("The total number of cities left is %v\n", cityPop.Len())
	// Output: The total number of cities left is 3
}
```
