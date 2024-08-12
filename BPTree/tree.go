package BPTree

import (
	"errors"
	"fmt"
)

type item struct {
	key   int
	value int
}

func newItem(key, value int) item {
	return item{
		key:   key,
		value: value,
	}
}

type Node struct {
	items     []item
	Children  []*Node
	nextNodeL *Node
	nextNodeR *Node
	pointer   int
}

func newNode(degree int) Node {
	return Node{
		items:    make([]item, degree+1),
		Children: make([]*Node, degree+2),
	}
}

func (n *Node) delete() {
	n.items = nil
	n.Children = nil
	n.pointer = 0
	n.nextNodeL = nil
	n.nextNodeR = nil
}

type Tree struct {
	root   *Node
	degree int
}

func New(degree int) *Tree {
	if degree < 3 {
		degree = 3
	}

	return &Tree{
		degree: degree,
		root: &Node{
			items:    make([]item, degree+1),
			Children: make([]*Node, degree+2),
		},
	}
}

func (n *Node) search(target int) (int, bool) {
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

type Stack struct {
	store []positionStr
}

type positionStr struct {
	node     *Node
	position int
}

func newStack() Stack {
	return Stack{
		store: make([]positionStr, 0, 4),
	}
}

func (s *Stack) Push(node *Node, position int) {
	s.store = append(s.store, positionStr{
		node:     node,
		position: position,
	})
}

func (s *Stack) Pop() (positionStr, error) {
	if len(s.store) == 0 {
		return positionStr{}, errors.New("stack is empty")
	}

	pop := s.store[len(s.store)-1]
	s.store = s.store[:len(s.store)-1]
	return pop, nil
}

func (t *Tree) Find(key int) (int, error) {
	for next := t.root; next != nil; {
		index, found := next.search(key)

		if found {
			return next.items[index-1].value, nil
		}

		next = next.Children[index]
	}

	return -1, fmt.Errorf("key %v not found", key)
}

func (t *Tree) Insert(key, value int) {
	stack, item := newStack(), newItem(key, value)
	position, found := findLeaf(t.root, &stack, key)

	current, _ := stack.Pop()
	//update just state state
	if found {
		current.node.items[position].value = value
		return
	}

	//insert to leaf and update state
	if middleKey, nodeChildren := insertLeaf(current.node, position, t.degree, item); nodeChildren != nil {
		for {
			temp := current //right leaf

			stack, err := stack.Pop() //get parent
			if err != nil {
				current = temp
				break
			}

			parent := stack.node

			//add to parent new item
			parent.pointer += insert(parent.items, middleKey, stack.position)
			chIndex := childrenIndex(middleKey.key, parent.items[stack.position].key, stack.position)
			//make good link
			insert(parent.Children, nodeChildren, chIndex) //insert pointer on children
			if parent.pointer < t.degree {
				return
			}

			middle := parent.pointer / 2
			middleKey = parent.items[middle]

			//split
			newNode := newNode(t.degree)
			newNode.pointer += migrate(newNode.items, parent.items[:middle], 0) //migrate half element to left child node
			migrate(newNode.Children, parent.Children[:middle+1], 0)

			parent.pointer -= deleteElement(parent.items, 0, parent.pointer-middle+1-t.degree&1) // parent.pointer-middle+1-t.degree%2
			migrate(parent.Children, parent.Children[middle+1:], 0)                              //
			nodeChildren = &newNode

			current = stack //fix this part
		}

		rootNode := newNode(t.degree)
		rootNode.pointer += insert(rootNode.items, middleKey, 0)
		rootNode.Children[0] = nodeChildren
		rootNode.Children[1] = current.node
		t.root = &rootNode
	}
}

func childrenIndex(key, value, index int) int {
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

func insertLeaf(current *Node, position, degree int, item item) (item, *Node) {
	current.pointer += insert(current.items, item, position)

	if current.pointer < degree {
		return item, nil //
	}

	//Split
	newNode := newNode(degree)
	middle := degree / 2 //Check

	newNode.pointer += migrate(newNode.items, current.items[:middle], 0)
	current.pointer -= deleteElement(current.items, 0, current.pointer-middle-degree%2)

	//update links between leafs
	if current.nextNodeL != nil {
		newNode.nextNodeL = current.nextNodeL  //left connection
		current.nextNodeL.nextNodeR = &newNode //right connection prevues left node to new node
	}

	current.nextNodeL = &newNode //left connection
	newNode.nextNodeR = current  //right connection

	return current.items[0], &newNode
}

// check //work fine
func minAllowed(degree, numElement int) bool {
	return (degree-1)/2 <= numElement
}

func findLeaf(root *Node, stack *Stack, key int) (int, bool) {
	position, found := 0, false

	for current := root; current != nil; {
		position, found = current.search(key)
		stack.Push(current, position)

		current = current.Children[position]
	}

	return position, found
}

func (t *Tree) Delete(key int) error {
	stack := newStack()                       //create stack
	_, found := findLeaf(t.root, &stack, key) //fill stack

	if !found {
		return fmt.Errorf("key %v is not exist", key)
	}

	current, _ := stack.Pop()

	for {
		current.node.pointer -= deleteElement(current.node.items, indexElement(current.position), 1)

		if minAllowed(t.degree, current.node.pointer) || (found && t.root.pointer == 0) {
			return nil
		}

		temp := current //current
		parent, err := stack.Pop()
		if err != nil {
			break
		}

		sibling, side, operation := sibling(parent, t.degree)

		if operation {
			transfer(parent, temp, sibling, found, side)
			return nil
		} else {
			merge(temp.node, sibling, parent, found, side)
		}

		if found {
			found = !found
		}

		current = parent
	}

	if t.root.pointer == 0 && len(t.root.Children) > 0 {
		// If the root is empty, promote the first child as the new root
		if t.root.Children[0] != nil {
			t.root = t.root.Children[0]
		} else {
			t.root = t.root.Children[1]
		}
	}

	//update root
	return nil
}

func siblingExist(parent positionStr, index int) (*Node, bool) {
	index = parent.position + index

	if index < 0 || index > parent.node.pointer {
		return nil, false
	}

	sibling := parent.node.Children[index]
	return sibling, true
}

// return sibling, side and [transfer/merge]
func sibling(parent positionStr, degree int) (*Node, bool, bool) {
	var (
		potential *Node
		side      bool
	)

	if sibling, isExist := siblingExist(parent, -1); isExist {
		if minAllowed(degree, sibling.pointer-1) {
			return sibling, true, true
		}

		potential, side = sibling, true
	}

	if sibling, isExist := siblingExist(parent, +1); isExist {
		if minAllowed(degree, sibling.pointer-1) { //When we delete item -> -1
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

func merge(current, sibling *Node, parent positionStr, leafInternal, side bool) {
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
		//add parent element from
		current.pointer += insert(current.items, parentElement, position)

		if !side {
			position++
		}

		insertSet(current.Children, sibling.Children[:sibling.pointer+1], position)
	}

	//delete parent element
	if side { //left side sibling delete
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

func transfer(parent, current positionStr, sibling *Node, leafInternal, side bool) {
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

	//delete sibling element
	sibling.pointer -= deleteElement(sibling.items, itemIndex, 1)
}

func (t *Tree) TestFunc() int {
	current := t.root
	for current.Children[0] != nil {
		current = current.Children[0]
	}

	var counter int

	for current != nil {
		for _, value := range current.items[:current.pointer] {
			counter++
			fmt.Println(counter, value)
		}
		fmt.Println("======")
		current = current.nextNodeR
	}

	return counter
}

func (t *Tree) GetRoot() *Node {
	return t.root
}
