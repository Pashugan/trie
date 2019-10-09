package trie

import (
	"reflect"
	"testing"
)

var initData = []struct {
	Key   string
	Value interface{}
}{
	{"foo", 11},
	{"foobar", 111},
	{"bar", 22},
}

func TestNewTrie(t *testing.T) {
	var trie interface{}

	trie = NewTrie()
	_, ok := trie.(*Node)
	if !ok {
		t.Errorf("Invalid trie type")
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
		{"", nil},
		{"foob", nil},
		{"foobarr", nil},
	}

	trie := NewTrie()
	for _, item := range initData {
		trie.Insert(item.Key, item.Value)
	}

	for _, item := range cases {
		value, _ := trie.Search(item.Key)
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

	trie := NewTrie()
	for _, item := range initData {
		trie.Insert(item.Key, item.Value)
	}

	ok := trie.Delete("foob")
	if ok {
		t.Errorf("Deleting unexisting key must return nil")
	}

	ok = trie.Delete("foobar")
	if !ok {
		t.Errorf("Deleting existing key must not return nil")
	}

	for _, item := range cases {
		value, _ := trie.Search(item.Key)
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

	trie := NewTrie()
	for _, item := range initData {
		trie.Insert(item.Key, item.Value)
	}

	for _, item := range cases {
		value := trie.HasPrefix(item.Key)
		if !reflect.DeepEqual(value, item.ExpectedValue) {
			t.Errorf("Invalid prefix values: expected %v, got %v", item.ExpectedValue, value)
		}
	}
}
