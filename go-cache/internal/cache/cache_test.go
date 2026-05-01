package cache

import (
	"testing"
	"time"
)

const (
	a = "first"
	b = "second"
	c = "third"
	d = "fourth"
	e = "fifth"
	f = "sixth"
)

func fillCache() *LRUCache[string] {
	cache := NewLRUCache[string](5, time.Hour)
	cache.Put("a", a)
	cache.Put("b", b)
	cache.Put("c", c)
	cache.Put("d", d)
	cache.Put("e", e)
	return cache
}

func containsInOrder(keys []string, expected []string) bool {
	if len(keys) != len(expected) {
		return false
	}

	expectedIdx := 0
	for _, key := range expected {
		if key == expected[expectedIdx] {
			expectedIdx++
		} else {
			return false
		}
	}

	return expectedIdx == len(expected)
}

func TestCache_Init(t *testing.T) {
	cache := NewLRUCache[string](5, time.Hour)

	if cache == nil {
		t.Error("Expected cache to be initialized, got nil")
	}
	if cache.Len() != 0 {
		t.Errorf("Expected cache length to be 0, got %d", cache.Len())
	}
	if cache.Capacity() != 5 {
		t.Errorf("Expected capacity to be 5, got %d", cache.capacity)
	}
	if cache.list.GetHead() != nil {
		t.Errorf("Expected head to be nil, got %v", cache.list.GetHead())
	}
	if cache.list.GetTail() != nil {
		t.Errorf("Expected tail to be nil, got %v", cache.list.GetTail())
	}
	if len(cache.cache) != 0 {
		t.Errorf("Expected cache map to be empty, got %d items", len(cache.cache))
	}
}

func TestCache_Keys(t *testing.T) {
	cache := fillCache()
	keys := cache.Keys()
	expectedKeys := []string{"e", "d", "c", "b", "a"}

	if !containsInOrder(keys, expectedKeys) {
		t.Errorf("Expected keys to be in order %v, got %v", expectedKeys, keys)
	}
}

func TestCache_Put(t *testing.T) {
	t.Run("should add items to cache and maintain order", func(t *testing.T) {
		cache := NewLRUCache[string](5, time.Hour)
		cache.Put(a, a)
		cache.Put(b, b)

		if cache.Len() != 2 {
			t.Errorf("Expected cache length to be 2, got %d", cache.Len())
		}

		keys := cache.Keys()
		expectedKeys := []string{"b", "a"}

		if !containsInOrder(keys, expectedKeys) {
			t.Errorf("Expected keys to be in order %v, got %v", expectedKeys, keys)
		}
	})

	t.Run("should not exceed capacity when cache is at capacity", func(t *testing.T) {

		cache := fillCache()
		if cache.Len() != 5 {
			t.Errorf("Expected cache length to be 5, got %d", cache.Len())
		}

		cache.Put(f, f)
		if cache.Len() != 5 {
			t.Errorf("Expected cache length to remain 5 after adding item beyond capacity, got %d", cache.Len())
		}
	})

}

func TestCache_Get(t *testing.T) {
	t.Run("should return value and true for existing key", func(t *testing.T) {
		cache := fillCache()
		value, exists := cache.Get("c")

		if !exists {
			t.Errorf("Expected key 'c' to exist, got false")
		}
		if value != c {
			t.Errorf("Expected value for key 'c' to be '%s', got '%s'", c, value)
		}
	})

	t.Run("should return zero value and false for non-existing key", func(t *testing.T) {
		cache := fillCache()
		value, exists := cache.Get("nonexistent")

		if exists {
			t.Errorf("Expected key 'nonexistent' to not exist, got true")
		}
		if value != "" {
			t.Errorf("Expected zero value for non-existing key, got '%s'", value)
		}
	})

	t.Run("should move accessed item to head of the list", func(t *testing.T) {
		cache := fillCache()
		cache.Get("c")

		keys := cache.Keys()
		expectedKeys := []string{"c", "e", "d", "b", "a"}

		if !containsInOrder(keys, expectedKeys) {
			t.Errorf("Expected keys to be in order %v after accessing 'c', got %v", expectedKeys, keys)
		}
	})

	t.Run("should get items properly after capacity", func(t *testing.T) {
		cache := fillCache()
		cache.Put(f, f)

		keys := cache.Keys()
		expected := []string{f, e, d, c, b}
		if !containsInOrder(keys, expected) {
			t.Errorf("Expected keys to be in order %v after putting 'f', got %v", expected, keys)
		}

		value, exists := cache.Get(a)
		if exists {
			t.Errorf("Expected 'a' to be evicted after 'f' was put, got '%s'", value)
		}
	})
}

