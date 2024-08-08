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

// check //work fine
func minAllowed(degree, numElement int) bool {
	return (degree/2)+degree%2-1 <= numElement
}

// left true, right false //work fine
func checkTransfer(parent *Node, position, degree int) (*Node, bool) {
	if position > 0 && minAllowed(degree, parent.children[position-1].pointer-1) {
		return parent.children[position-1], true
	} else if position <= parent.pointer && minAllowed(degree, parent.children[position+1].pointer-1) {
		return parent.children[position+1], false
	} else {
		return nil, false
	}
}

// left right transfer
func transferLeaf(transferNode, parent, current *Node, side bool, position int) {
	//left
	transferNode.pointer-- //mask the last element

	if position > 0 {
		position--
	}

	if side {
		current.pointer += insert(current.items, transferNode.items[transferNode.pointer], 0) //current leaf with deleted key
		parent.items[position] = transferNode.items[transferNode.pointer]                     // update parent on position right

	} else {
		current.pointer += insert(current.items, transferNode.items[0], current.pointer) //current leaf with deleted key
		deleteElement(transferNode.items, 0, 1)                                          //real delete element
		parent.items[position+1] = transferNode.items[0]
	}
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
	stack := newStack()                       //create stack
	_, found := findLeaf(t.root, &stack, key) //fill stack

	current, _ := stack.Pop()

	//Check if key exist
	if !found {
		return fmt.Errorf("not found key %d", key)
	}

	//delete key //fixPosition -> current.position-1
	current.node.pointer -= deleteElement(current.node.items, current.position-1, 1)

	if minAllowed(t.degree, current.node.pointer) {
		return nil
	}

	temp := current //leaf node

	//parent
	current, err := stack.Pop()
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	if transferNode, side := checkTransfer(current.node, current.position, t.degree); transferNode != nil {
		transferLeaf(transferNode, current.node, temp.node, side, fixPosition(current.node, current.position, key)) //check is necessary fixPosition()
		return nil
	}

	mergeLeaf(temp.node, current.node, current.position)

	// for {
	current.node.pointer -= deleteElement(current.node.items, current.position, 1)

	if minAllowed(t.degree, current.node.pointer) {
		return nil
	}

	temp = current //old parent

	// //parent
	current, err = stack.Pop() //new parent
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	// }

	// transferNode, side := canIMergeLeaf(current.node, current.position)
	// fmt.Println(transferNode)
	// fmt.Println(side)

	transferInternal(current.node, temp.node, current.position, t.degree)

	return nil
}

// nista ne valja!!! <<<
func transferInternal(parent, current *Node, position, degree int) {
	transferNode, side := checkTransfer(parent, position, degree)

	if side {
		current.items[0] = parent.items[position-1] //nece se stalno na 0 obrisati key, mora se pomjerati sve
		parent.items[position-1] = transferNode.items[transferNode.pointer-1]
		transferNode.pointer -= deleteElement(transferNode.items, transferNode.pointer, 1) //!number of element not good
		current.children[0] = transferNode.children[transferNode.pointer+1]                //children last and second last update
	} else {

	}

	transferNode.pointer--
}

// if we find the same key in the internal node then just decries
func fixPosition(node *Node, position, key int) int {
	if node.items[position-1].key == key {
		return position - 1
	}

	return position
}

// i cant find right way
// just need to return sibling
func canIMergeLeaf(parent *Node, position int) (*Node, bool) {
	if position > 0 {
		return parent.children[position-1], true
	}

	return parent.children[position+1], false
}

func mergeLeaf(current, parent *Node, position int) {
	migrateNode, side := canIMergeLeaf(parent, position)
	//left
	if side {
		current.pointer += insertSet(current.items, migrateNode.items[:migrateNode.pointer], 0)
		if migrateNode.nextNodeL != nil {
			current.nextNodeL = migrateNode.nextNodeL
		} else {
			current.nextNodeL = nil
		}

		return
	}

	current.pointer += insertSet(current.items, migrateNode.items[:migrateNode.pointer], current.pointer)
	if migrateNode.nextNodeR != nil {
		current.nextNodeR = migrateNode.nextNodeR
	} else {
		current.nextNodeR = nil
	}
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

// return position for delete(transferNode), update(parentNode),  insert(currentNode)
func transferPosition(parent, transfer, current int, side bool) (int, int, int) {
	if side {
		return transfer - 1, parent - 1, 0
	}

	return 0, parent, current
}

func transfer() {
	//delete transfer element
	//update parent
	//insert current
	//insert children if exist
}
