package BPTree

import (
	"errors"
	"fmt"
)

type Tree struct {
	root   *Node
	degree int
	//length int
}

func New(degree int) *Tree {
	if degree < 3 {
		degree = 3
	}

	return &Tree{
		degree: degree,
		root: &Node{
			items:    make([]item, degree+1),
			children: make([]*Node, degree+2),
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

		next = next.children[index]
	}

	return -1, errors.New("key not found")
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
			insert(parent.children, nodeChildren, chIndex) //insert pointer on children

			if parent.pointer < t.degree {
				return
			}

			middle := t.degree / 2
			middleKey = parent.items[middle]

			//split
			newNode := newNode(t.degree)
			newNode.pointer += migrate(newNode.items, parent.items[:middle], 0) //migrate half element to left child node
			migrate(newNode.children, parent.children[:middle+1], 0)

			parent.pointer -= deleteElement(parent.items, 0, parent.pointer-middle)
			migrate(parent.children, parent.children[middle+t.degree%2:], 0)

			nodeChildren = &newNode

			current = stack //fix this part
		}

		rootNode := newNode(t.degree)
		rootNode.pointer += insert(rootNode.items, middleKey, 0)
		rootNode.children[0] = nodeChildren
		rootNode.children[1] = current.node
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

func insertSet[T any](list []T, insert []T, position int) int {
	copy(list[position+len(insert):], list[position:])
	return copy(list[position:], insert)
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

	// fmt.Println(current)

	return current.items[0], &newNode
}

// check //work fine
func minAllowed(degree, numElement int) bool {
	return (degree/2)+degree%2-1 <= numElement
}

func findLeaf(root *Node, stack *Stack, key int) (int, bool) {
	position, found := 0, false

	for current := root; current != nil; {
		position, found = current.search(key)
		stack.Push(current, position)

		current = current.children[position]
	}

	return position, found
}

func (t *Tree) Delete(key int) error {
	// stack := newStack()                       //create stack
	// _, found := findLeaf(t.root, &stack, key) //fill stack

	return nil
}

// ==================================================================
func siblingExist(parent *Node, index int) (*Node, bool) {
	sibling := parent.children[parent.pointer+index]
	return sibling, sibling != nil
}

// return sibling, side and [transfer/merge]
func sibling(parent positionStr, degree int) (*Node, bool, bool) {
	var (
		potential *Node
		side      bool
	)

	if sibling, isExist := siblingExist(parent.node, -1); isExist {
		if minAllowed(degree, sibling.pointer) {
			return sibling, true, false
		}

		potential, side = sibling, true
	}

	if sibling, isExist := siblingExist(parent.node, +1); isExist {
		if minAllowed(degree, sibling.pointer) {
			return sibling, false, false
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

func merge(currentNode, mergeNode *Node, parentElement item, leafInternal, side bool) {
	position := sideFn(side, currentNode.pointer)

	if leafInternal {
		if currentNode.nextNodeL == mergeNode {
			currentNode.nextNodeL = mergeNode.nextNodeL
			mergeNode.nextNodeL.nextNodeR = currentNode
		} else {
			currentNode.nextNodeR = mergeNode.nextNodeR
			mergeNode.nextNodeR.nextNodeL = currentNode
		}
	} else {
		//update children
		migrate(currentNode.children, mergeNode.children[:mergeNode.pointer], position+1)
		//insert parent node
		currentNode.pointer += insert(currentNode.items, parentElement, position)
	}

	currentNode.pointer += migrate(currentNode.items, mergeNode.items[:mergeNode.pointer], position)
	mergeNode = nil
}

// side true left sibling
func transfer(parent, current, sibling positionStr, leafInternal, side bool) {
	itemIndex := 0
	parentPosition := parent.position - 1
	childInsertPosition := current.node.pointer
	insertPosition := current.node.pointer

	if side {
		itemIndex = sibling.node.pointer - 1
		childInsertPosition = 0
		insertPosition = 0
	} else {
		parentPosition++
	}

	if leafInternal {
		siblingItem := sibling.node.items[itemIndex]
		if !side {
			siblingItem = sibling.node.items[1]
		}
		parent.node.items[parentPosition] = siblingItem
		current.node.pointer += insert(current.node.items, sibling.node.items[itemIndex], insertPosition)
	} else {
		current.node.pointer += insert(current.node.items, parent.node.items[parentPosition], childInsertPosition)
		parent.node.items[parentPosition] = sibling.node.items[itemIndex]
		insert(current.node.children, sibling.node.children[itemIndex], childInsertPosition)
		if !side {
			deleteElement(sibling.node.children, 0, 1)
		}
	}

	sibling.node.pointer -= deleteElement(sibling.node.items, itemIndex, 1)
	// if leafInternal { //leaf
	// 	itemIndex := 0
	// 	if side {
	// 		itemIndex = sibling.position - 1
	// 		parent.node.items[parent.position-1] = sibling.node.items[itemIndex]
	// 		current.node.pointer += insert(current.node.items, sibling.node.items[itemIndex], 0)
	// 	} else {
	// 		current.node.pointer += insert(current.node.items, sibling.node.items[0], current.node.pointer)
	// 		parent.node.items[parent.position-1] = sibling.node.items[1]
	// 	}
	// 	sibling.node.pointer -= deleteElement(sibling.node.items, itemIndex, 1)
	// } else {
	// 	parentPosition := parent.position - 1
	// 	itemIndex := 0
	// 	childInsertPosition := current.node.pointer

	// 	if side {
	// 		itemIndex = sibling.node.pointer - 1
	// 		childInsertPosition = 0
	// 	} else {
	// 		parentPosition++
	// 	}

	// 	current.node.pointer += insert(current.node.items, parent.node.items[parentPosition], childInsertPosition)
	// 	parent.node.items[parentPosition] = sibling.node.items[itemIndex]
	// 	sibling.node.pointer -= deleteElement(sibling.node.items, itemIndex, 1)

	// 	insert(current.node.children, sibling.node.children[itemIndex], childInsertPosition)

	// 	if !side {
	// 		deleteElement(sibling.node.children, 0, 1)
	// 	}

	// if side {
	// 	parentPosition := parent.position - 1
	// 	current.node.pointer += insert(current.node.items, parent.node.items[parentPosition], 0)
	// 	parent.node.items[parentPosition] = sibling.node.items[sibling.node.pointer-1]
	// 	sibling.node.pointer -= deleteElement(sibling.node.items, sibling.node.pointer-1, 1)
	// 	insert(current.node.children, sibling.node.children[sibling.node.pointer], 0)
	// } else {
	// 	parentPosition := parent.position
	// 	current.node.pointer += insert(current.node.items, parent.node.items[parentPosition], current.node.pointer)
	// 	parent.node.items[parentPosition] = sibling.node.items[0]
	// 	sibling.node.pointer -= deleteElement(sibling.node.items, 0, 1)
	// 	insert(current.node.children, sibling.node.children[0], current.node.pointer)
	// 	deleteElement(sibling.node.children, 0, 1)
	// }
	// }
}

func (t *Tree) TestFunc() {
	current := t.root
	for current.children[0] != nil {
		current = current.children[0]
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
}
