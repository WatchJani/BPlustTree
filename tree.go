package BPTree

import (
	"fmt"
)

type KeyType interface {
	int | string | float64 | float32 | int16 | int8 | int32 | int64
}

type Tree[K KeyType, V any] struct {
	root   *Node[K, V]
	degree int
}

type item[K KeyType, V any] struct {
	key   K
	value V
}

type Node[K KeyType, V any] struct {
	items     []item[K, V]
	Children  []*Node[K, V]
	nextNodeL *Node[K, V]
	nextNodeR *Node[K, V]
	pointer   int
}

type Stack[T any] struct {
	store []T
}

type positionStr[K KeyType, V any] struct {
	node     *Node[K, V]
	position int
}

func newItem[K KeyType, V any](key K, value V) item[K, V] {
	return item[K, V]{
		key:   key,
		value: value,
	}
}

func newNode[K KeyType, V any](degree int) Node[K, V] {
	return Node[K, V]{
		items:    make([]item[K, V], degree+1),
		Children: make([]*Node[K, V], degree+2),
	}
}

func (n *Node[K, V]) delete() {
	n.items = nil
	n.Children = nil
	n.pointer = 0
	n.nextNodeL = nil
	n.nextNodeR = nil
}

func New[K KeyType, V any](degree int) *Tree[K, V] {
	if degree < 3 {
		degree = 3
	}

	return &Tree[K, V]{
		degree: degree,
		root: &Node[K, V]{
			items:    make([]item[K, V], degree+1),
			Children: make([]*Node[K, V], degree+2),
		},
	}
}

