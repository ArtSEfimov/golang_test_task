package linked_list

import (
	"go_text_task/pkg/files"
)

type DoubleLinkedList struct {
	Head *Node
	Tail *Node
	Size uint64
}

func NewDoubleLinkedList() *DoubleLinkedList {
	OrderedMap = make(map[uint64]*Node)

	//tasks := make(chan func())
	//go func(tasks chan func()) {
	//	for nextTask := range tasks {
	//		nextTask()
	//	}
	//}(tasks)

	dll := &DoubleLinkedList{}
	//dll.tasks = tasks
	if files.IsFileExists(createPath(LinkedListFileName)) {
		recoverOrderedMap(dll)
	}

	return dll
}

func (dl *DoubleLinkedList) GetSize() uint64 {
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
		storeOrderedMap(dl)
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
	dl.Size--

	defer func() {
		storeOrderedMap(dl)
	}()

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
