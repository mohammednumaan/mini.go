package cache

import "time"

type cacheItem struct {
	Value string
	Expiry time.Time
}

type cache struct {
	items map[string]*cacheItem
}

func (c *cache) Set(key, value string) *cacheItem {
	c.items[key] = &cacheItem{ Value: value, Expiry: time.Now().Add(5 * time.Second) }
	return c.items[key];
}

func (c *cache) Get(key string) *cacheItem {
	value, exists := c.items[key]
	if !exists {
		return nil
	}
	return value
}

func (c *cache) Delete(key string) bool {
	delete(c.items, key)
	return true
}

func (c *cache) CleanUp() {
	// this is a simple TTL-based 
	// key expiration
	now := time.Now()
	for key, item := range c.items {
		keyTime := item.Expiry
		if now.After(keyTime) {
			c.Delete(key)
		}
	}
}

func NewCache() *cache {
	c := cache{ items: make(map[string]*cacheItem), }
	return &c;
}
