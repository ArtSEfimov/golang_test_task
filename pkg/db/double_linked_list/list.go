package double_linked_list

type DoubleLinkedList struct {
	Head *Node
	Tail *Node
	Size uint
}

func NewDoubleLinkedList() *DoubleLinkedList {
	return &DoubleLinkedList{}
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

	dl.Size++

	if dl.Head == nil {
		dl.Head = newNode
		dl.Tail = newNode
		return
	}

	dl.Tail.Next = newNode
	newNode.Prev = dl.Tail
	dl.Tail = newNode

	OrderedMap[id] = newNode
}

func (dl *DoubleLinkedList) Remove(id uint64) {
	if dl.Size == 0 {
		return
	}

	dl.Size--

	node := OrderedMap[id]

	delete(OrderedMap, id)

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
