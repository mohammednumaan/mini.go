package cache

import (
	"github.com/mohammednumaan/mini.go/go-cache/internal/linked_list"
	"sync"
	"time"
)

type cacheStats struct {
	Hits      int
	Misses    int
	Evictions int
}

type cacheEntry[T any] struct {
	Key    string
	Value  T
	Expiry time.Time
}

type LRUCache[T any] struct {
	capacity int
	cache    map[string]*linkedList.DoubleNode[cacheEntry[T]]
	list     *linkedList.DoublyLinkedListImpl[cacheEntry[T]]
	metrics  cacheStats
	entryTTL time.Duration
	mu       sync.Mutex
}

func NewLRUCache[T any](capacity int, ttl time.Duration) *LRUCache[T] {
	return &LRUCache[T]{
		capacity: capacity,
		cache:    make(map[string]*linkedList.DoubleNode[cacheEntry[T]]),
		list:     linkedList.NewDoublyLinkedList[cacheEntry[T]](),
		entryTTL: ttl,
	}
}

func (lru *LRUCache[T]) Get(key string) (T, bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	node, exists := lru.cache[key]
	if !exists {
		lru.metrics.Misses++
		var zero T
		return zero, false
	}

	lru.list.MoveToHead(node)
	lru.metrics.Hits++
	return node.Data.Value, true
}

func (lru *LRUCache[T]) Put(key string, value T) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	node, exists := lru.cache[key]
	if exists {
		node.Data = cacheEntry[T]{
			Key:    key,
			Value:  value,
			Expiry: time.Now().Add(lru.entryTTL),
		}

		lru.list.MoveToHead(node)
	} else {
		newData := cacheEntry[T]{
			Key:    key,
			Value:  value,
			Expiry: time.Now().Add(lru.entryTTL),
		}

		newNode := lru.list.Push(newData)
		lru.cache[key] = newNode
	}

	if len(lru.cache) > lru.capacity {
		nodeRemoved := lru.list.RemoveTail()
		if nodeRemoved != nil {
			lru.metrics.Evictions++
			delete(lru.cache, nodeRemoved.Data.Key)
		}
	}
}

func (lru *LRUCache[T]) Delete(key string) bool {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	node, exists := lru.cache[key]
	if !exists {
		return false
	}

	lru.list.Remove(node)
	delete(lru.cache, key)
	return true
}

func (lru *LRUCache[T]) Keys() []string {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	keys := make([]string, 0, len(lru.cache))
	curr := lru.list.GetHead()

	for curr != nil {
		keys = append(keys, curr.Data.Key)
		curr = curr.Next
	}

	return keys
}

func (lru *LRUCache[T]) Clear() {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	clear(lru.cache)
	lru.metrics = cacheStats{}
	lru.list = linkedList.NewDoublyLinkedList[cacheEntry[T]]()
}

func (lru *LRUCache[T]) Cleanup() {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	now := time.Now()
	curr := lru.list.GetHead()

	for curr != nil {
		if now.After(curr.Data.Expiry) {
			next := curr.Next
			lru.list.Remove(curr)
			delete(lru.cache, curr.Data.Key)
			curr = next
		} else {
			curr = curr.Next
		}
	}
}

func (lru *LRUCache[T]) Len() int {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	return len(lru.cache)
}

func (lru *LRUCache[T]) Capacity() int {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	return lru.capacity
}

func (lru *LRUCache[T]) Stats() cacheStats {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	return lru.metrics
}
