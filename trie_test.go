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

// PrefixMap searches for prefixes in a map,
// ignoring thread-safety for simplicity
type PrefixMap map[string]interface{}

func (m PrefixMap) WithPrefix(prefix string) map[string]interface{} {
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

func getTestTrie() *Trie {
	trie := NewTrie()
	for _, item := range testData {
		trie.Insert(item.Key, item.Value)
	}
	return trie
}

func getTestPrefixMap() PrefixMap {
	m := make(PrefixMap)
	for _, item := range testData {
		m[item.Key] = item.Value
	}
	return m
}

func TestMapWithPrefix(t *testing.T) {
	tests := []struct {
		key  string
		want map[string]interface{}
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

	m := getTestPrefixMap()

	for _, test := range tests {
		value := m.WithPrefix(test.key)
		if !reflect.DeepEqual(value, test.want) {
			t.Errorf("Invalid Map prefix values: expected %v, got %v", test.want, value)
		}
	}
}

func TestNewTrie(t *testing.T) {
	var trie interface{} = NewTrie()
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
	got := trie.WithPrefix("")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WithPrefix must include root node data if available")
	}
}

func TestInsertAndSearch(t *testing.T) {
	tests := []struct {
		key  string
		want interface{}
	}{
		{"foo", 11},
		{"foobar", 111},
		{"bar", 22},
		{"foob", nil},
		{"foobarr", nil},
	}

	trie := getTestTrie()

	for _, test := range tests {
		value := trie.Search(test.key)
		if value != test.want {
			t.Errorf("Invalid value: expected %v, got %v", test.want, value)
		}
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		key  string
		want interface{}
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

	for _, test := range tests {
		value := trie.Search(test.key)
		if value != test.want {
			t.Errorf("Invalid value: expected %v, got %v", test.want, value)
		}
	}
}

func TestWithPrefix(t *testing.T) {
	tests := []struct {
		key  string
		want map[string]interface{}
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

	for _, test := range tests {
		value := trie.WithPrefix(test.key)
		if !reflect.DeepEqual(value, test.want) {
			t.Errorf("Invalid prefix values: expected %v, got %v", test.want, value)
		}
	}
}

func TestCounters(t *testing.T) {
	tests := []struct {
		wantLen     int
		wantNodeNum int
	}{
		{1, 3 + 1}, // +1 includes the root node
		{2, 6 + 1},
		{3, 9 + 1},
	}
	trie := NewTrie()
	for i, item := range testData {
		trie.Insert(item.Key, item.Value)
		if trie.Len() != tests[i].wantLen {
			t.Errorf("Invalid trie length: expected %v, got %v", tests[i].wantLen, trie.Len())
		}
		if trie.NodeNum() != tests[i].wantNodeNum {
			t.Errorf("Invalid trie node number: expected %v, got %v", tests[i].wantNodeNum, trie.NodeNum())
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

func getBenchPrefixMap() PrefixMap {
	m := make(PrefixMap)
	for _, key := range benchData {
		m[key] = struct{}{}
	}
	return m
}

func BenchmarkWithPrefixTrie(b *testing.B) {
	b.ReportAllocs()
	trie := getBenchTrie()

	length := len(benchData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = trie.WithPrefix(benchData[i%length])
	}
}

func BenchmarkWithPrefixMap(b *testing.B) {
	b.ReportAllocs()
	m := getBenchPrefixMap()

	length := len(benchData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.WithPrefix(benchData[i%length])
	}
}

func BenchmarkInsertTrie(b *testing.B) {
	b.ReportAllocs()
	trie := NewTrie()
	length := len(benchData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.Insert(benchData[i%length], struct{}{})
	}
}

func BenchmarkInsertMap(b *testing.B) {
	b.ReportAllocs()
	m := make(map[string]interface{})
	length := len(benchData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m[benchData[i%length]] = struct{}{}
	}
}

func BenchmarkSearchTrie(b *testing.B) {
	b.ReportAllocs()
	trie := getBenchTrie()

	length := len(benchData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = trie.Search(benchData[i%length])
	}
}

func BenchmarkSearchMap(b *testing.B) {
	b.ReportAllocs()
	m := getBenchPrefixMap()

	length := len(benchData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m[benchData[i%length]]
	}
}

func BenchmarkDeleteTrie(b *testing.B) {
	b.ReportAllocs()
	trie := getBenchTrie()

	length := len(benchData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.Delete(benchData[i%length])
	}
}

func BenchmarkDeleteMap(b *testing.B) {
	b.ReportAllocs()
	m := getBenchPrefixMap()

	length := len(benchData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		delete(m, benchData[i%length])
	}
}

func BenchmarkSearchWhileInsert(b *testing.B) {
	b.ReportAllocs()
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
	b.ReportAllocs()
	trie := getBenchTrie()

	length := len(benchData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := benchData[i%length]
		go trie.Search(key)
		trie.Insert(key, struct{}{})
	}
}

func BenchmarkSearchWhileInsertParallel(b *testing.B) {
	b.ReportAllocs()
	trie := getBenchTrie()

	length := len(benchData)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := benchData[i%length]
			go trie.Insert(key, struct{}{})
			trie.Search(key)
			i++
		}
	})
}

func BenchmarkInsertWhileSearchParallel(b *testing.B) {
	b.ReportAllocs()
	trie := getBenchTrie()

	length := len(benchData)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := benchData[i%length]
			go trie.Search(key)
			trie.Insert(key, struct{}{})
			i++
		}
	})
}