func TestCache_Delete(t *testing.T) {
	t.Run("should return true and remove existing key", func(t *testing.T) {
		cache := fillCache()

		deleted := cache.Delete("c")
		if !deleted {
			t.Errorf("Expected delete to return true for existing key 'c'")
		}
		if cache.Len() != 4 {
			t.Errorf("Expected cache length to be 4 after delete, got %d", cache.Len())
		}

		keys := cache.Keys()
		expectedKeys := []string{"e", "d", "b", "a"}
		if !containsInOrder(keys, expectedKeys) {
			t.Errorf("Expected keys to be in order %v, got %v", expectedKeys, keys)
		}
	})

	t.Run("should return false for non-existing key", func(t *testing.T) {
		cache := fillCache()

		deleted := cache.Delete("nonexistent")
		if deleted {
			t.Errorf("Expected delete to return false for non-existing key")
		}
		if cache.Len() != 5 {
			t.Errorf("Expected cache length to remain 5, got %d", cache.Len())
		}
	})

	t.Run("should update order after delete", func(t *testing.T) {
		cache := fillCache()
		cache.Delete("e")

		keys := cache.Keys()
		expectedKeys := []string{"d", "c", "b", "a"}
		if !containsInOrder(keys, expectedKeys) {
			t.Errorf("Expected keys to be in order %v after delete, got %v", expectedKeys, keys)
		}
	})
}

func TestCache_Clear(t *testing.T) {
	t.Run("should remove all items from cache", func(t *testing.T) {
		cache := fillCache()

		cache.Clear()
		if cache.Len() != 0 {
			t.Errorf("Expected cache length to be 0 after clear, got %d", cache.Len())
		}
		if len(cache.Keys()) != 0 {
			t.Errorf("Expected keys to be empty after clear, got %v", cache.Keys())
		}
	})

	t.Run("should handle clearing empty cache", func(t *testing.T) {
		cache := NewLRUCache[string](5, time.Hour)

		cache.Clear()
		if cache.Len() != 0 {
			t.Errorf("Expected cache length to be 0 after clear, got %d", cache.Len())
		}
	})
}

func TestCache_Cleanup(t *testing.T) {
	t.Run("should remove expired entries", func(t *testing.T) {
		cache := NewLRUCache[string](5, 50*time.Millisecond)
		cache.Put(a, a)
		cache.Put(b, b)

		time.Sleep(100 * time.Millisecond)

		cache.Cleanup()
		if cache.Len() != 0 {
			t.Errorf("Expected cache length to be 0 after cleanup, got %d", cache.Len())
		}
	})

	t.Run("should keep non-expired entries when called before expiry", func(t *testing.T) {
		cache := NewLRUCache[string](5, 10*time.Second)
		cache.Put(a, a)
		cache.Put(b, b)

		cache.Cleanup()
		if cache.Len() != 2 {
			t.Errorf("Expected cache length to be 2, got %d", cache.Len())
		}

		keys := cache.Keys()
		expectedKeys := []string{b, a}
		if !containsInOrder(keys, expectedKeys) {
			t.Errorf("Expected keys to be in order %v, got %v", expectedKeys, keys)
		}

		value, exists := cache.Get(a)
		if !exists || value != a {
			t.Errorf("Expected to get value '%s' for key 'a', got '%s', exists: %v", a, value, exists)
		}
	})

	t.Run("should remove only expired entries", func(t *testing.T) {
		cache := NewLRUCache[string](5, 50*time.Millisecond)
		cache.Put(a, a)
		cache.Put(b, b)

		time.Sleep(100 * time.Millisecond)

		cache.Cleanup()
		cache.Put(c, c)

		keys := cache.Keys()
		expectedKeys := []string{c}
		if !containsInOrder(keys, expectedKeys) {
			t.Errorf("Expected keys to be only 'c' after cleanup, got %v", keys)
		}
		if cache.Len() != 1 {
			t.Errorf("Expected cache length to be 1, got %d", cache.Len())
		}
	})

	t.Run("should preserve non-expired entries after repeated cleanup calls", func(t *testing.T) {
		cache := NewLRUCache[string](5, 10*time.Second)
		cache.Put(a, a)
		cache.Put(b, b)
		cache.Put(c, c)

		cache.Cleanup()
		cache.Cleanup()
		cache.Cleanup()

		if cache.Len() != 3 {
			t.Errorf("Expected cache length to be 3, got %d", cache.Len())
		}

		keys := cache.Keys()
		expectedKeys := []string{c, b, a}
		if !containsInOrder(keys, expectedKeys) {
			t.Errorf("Expected keys to be in order %v, got %v", expectedKeys, keys)
		}
	})
}
