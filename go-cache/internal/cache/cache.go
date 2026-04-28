package cache

type cacheItem struct {
	Value string
}

type cache struct {
	items map[string]*cacheItem
}

func (c *cache) Set(key, value string) *cacheItem {
	c.items[key] = &cacheItem{ Value: value }
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

func NewCache() *cache {
	c := cache{ items: make(map[string]*cacheItem), }
	return &c;
}
