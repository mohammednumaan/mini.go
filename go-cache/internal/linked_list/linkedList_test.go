package linkedList

import "testing"

func newIntList(values ...int) (*DoublyLinkedListImpl[int], []*DoubleNode[int]) {
	dl := NewDoublyLinkedList[int]()
	nodes := make([]*DoubleNode[int], 0, len(values))

	for _, value := range values {
		nodes = append(nodes, dl.Append(value))
	}

	return dl, nodes
}

func TestDoublyLinkedList_Init(t *testing.T) {
	dl, _ := newIntList()
	if dl.GetHead() != nil {
		t.Fatalf("expected nil head")
	}
	if dl.GetTail() != nil {
		t.Fatalf("expected nil tail")
	}
	if dl.Length() != 0 {
		t.Fatalf("expected length 0, got %d", dl.Length())
	}
}

func TestDoublyLinkedList_Append(t *testing.T) {
	dl, _ := newIntList()

	first := dl.Append(100)
	if first.Data != 100 {
		t.Fatalf("expected appended node data 100, got %d", first.Data)
	}
	if first.Prev != nil {
		t.Fatalf("expected first node prev to be nil")
	}
	if first.Next != nil {
		t.Fatalf("expected first node next to be nil")
	}
	if dl.GetHead() != first {
		t.Fatalf("expected head to be first node")
	}
	if dl.GetTail() != first {
		t.Fatalf("expected tail to be first node")
	}
	if dl.Length() != 1 {
		t.Fatalf("expected length 1, got %d", dl.Length())
	}

	second := dl.Append(200)
	if second.Prev != first {
		t.Fatalf("expected second node prev to point to first node")
	}
	if second.Next != nil {
		t.Fatalf("expected second node next to be nil")
	}
	if dl.GetHead() != first {
		t.Fatalf("expected head to remain first node")
	}
	if dl.GetTail() != second {
		t.Fatalf("expected tail to be second node")
	}
	if first.Next != second {
		t.Fatalf("expected first node next to be second node")
	}
	if dl.Length() != 2 {
		t.Fatalf("expected length 2, got %d", dl.Length())
	}
}

func TestDoublyLinkedList_Push(t *testing.T) {
	dl, nodes := newIntList(20, 30)

	head := dl.Push(10)
	if head.Data != 10 {
		t.Fatalf("expected pushed node data 10, got %d", head.Data)
	}
	if head.Prev != nil {
		t.Fatalf("expected pushed head prev to be nil")
	}
	if head.Next != nodes[0] {
		t.Fatalf("expected pushed head next to point to previous head")
	}
	if dl.GetHead() != head {
		t.Fatalf("expected head to be pushed node")
	}
	if dl.GetTail() != nodes[1] {
		t.Fatalf("expected tail to remain last node")
	}
	if nodes[0].Prev != head {
		t.Fatalf("expected previous head prev to point to new head")
	}
	if dl.Length() != 3 {
		t.Fatalf("expected length 3, got %d", dl.Length())
	}
}

