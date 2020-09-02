package lru

import (
	"container/list"
	"errors"
)

type ICache interface {
	Get(string) (interface{}, bool)
	Set(string, Value)
	Del(string) error
	RemoveOldest()
	Len() int64
}

type OnEvictedFunc func(string, interface{})

type cache struct {
	maxBytes  int64
	nBytes    int64
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted OnEvictedFunc
}

type Entry struct {
	Key   string
	Value Value
}

type Value interface {
	Len() int64
}

func NewCache(maxBytes int64, onEvicted OnEvictedFunc) ICache {
	return &cache{
		maxBytes:  maxBytes,
		nBytes:    0,
		ll:        new(list.List),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

//Get get value from cache
func (c *cache) Get(k string) (v interface{}, ok bool) {
	if val, ok := c.cache[k]; ok {
		c.ll.MoveToFront(val)
		return val.Value.(*Entry).Value, ok
	}
	return
}

func (c *cache) Set(k string, value Value) {
	if val, ok := c.cache[k]; ok {
		kv := val.Value.(*Entry)
		val.Value = &Entry{Key: k, Value: value}
		c.nBytes -= kv.Value.Len() + value.Len()
		c.ll.MoveToFront(val)
	} else {
		elm := c.ll.PushFront(&Entry{Key: k, Value: value})
		c.cache[k] = elm
		c.nBytes += int64(len(k)) + value.Len()
	}

	for c.nBytes > 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

func (c *cache) RemoveOldest() {
	elm := c.ll.Back()
	if elm != nil {
		kv := elm.Value.(*Entry)
		c.ll.Remove(elm)
		delete(c.cache, kv.Key)
		c.nBytes -= int64(len(kv.Key)) + kv.Value.Len()
		if c.OnEvicted != nil {
			c.OnEvicted(kv.Key, kv.Value)
		}
	}
}

func (c *cache) Del(k string) error {
	if val, ok := c.cache[k]; ok {
		kv := val.Value.(*Entry)
		delete(c.cache, kv.Key)
		c.ll.Remove(val)

		c.nBytes -= int64(len(k)) - kv.Value.Len()
		return nil
	} else {
		return errors.New("can not found k")
	}
}

func (c *cache) Len() int64 {
	return int64(c.ll.Len())
}
