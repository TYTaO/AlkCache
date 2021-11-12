package lru

import "container/list" // 双向链表

// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	maxBytes int64
	nBytes   int64
	ll       *list.List
	cache    map[string]*list.Element

	OnEvicted func(key string, value Value)
}

// *list.Element
type entry struct {
	key string
	val Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get look ups a key's value
func (c *Cache) Get(key string) (Value, bool) {

	if element, ok := c.cache[key]; ok {
		c.ll.MoveToFront(element)
		e := element.Value.(*entry)
		return e.val, true
	}

	return nil, false
}

func (c *Cache) RemoveOldest() {
	if back := c.ll.Back(); back != nil {
		kv := back.Value.(*entry)
		delete(c.cache, kv.key)
		c.ll.Remove(back)
		c.nBytes -= int64(len(kv.key)) + int64(kv.val.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.val)
		}
	}
}

// Add adds a value to the cache.
func (c *Cache) Add(key string, value Value) {
	ele, ok := c.cache[key]
	if ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		kv.val = value
		c.nBytes += int64(value.Len()) - int64(kv.val.Len())
	} else {
		element := c.ll.PushFront(&entry{key: key, val: value})
		c.cache[key] = element
		c.nBytes += int64(value.Len()) + int64(len(key))
	}
	for c.maxBytes != 0 && c.nBytes > c.maxBytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}