func TestDoublyLinkedList_Remove(t *testing.T) {
	t.Run("remove only node", func(t *testing.T) {
		dl, nodes := newIntList(10)
		removed := dl.Remove(nodes[0])

		if removed != nodes[0] {
			t.Fatalf("expected removed node to be returned")
		}
		if dl.GetHead() != nil {
			t.Fatalf("expected nil head")
		}
		if dl.GetTail() != nil {
			t.Fatalf("expected nil tail")
		}
		if dl.Length() != 0 {
			t.Fatalf("expected length 0, got %d", dl.Length())
		}
	})

	t.Run("remove head", func(t *testing.T) {
		dl, nodes := newIntList(10, 20, 30)
		removed := dl.Remove(nodes[0])

		if removed != nodes[0] {
			t.Fatalf("expected removed node to be returned")
		}
		if dl.GetHead() != nodes[1] {
			t.Fatalf("expected head to be second node")
		}
		if dl.GetHead().Prev != nil {
			t.Fatalf("expected new head prev to be nil")
		}
		if dl.GetTail() != nodes[2] {
			t.Fatalf("expected tail to remain third node")
		}
		if dl.Length() != 2 {
			t.Fatalf("expected length 2, got %d", dl.Length())
		}
	})

	t.Run("remove tail", func(t *testing.T) {
		dl, nodes := newIntList(10, 20, 30)
		removed := dl.Remove(nodes[2])

		if removed != nodes[2] {
			t.Fatalf("expected removed node to be returned")
		}
		if dl.GetHead() != nodes[0] {
			t.Fatalf("expected head to remain first node")
		}
		if dl.GetTail() != nodes[1] {
			t.Fatalf("expected tail to be second node")
		}
		if nodes[1].Next != nil {
			t.Fatalf("expected new tail next to be nil")
		}
		if dl.Length() != 2 {
			t.Fatalf("expected length 2, got %d", dl.Length())
		}
	})

	t.Run("remove middle", func(t *testing.T) {
		dl, nodes := newIntList(10, 20, 30)
		removed := dl.Remove(nodes[1])

		if removed != nodes[1] {
			t.Fatalf("expected removed node to be returned")
		}
		if dl.GetHead() != nodes[0] {
			t.Fatalf("expected head to remain first node")
		}
		if dl.GetTail() != nodes[2] {
			t.Fatalf("expected tail to remain third node")
		}
		if nodes[0].Next != nodes[2] {
			t.Fatalf("expected first node next to be third node")
		}
		if nodes[2].Prev != nodes[0] {
			t.Fatalf("expected third node prev to be first node")
		}
		if dl.Length() != 2 {
			t.Fatalf("expected length 2, got %d", dl.Length())
		}
	})
}

func TestDoublyLinkedList_RemoveTail(t *testing.T) {
	dl, nodes := newIntList(10, 20, 30)
	removed := dl.RemoveTail()

	if removed != nodes[2] {
		t.Fatalf("expected removed tail to be returned")
	}
	if dl.GetTail() != nodes[1] {
		t.Fatalf("expected tail to be second node")
	}
	if nodes[1].Next != nil {
		t.Fatalf("expected new tail next to be nil")
	}
	if dl.Length() != 2 {
		t.Fatalf("expected length 2, got %d", dl.Length())
	}
}

func TestDoublyLinkedList_MoveToHead(t *testing.T) {
	t.Run("move tail to head", func(t *testing.T) {
		dl, nodes := newIntList(10, 20, 30)
		dl.MoveToHead(nodes[2])
		if dl.GetHead() != nodes[2] {
			t.Fatalf("expected head to be previous tail")
		}
		if nodes[2].Prev != nil {
			t.Fatalf("expected new head prev to be nil")
		}
		if nodes[2].Next != nodes[0] {
			t.Fatalf("expected new head next to be first node")
		}
		if nodes[0].Next != nodes[1] {
			t.Fatalf("expected first node next to be second node")
		}
		if dl.GetTail() != nodes[1] {
			t.Fatalf("expected tail to be previous middle node")
		}
	})

	t.Run("move middle to head", func(t *testing.T) {
		dl, nodes := newIntList(10, 20, 30)
		dl.MoveToHead(nodes[1])
		if dl.GetHead() != nodes[1] {
			t.Fatalf("expected head to be previous middle node")
		}
		if nodes[1].Prev != nil {
			t.Fatalf("expected new head prev to be nil")
		}
		if nodes[1].Next != nodes[0] {
			t.Fatalf("expected new head next to be first node")
		}
		if nodes[0].Next != nodes[2] {
			t.Fatalf("expected first node next to be third node")
		}
		if dl.GetTail() != nodes[2] {
			t.Fatalf("expected tail to remain third node")
		}
	})
}

func TestDoublyLinkedList_Length(t *testing.T) {
	dl, _ := newIntList()

	if dl.Length() != 0 {
		t.Fatalf("expected length 0, got %d", dl.Length())
	}

	dl.Append(10)
	dl.Append(20)
	dl.Push(5)

	if dl.Length() != 3 {
		t.Fatalf("expected length 3, got %d", dl.Length())
	}

	dl.RemoveTail()

	if dl.Length() != 2 {
		t.Fatalf("expected length 2, got %d", dl.Length())
	}
}
