package cache

import "testing"

func TestCache_Init(t *testing.T) {
	cache := NewLRUCache[string](5)

	if cache == nil {
		t.Error("Expected cache to be initialized, got nil")
	}

	if cache.capacity != 5 {
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

func TestCache_Put(t *testing.T) {
	cache := NewLRUCache[string](2)

	cache.Put("a", "1")
	cache.Put("b", "2")

	if len(cache.cache) != 2 {
		t.Errorf("Expected cache map to have 2 items, got %d", len(cache.cache))
	}

	if cache.list.GetHead().Data.Key != "b" {
		t.Errorf("Expected head key to be 'b', got '%s'", cache.list.GetHead().Data.Key)
	}

	if cache.list.GetTail().Data.Key != "a" {
		t.Errorf("Expected tail key to be 'a', got '%s'", cache.list.GetTail().Data.Key)
	}

	cache.Put("c", "3")

	if len(cache.cache) != 2 {
		t.Errorf("Expected cache map to have 2 items after eviction, got %d", len(cache.cache))
	}

	if cache.list.GetHead().Data.Key != "c" {
		t.Errorf("Expected head key to be 'c', got '%s'", cache.list.GetHead().Data.Key)
	}

	if cache.list.GetTail().Data.Key != "b" {
		t.Errorf("Expected tail key to be 'b', got '%s'", cache.list.GetTail().Data.Key)
	}
}


func TestCache_Get(t *testing.T) {
	cache := NewLRUCache[string](2)

	cache.Put("a", "1")
	cache.Put("b", "2")
	
	value, found := cache.Get("a")
	if !found {
		t.Error("Expected to find key 'a', but it was not found")
	}

	if value != "1" {
		t.Errorf("Expected value for key 'a' to be '1', got '%s'", value)
	}

	if cache.list.GetHead().Data.Key != "a" {
		t.Errorf("Expected head key to be 'a' after access, got '%s'", cache.list.GetHead().Data.Key)
	}

	cache.Put("c", "3")

	if _, found := cache.Get("b"); found {
		t.Error("Expected key 'b' to be evicted, but it was found")
	}

	if cache.list.GetHead().Data.Key != "c" {
		t.Errorf("Expected head key to be 'c', got '%s'", cache.list.GetHead().Data.Key)
	}

	if cache.list.GetTail().Data.Key != "a" {
		t.Errorf("Expected tail key to be 'a', got '%s'", cache.list.GetTail().Data.Key)
	}
}


func TestCache_Delete(t *testing.T) {
	cache := NewLRUCache[string](2)

	cache.Put("a", "1")
	cache.Put("b", "2")

	cache.Delete("a")

	if _, found := cache.Get("a"); found {
		t.Error("Expected key 'a' to be deleted, but it was found")
	}

	if cache.list.GetHead().Data.Key != "b" {
		t.Errorf("Expected head key to be 'b' after deletion, got '%s'", cache.list.GetHead().Data.Key)
	}

	if cache.list.GetTail().Data.Key != "b" {
		t.Errorf("Expected tail key to be 'b' after deletion, got '%s'", cache.list.GetTail().Data.Key)
	}
}

func TestCache_Clear(t *testing.T) {
	cache := NewLRUCache[string](2)

	cache.Put("a", "1")
	cache.Put("b", "2")

	cache.Clear()

	if len(cache.cache) != 0 {
		t.Errorf("Expected cache map to be empty after clear, got %d items", len(cache.cache))
	}

	if cache.list.GetHead() != nil {
		t.Errorf("Expected head to be nil after clear, got %v", cache.list.GetHead())
	}

	if cache.list.GetTail() != nil {
		t.Errorf("Expected tail to be nil after clear, got %v", cache.list.GetTail())
	}
}
