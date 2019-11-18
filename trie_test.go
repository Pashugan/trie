// Copyright 2019 Pavel Knoblokh. All rights reserved.
// Use of this source code is governed by MIT License
// that can be found in the LICENSE file.
// The fixtures data was kindly borrowed and sampled
// from /usr/share/dict/web2 on MacOS and,
// according to their README, is copyright free.

package trie

import (
	"bufio"
	"compress/bzip2"
	"log"
	"os"
	"reflect"
	"testing"
)

var testData = []struct {
	Key   string
	Value interface{}
}{
	{"foo", 11},
	{"foobar", 111},
	{"bar", 22},
}

var benchData []string

func init() {
	file, err := os.Open("fixtures/words.bz2")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(bzip2.NewReader(file))
	for scanner.Scan() {
		benchData = append(benchData, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

type PrefixMap map[string]interface{}

func (m PrefixMap) HasPrefix(prefix string) map[string]interface{} {
	results := make(map[string]interface{})

	prefixLen := len(prefix)
	for key, value := range m {
		if len(key) >= prefixLen {
			k := key[:prefixLen]
			if k == prefix {
				results[key] = value
			}
		}
	}

	return results
}

func TestMapHasPrefix(t *testing.T) {
	cases := []struct {
		Key           string
		ExpectedValue map[string]interface{}
	}{
		{"f", map[string]interface{}{
			"foo":    11,
			"foobar": 111,
		}},
		{"foo", map[string]interface{}{
			"foo":    11,
			"foobar": 111,
		}},
		{"foob", map[string]interface{}{
			"foobar": 111,
		}},
		{"ba", map[string]interface{}{
			"bar": 22,
		}},
		{"xyz", map[string]interface{}{}},
	}

	m := make(PrefixMap)
	for _, item := range testData {
		m[item.Key] = item.Value
	}

	for _, item := range cases {
		value := m.HasPrefix(item.Key)
		if !reflect.DeepEqual(value, item.ExpectedValue) {
			t.Errorf("Invalid Map prefix values: expected %v, got %v", item.ExpectedValue, value)
		}
	}
}

func getTestTrie() *Trie {
	trie := NewTrie()
	for _, item := range testData {
		trie.Insert(item.Key, item.Value)
	}
	return trie
}

func TestNewTrie(t *testing.T) {
	var trie interface{}

	trie = NewTrie()
	_, ok := trie.(*Trie)
	if !ok {
		t.Errorf("Invalid trie type")
	}
}

func TestEdgeCases(t *testing.T) {
	trie := NewTrie()
	res := trie.Search("xyz")
	if res != nil {
		t.Errorf("Search in empty trie must return nil")
	}

	ok := trie.Delete("xyz")
	if ok {
		t.Errorf("Delete in empty trie must fail")
	}

	res = trie.Search("")
	if res != nil {
		t.Errorf("Empty key in empty trie must return nil")
	}

	trie.Insert("", 1234)
	res = trie.Search("")
	if res != 1234 {
		t.Errorf("Root node must also be able to store data")
	}

	want := map[string]interface{}{
		"":    1234,
		"xyz": "xyz",
	}
	trie.Insert("xyz", "xyz")
	got := trie.HasPrefix("")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("HasPrefix must include root node data if available")
	}
}

func TestInsertAndSearch(t *testing.T) {
	cases := []struct {
		Key           string
		ExpectedValue interface{}
	}{
		{"foo", 11},
		{"foobar", 111},
		{"bar", 22},
		{"foob", nil},
		{"foobarr", nil},
	}

	trie := getTestTrie()

	for _, item := range cases {
		value := trie.Search(item.Key)
		if value != item.ExpectedValue {
			t.Errorf("Invalid value: expected %v, got %v", item.ExpectedValue, value)
		}
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		Key           string
		ExpectedValue interface{}
	}{
		{"foo", 11},
		{"bar", 22},
		{"foobar", nil},
		{"fooba", nil},
		{"foob", nil},
	}

	trie := getTestTrie()

	ok := trie.Delete("foob")
	if ok {
		t.Errorf("Deleting unexisting key must return nil")
	}

	ok = trie.Delete("foobar")
	if !ok {
		t.Errorf("Deleting existing key must not return nil")
	}

	for _, item := range cases {
		value := trie.Search(item.Key)
		if value != item.ExpectedValue {
			t.Errorf("Invalid value: expected %v, got %v", item.ExpectedValue, value)
		}
	}
}

func TestHasPrefix(t *testing.T) {
	cases := []struct {
		Key           string
		ExpectedValue map[string]interface{}
	}{
		{"f", map[string]interface{}{
			"foo":    11,
			"foobar": 111,
		}},
		{"foo", map[string]interface{}{
			"foo":    11,
			"foobar": 111,
		}},
		{"foob", map[string]interface{}{
			"foobar": 111,
		}},
		{"ba", map[string]interface{}{
			"bar": 22,
		}},
		{"xyz", map[string]interface{}{}},
	}

	trie := getTestTrie()

	for _, item := range cases {
		value := trie.HasPrefix(item.Key)
		if !reflect.DeepEqual(value, item.ExpectedValue) {
			t.Errorf("Invalid prefix values: expected %v, got %v", item.ExpectedValue, value)
		}
	}
}

func TestCounters(t *testing.T) {
	cases := []struct {
		ExpectedLen     int
		ExpectedNodeNum int
	}{
		{1, 3 + 1}, // +1 includes the root node
		{2, 6 + 1},
		{3, 9 + 1},
	}
	trie := NewTrie()
	for i, item := range testData {
		trie.Insert(item.Key, item.Value)
		if trie.Len() != cases[i].ExpectedLen {
			t.Errorf("Invalid trie length: expected %v, got %v", cases[i].ExpectedLen, trie.Len())
		}
		if trie.NodeNum() != cases[i].ExpectedNodeNum {
			t.Errorf("Invalid trie node number: expected %v, got %v", cases[i].ExpectedNodeNum, trie.NodeNum())
		}
	}
}

func getBenchTrie() *Trie {
	trie := NewTrie()
	for _, key := range benchData {
		trie.Insert(key, struct{}{})
	}
	return trie
}

func getBenchMap() map[string]interface{} {
	m := make(map[string]interface{})
	for _, key := range benchData {
		m[key] = struct{}{}
	}
	return m
}

func BenchmarkHasPrefixTrie(b *testing.B) {
	trie := getBenchTrie()

	length := len(benchData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = trie.HasPrefix(benchData[i%length])
	}
}

func BenchmarkHasPrefixMap(b *testing.B) {
	m := PrefixMap(getBenchMap())

	length := len(benchData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.HasPrefix(benchData[i%length])
	}
}

func BenchmarkInsertTrie(b *testing.B) {
	trie := NewTrie()
	length := len(benchData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.Insert(benchData[i%length], struct{}{})
	}
}

func BenchmarkInsertMap(b *testing.B) {
	m := make(map[string]interface{})
	length := len(benchData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m[benchData[i%length]] = struct{}{}
	}
}

func BenchmarkSearchTrie(b *testing.B) {
	trie := getBenchTrie()

	length := len(benchData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = trie.Search(benchData[i%length])
	}
}

func BenchmarkSearchMap(b *testing.B) {
	m := getBenchMap()

	length := len(benchData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m[benchData[i%length]]
	}
}

func BenchmarkDeleteTrie(b *testing.B) {
	trie := getBenchTrie()

	length := len(benchData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.Delete(benchData[i%length])
	}
}

func BenchmarkDeleteMap(b *testing.B) {
	m := getBenchMap()

	length := len(benchData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		delete(m, benchData[i%length])
	}
}

func BenchmarkSearchWhileInsert(b *testing.B) {
	trie := getBenchTrie()

	length := len(benchData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := benchData[i%length]
		go trie.Insert(key, struct{}{})
		trie.Search(key)
	}
}

func BenchmarkInsertWhileSearch(b *testing.B) {
	trie := getBenchTrie()

	length := len(benchData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := benchData[i%length]
		go trie.Search(key)
		trie.Insert(key, struct{}{})
	}
}

func BenchmarkParallel(b *testing.B) {
	trie := getBenchTrie()

	length := len(benchData)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := benchData[i%length]
			trie.Insert(key, struct{}{})
			trie.Search(key)
			i++
		}
	})
}
