package BPTree

type Node struct {
	items    []item
	children []*Node
	nextNode *Node
	pointer  int
}

func newNode(degree int) Node {
	return Node{
		items:    make([]item, degree+1),
		children: make([]*Node, degree+2),
	}
}