func (n *Node[K, V]) search(target K) (int, bool) {
	low, high := 0, n.pointer-1

	for low <= high {
		mid := (low + high) / 2

		if n.items[mid].key == target {
			return mid + 1, true
		} else if n.items[mid].key < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return low, false
}

func newStack[T any](capacity int) Stack[T] {
	return Stack[T]{
		store: make([]T, 0, capacity),
	}
}

// func (s *Stack[K, V]) Push(node *Node[K, V], position int) {
// 	s.store = append(s.store, positionStr[K, V]{
// 		node:     node,
// 		position: position,
// 	})
// }

func (s *Stack[T]) Push(value T) {
	s.store = append(s.store, value)
}

func (s *Stack[T]) Pop() (T, error) {
	if len(s.store) == 0 {
		var zeroValue T
		return zeroValue, fmt.Errorf("stack is empty")
	}
	value := s.store[len(s.store)-1]
	s.store = s.store[:len(s.store)-1]
	return value, nil
}

func (t *Tree[K, V]) Find(key K) (V, error) {
	for next := t.root; next != nil; {
		index, found := next.search(key)

		if found {
			return next.items[index-1].value, nil
		}

		next = next.Children[index]
	}

	var res V
	return res, fmt.Errorf("key %v not found", key)
}

func (t *Tree[K, V]) Insert(key K, value V) {
	stack, item := newStack[positionStr[K, V]](4), newItem(key, value)
	position, found := findLeaf(t.root, &stack, key)

	current, _ := stack.Pop()
	//update just state state
	if found {
		current.node.items[position].value = value
		return
	}

	if middleKey, nodeChildren := insertLeaf(current.node, position, t.degree, item); nodeChildren != nil {
		for {
			if nextCurrent, err := stack.Pop(); err != nil {
				break
			} else {
				current = nextCurrent
			}

			current.node.pointer += insert(current.node.items, middleKey, current.position)
			chIndex := childrenIndex(middleKey.key, current.node.items[current.position].key, current.position)

			insert(current.node.Children, nodeChildren, chIndex)
			if current.node.pointer < t.degree {
				return
			}

			middle := current.node.pointer / 2
			middleKey = current.node.items[middle]

			newNode := newNode[K, V](t.degree)
			newNode.pointer += migrate(newNode.items, current.node.items[:middle], 0) // migrate half elements to new node
			migrate(newNode.Children, current.node.Children[:middle+1], 0)

			current.node.pointer -= deleteElement(current.node.items, 0, current.node.pointer-middle+1-t.degree&1)
			migrate(current.node.Children, current.node.Children[middle+1:], 0)
			nodeChildren = &newNode
		}

		rootNode := newNode[K, V](t.degree)
		rootNode.pointer += insert(rootNode.items, middleKey, 0)
		rootNode.Children[0] = nodeChildren
		rootNode.Children[1] = current.node
		t.root = &rootNode
	}
}

func childrenIndex[K KeyType](key, value K, index int) int {
	if value < key {
		return index + 1
	}

	return index
}

func insert[T any](list []T, insert T, position int) int {
	copy(list[position+1:], list[position:])
	return copy(list[position:], []T{insert})
}

func migrate[T any](list, migrateElement []T, position int) int {
	return copy(list[position:], migrateElement)
}

func deleteElement[T any](list []T, position, deletion int) int {
	copy(list[position:], list[position+deletion:])
	return deletion
}

func insertLeaf[K KeyType, V any](current *Node[K, V], position, degree int, item item[K, V]) (item[K, V], *Node[K, V]) {
	current.pointer += insert(current.items, item, position)

	if current.pointer < degree {
		return item, nil
	}

	newNode := newNode[K, V](degree)
	middle := degree / 2

	newNode.pointer += migrate(newNode.items, current.items[:middle], 0)
	current.pointer -= deleteElement(current.items, 0, current.pointer-middle-degree%2)

	//update links between leafs
	if current.nextNodeL != nil {
		newNode.nextNodeL = current.nextNodeL
		current.nextNodeL.nextNodeR = &newNode
	}

	current.nextNodeL = &newNode
	newNode.nextNodeR = current

	return current.items[0], &newNode
}

func minAllowed(degree, numElement int) bool {
	return (degree-1)/2 <= numElement
}

func findLeaf[K KeyType, V any](root *Node[K, V], stack *Stack[positionStr[K, V]], key K) (int, bool) {
	position, found := 0, false

	for current := root; current != nil; {
		position, found = current.search(key)
		stack.Push(positionStr[K, V]{
			current,
			position,
		})

		current = current.Children[position]
	}

	return position, found
}

func (t *Tree[K, V]) Delete(key K) error {
	stack := newStack[positionStr[K, V]](4)
	_, found := findLeaf(t.root, &stack, key)

	if !found {
		return fmt.Errorf("key %v does not exist", key)
	}

	current, _ := stack.Pop()

	for {
		current.node.pointer -= deleteElement(current.node.items, indexElement(current.position), 1)

		// Check if the current node has fallen below the minimum allowed size
		if minAllowed(t.degree, current.node.pointer) || (found && t.root.pointer == 0) {
			return nil
		}

		// Attempt to pop the parent node off the stack
		if nextCurrent, err := stack.Pop(); err != nil {
			break
		} else {
			parent := nextCurrent
			sibling, side, operation := sibling(parent, t.degree)

			// Perform either a transfer or a merge based on the sibling node
			if operation {
				transfer(parent, current, sibling, found, side)
				return nil
			} else {
				merge(current.node, sibling, parent, found, side)
			}

			// Update the current node to the parent node
			current = parent

			// Once found, we don't want to toggle it back
			if found {
				found = !found
			}
		}
	}

	if t.root.pointer == 0 && len(t.root.Children) > 0 {
		t.root = t.root.Children[0]
	}

	return nil
}

func siblingExist[K KeyType, V any](parent positionStr[K, V], index int) (*Node[K, V], bool) {
	index = parent.position + index

	if index < 0 || index > parent.node.pointer {
		return nil, false
	}

	sibling := parent.node.Children[index]
	return sibling, true
}

func sibling[K KeyType, V any](parent positionStr[K, V], degree int) (*Node[K, V], bool, bool) {
	var (
		potential *Node[K, V]
		side      bool
	)

	if sibling, isExist := siblingExist(parent, -1); isExist {
		if minAllowed(degree, sibling.pointer-1) {
			return sibling, true, true
		}

		potential, side = sibling, true
	}

	if sibling, isExist := siblingExist(parent, +1); isExist {
		if minAllowed(degree, sibling.pointer-1) {
			return sibling, false, true
		}

		if potential == nil {
			potential, side = sibling, false
		}
	}

	return potential, side, false
}

func sideFn(side bool, pointer int) int {
	if side {
		return 0
	}

	return pointer
}

func indexElement(index int) int {
	if index > 0 {
		return index - 1
	}

	return index
}

func merge[K KeyType, V any](current, sibling *Node[K, V], parent positionStr[K, V], leafInternal, side bool) {
	parentElement := parent.node.items[indexElement(parent.position)]
	position := sideFn(side, current.pointer)

	if leafInternal {
		if current.nextNodeL == sibling {
			current.nextNodeL = sibling.nextNodeL
			if sibling.nextNodeL != nil {
				sibling.nextNodeL.nextNodeR = current
			}
		} else {
			current.nextNodeR = sibling.nextNodeR
			if sibling.nextNodeR != nil {
				sibling.nextNodeR.nextNodeL = current
			}
		}
	} else {
		current.pointer += insert(current.items, parentElement, position)

		if !side {
			position++
		}

		insertSet(current.Children, sibling.Children[:sibling.pointer+1], position)
	}

	if side {
		deleteElement(parent.node.Children, parent.position-1, 1)
	} else {
		deleteElement(parent.node.Children, parent.position+1, 1)
	}

	current.pointer += insertSet(current.items, sibling.items[:sibling.pointer], position)

	sibling.delete()
}

func insertSet[T any](list []T, insert []T, position int) int {
	copy(list[position+len(insert):], list[position:])
	return copy(list[position:], insert)
}

func transfer[K KeyType, V any](parent, current positionStr[K, V], sibling *Node[K, V], leafInternal, side bool) {
	itemIndex := 0
	parentPosition := parent.position - 1
	childInsertPosition := current.node.pointer
	insertPosition := current.node.pointer

	if side {
		itemIndex = sibling.pointer - 1
		childInsertPosition = 0
		insertPosition = 0
	} else {
		parentPosition++
	}

	if leafInternal {
		siblingItem := sibling.items[itemIndex]
		if !side {
			siblingItem = sibling.items[1]
		}
		parent.node.items[parentPosition] = siblingItem
		current.node.pointer += insert(current.node.items, sibling.items[itemIndex], insertPosition)
	} else {
		current.node.pointer += insert(current.node.items, parent.node.items[parentPosition], childInsertPosition)
		parent.node.items[parentPosition] = sibling.items[itemIndex]

		if !side {
			insert(current.node.Children, sibling.Children[0], current.node.pointer) //check right side -> itemIndex+1(work for left) -> itemIndex
			deleteElement(sibling.Children, 0, 1)
		} else {
			insert(current.node.Children, sibling.Children[itemIndex+1], childInsertPosition) //check right side -> itemIndex+1(work for left) -> itemIndex
		}
	}

	sibling.pointer -= deleteElement(sibling.items, itemIndex, 1)
}

func (t *Tree[K, V]) TestFunc() int {
	current := t.root
	for current.Children[0] != nil {
		current = current.Children[0]
	}

	var counter int

	for current != nil {
		for _, value := range current.items[:current.pointer] {
			counter++
			// fmt.Println(counter, value)
			_ = value
		}
		// fmt.Println("======")
		current = current.nextNodeR
	}

	return counter
}

func (t *Tree[K, V]) GetRoot() *Node[K, V] {
	return t.root
}
