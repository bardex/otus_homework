package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity   int
	queue      List
	emptyQueue List
	mx         sync.Mutex
	items      map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mx.Lock()
	defer c.mx.Unlock()

	val, exists := c.items[key]
	if !exists {
		return nil, false
	}
	c.queue.MoveToFront(val)
	if v, ok := val.Value.(cacheItem); ok {
		return v.value, true
	}
	return val.Value, true
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mx.Lock()
	defer c.mx.Unlock()

	val, exists := c.items[key]
	if exists {
		val.Value = cacheItem{key: key, value: value}
		c.queue.MoveToFront(val)
		return exists
	}

	if c.queue.Len() >= c.capacity {
		last := c.queue.Back()
		if v, ok := last.Value.(cacheItem); ok {
			delete(c.items, v.key)
		}
		c.queue.Remove(last)
	}
	c.items[key] = c.queue.PushFront(cacheItem{key: key, value: value})
	return exists
}

func (c *lruCache) Clear() {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue = c.emptyQueue
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity:   capacity,
		queue:      NewList(),
		emptyQueue: NewList(),
		items:      make(map[Key]*ListItem, capacity),
	}
}
