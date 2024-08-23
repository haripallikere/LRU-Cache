package internal

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type CacheItem struct {
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
	ExpiresAt time.Time   `json:"expiresAt"`
}

type LRUCache struct {
	capacity int
	items    map[string]*list.Element
	list     *list.List
	mutex    sync.RWMutex
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if element, exists := c.items[key]; exists {
		item := element.Value.(*CacheItem)
		if time.Now().After(item.ExpiresAt) {
			c.removeElement(element)
			return nil, false
		}
		c.list.MoveToFront(element)
		return item.Value, true
	}
	return nil, false
}

func (c *LRUCache) Set(key string, value interface{}, expiration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, exists := c.items[key]; exists {
		c.list.MoveToFront(element)
		item := element.Value.(*CacheItem)
		item.Value = value
		item.ExpiresAt = time.Now().Add(expiration)
	} else {
		if c.list.Len() >= c.capacity {
			c.removeOldest()
		}
		item := &CacheItem{
			Key:       key,
			Value:     value,
			ExpiresAt: time.Now().Add(expiration),
		}
		fmt.Println("CacheItem CacheItem:", item)
		element := c.list.PushFront(item)
		c.items[key] = element
	}
}

func (c *LRUCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, exists := c.items[key]; exists {
		c.removeElement(element)
	} else {
		fmt.Println("Key not found in cache:", key)
	}
}

func (c *LRUCache) GetAll() []CacheItem {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	result := make([]CacheItem, 0, len(c.items))
	for _, element := range c.items {
		item := element.Value.(*CacheItem)
		if time.Now().Before(item.ExpiresAt) {
			result = append(result, *item)
		}
	}
	fmt.Println("GetAll:", result)
	return result
}

func (c *LRUCache) removeElement(element *list.Element) {
	c.list.Remove(element)
	item := element.Value.(*CacheItem)
	delete(c.items, item.Key)
}

func (c *LRUCache) removeOldest() {
	element := c.list.Back()
	if element != nil {
		c.removeElement(element)
	}
}

func (c *LRUCache) StartCleanupTask() {
	go func() {
		for {
			time.Sleep(time.Second)
			c.cleanup()
		}
	}()
}

func (c *LRUCache) cleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for key, element := range c.items {
		item := element.Value.(*CacheItem)
		if time.Now().After(item.ExpiresAt) {
			c.removeElement(element)
			delete(c.items, key)
		}
	}
}
