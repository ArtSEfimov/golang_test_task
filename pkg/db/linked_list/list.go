package linked_list

import (
	"go_text_task/pkg/files"
	"sync"
)

type DoubleLinkedList struct {
	Head  *Node
	Tail  *Node
	Size  uint
	mutex *sync.Mutex
}

func NewDoubleLinkedList() *DoubleLinkedList {
	OrderedMap = make(map[uint64]*Node)
	dll := &DoubleLinkedList{}
	dll.mutex = new(sync.Mutex)
	if files.IsFileExists(createPath(LinkedListFileName)) {
		recoverOrderedMap(dll)
	}

	return dll
}

func (dl *DoubleLinkedList) GetSize() uint {
	return dl.Size
}

func (dl *DoubleLinkedList) GetHead() *Node {
	return dl.Head
}

func (dl *DoubleLinkedList) GetTail() *Node {
	return dl.Tail
}

func (dl *DoubleLinkedList) Append(id uint64) {
	newNode := &Node{
		Value: id,
		Next:  nil,
		Prev:  nil,
	}
	defer func() {
		OrderedMap[id] = newNode
		dl.Size++
		go storeOrderedMap(dl)
	}()

	if dl.Head == nil {
		dl.Head = newNode
		dl.Tail = newNode
		return
	}

	dl.Tail.Next = newNode
	newNode.Prev = dl.Tail
	dl.Tail = newNode

}

func (dl *DoubleLinkedList) Remove(id uint64) {
	if dl.Size == 0 {
		return
	}

	node := OrderedMap[id]

	delete(OrderedMap, id)
	go storeOrderedMap(dl)

	defer func() { dl.Size-- }()

	if node == dl.Head {
		if node == dl.Tail {
			dl.Head = nil
			dl.Tail = nil
			return
		}
		next := node.Next
		node.Next = nil
		next.Prev = nil
		dl.Head = next
		return
	}
	if node == dl.Tail {
		prev := node.Prev
		node.Prev = nil
		prev.Next = nil
		dl.Tail = prev
		return
	}
	prev := node.Prev
	next := node.Next
	next.Prev = prev
	prev.Next = next
}
