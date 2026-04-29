package linkedList

type DoublyLinkedList[T any] interface {
	Head() *DoubleNode[T]
	Tail() *DoubleNode[T]

	Append(data T) *DoubleNode[T]
	Push(data T) *DoubleNode[T]

	Remove(node *DoubleNode[T]) *DoubleNode[T]
	RemoveTail() *DoubleNode[T]

	MoveToHead(node *DoubleNode[T])
	Length() int
}

type DoubleNode[T any] struct {
	Data T
	Prev *DoubleNode[T]
	Next *DoubleNode[T]
}

type DoublyLinkedListImpl[T any] struct {
	Head *DoubleNode[T]
	Tail *DoubleNode[T]
}

func NewDoublyLinkedList[T any]() *DoublyLinkedListImpl[T] {
	return &DoublyLinkedListImpl[T]{}
}

func (dl *DoublyLinkedListImpl[T]) GetHead() *DoubleNode[T] {
	return dl.Head
}

func (dl *DoublyLinkedListImpl[T]) GetTail() *DoubleNode[T] {
	return dl.Tail
}

func (dl *DoublyLinkedListImpl[T]) Append(data T) *DoubleNode[T] {
	newNode := &DoubleNode[T]{Data: data, Prev: nil, Next: nil}

	if dl.Head == nil {
		dl.Head = newNode
		dl.Tail = newNode
		return newNode
	}

	newNode.Prev = dl.Tail
	dl.Tail.Next = newNode
	dl.Tail = newNode
	return newNode
}

func (dl *DoublyLinkedListImpl[T]) Push(data T) *DoubleNode[T] {
	newNode := &DoubleNode[T]{Data: data, Prev: nil, Next: nil}

	if dl.Head == nil {
		dl.Head = newNode
		dl.Tail = newNode
		return newNode
	}

	newNode.Next = dl.Head
	dl.Head.Prev = newNode
	dl.Head = newNode
	return newNode
}

func (dl *DoublyLinkedListImpl[T]) Remove(node *DoubleNode[T]) *DoubleNode[T] {

	switch {
	case dl.Head == node && dl.Tail == node:
		dl.Head = nil
		dl.Tail = nil
		return node 

	case dl.Head == node:
		dl.Head = dl.Head.Next
		dl.Head.Prev = nil
		return node

	case dl.Tail == node:
		dl.Tail = dl.Tail.Prev
		dl.Tail.Next = nil
		return node

	default:
		node.Prev.Next = node.Next
		node.Next.Prev = node.Prev
		return node
	}
}

func (dl *DoublyLinkedListImpl[T]) RemoveTail() *DoubleNode[T] {
	return dl.Remove(dl.Tail)
}

func (dl *DoublyLinkedListImpl[T]) MoveToHead(node *DoubleNode[T]) {
	if dl.Head == nil {
		return
	}

	dl.Remove(node)
	node.Next = dl.Head

	dl.Head.Prev = node
	dl.Head = node 
	dl.Head.Prev = nil
}

func (dl *DoublyLinkedListImpl[T]) Length() int {
	length := 0
	curr := dl.Head
	
	for curr != nil {
		curr = curr.Next
		length++
	} 

	return length
}
