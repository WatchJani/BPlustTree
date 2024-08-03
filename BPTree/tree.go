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
			return mid, true
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
		store: make([]positionStr, 4),
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

func (t *Tree) Insert(key, value int) {
	var (
		position int
		found    bool
		stack    = newStack()
		item     = item{key: key, value: value}
	)

	for current := t.root; current != nil; {
		position, found = current.search(key)
		stack.Push(current, position)

		current = current.children[position]
	}

	current, _ := stack.Pop()

	//update just state state
	if found {
		current.node.items[position].value = value
		return
	}

	//insert to leaf and update state
	if middleKey, nodeChildren := insertLeaf(current.node, position, t.degree, item); nodeChildren != nil {
		for {
			temp := current
			current, err := stack.Pop()
			if err != nil {
				current = temp
				break
			}

			current.node.pointer += insert(current.node.items, middleKey, current.position)
			chIndex := childrenIndex(middleKey.key, current.node.items[current.position].key, current.position)
			insertChildren(current.node.children, temp.node, chIndex) //insert pointer on children

			if current.node.pointer < t.degree {
				fmt.Println("yes")
				return
			}

			//split
		}

		rootNode := newNode(t.degree) // new root node
		rootNode.pointer += insert(rootNode.items, middleKey, 0)
		rootNode.children[0] = nodeChildren
		rootNode.children[1] = current.node
		t.root = &rootNode

		fmt.Println("root:", rootNode)
		fmt.Println("=============================")
	}
}

func childrenIndex(key, value, index int) int {
	if value < key {
		return index + 1
	}

	return index
}

func insertChildren(list []*Node, insert *Node, position int) int {
	copy(list[position+1:], list[position:])
	return copy(list[position:], []*Node{insert})
}

func insert(list []item, insert item, position int) int {
	copy(list[position+1:], list[position:])
	return copy(list[position:], []item{insert})
}

func migrateElement(list, migrateElement []item, position int) int {
	return copy(list[position:], migrateElement)
}

func deleteElement(list []item, position, deletion int) int {
	copy(list[position:], list[position+deletion:])
	return deletion
}

func insertLeaf(current *Node, position, degree int, item item) (item, *Node) {
	current.pointer += insert(current.items, item, position)

	if current.pointer < degree {
		fmt.Println(current)
		return item, nil //
	}
	fmt.Println(current)

	//Split
	newNode := newNode(degree)
	middle := degree / 2 //Check

	newNode.pointer += migrateElement(newNode.items, current.items[:middle], 0)
	current.pointer -= deleteElement(current.items, 0, current.pointer-middle-1)

	//fix update with all nodes around current node
	newNode.nextNode = current

	return current.items[0], &newNode
}
