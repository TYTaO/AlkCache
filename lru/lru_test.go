package lru

import (
	"reflect"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestGet(t *testing.T) {
	lru := New(0, nil)

	lru.Add("key1", String("1234"))
	if ele, ok := lru.Get("key1"); !ok || string(ele.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}

	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestRemoveoldest(t *testing.T) {
	lru := New(10, nil)
	lru.Add("key1", String("1234"))
	lru.Add("key2", String("1234"))
	lru.Add("key3", String("1234"))

	if ele, ok := lru.Get("key3"); !ok || string(ele.(String)) != "1234" {
		t.Fatalf("cache get key3 failed")
	}

	if _, ok := lru.Get("key1"); ok {
		t.Fatalf("cache remove oldest key1 failed")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	var callback func(string, Value)
	callback = func(s string, value Value) {
		keys = append(keys, s)
	}
	lru := New(9, callback)
	lru.Add("key1", String("1234"))
	lru.Add("key2", String("1234"))
	lru.Add("key3", String("1234"))

	expect := []string{"key1", "key2"}
	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s, now keys: %s", expect, keys)
	}
}
