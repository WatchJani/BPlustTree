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

// check
func minAllowed(degree, numElement int) bool {
	return (degree/2)+degree%2-1 <= numElement
}

// left true, right false
func checkTransfer(parent positionStr, degree int) (*Node, bool) {
	//left
	node, position := parent.node, parent.position
	if position > 0 && minAllowed(degree, node.children[position-1].pointer-1) {
		return node.children[position-1], true
	} else if position < node.pointer && minAllowed(degree, node.children[position+1].pointer-1) {
		return node.children[position+1], false
	} else {
		return nil, false
	}
}

// left right transfer
func transferLeaf(transferNode, currentNode, parent *Node, side bool, position int) {
	//left
	transferNode.pointer--
	fmt.Println(transferNode)
	if side {
		parent.pointer += insert(parent.items, transferNode.items[transferNode.pointer], 0) //current leaf with deleted key
		currentNode.items[position-1] = transferNode.items[transferNode.pointer]            // update parent on position right

		//check -1 -> [position-1]
	} else {
		parent.pointer += insert(parent.items, transferNode.items[0], parent.pointer) //current leaf with deleted key
		currentNode.items[position-1] = transferNode.items[0]                         // update parent on position right
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
	stack := newStack()
	position, found := findLeaf(t.root, &stack, key)

	current, _ := stack.Pop() //cant be error because we have always root

	if !found {
		return fmt.Errorf("not found key %d", key)
	}

	current.node.pointer -= deleteElement(current.node.items, position, 1)

	if minAllowed(t.degree, current.node.pointer) {
		return nil
	}
	currentLeaf := current.node

	//parent
	current, err := stack.Pop()
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	transferNode, side := checkTransfer(current, t.degree)
	if transferNode != nil {
		fmt.Println(key, side)
		transferLeaf(transferNode, current.node, currentLeaf, side, current.position)
		return nil
	}

	return nil
}
