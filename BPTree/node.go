package BPTree

type Node struct {
	items     []item
	children  []*Node
	nextNodeL *Node
	nextNodeR *Node
	pointer   int
}

func newNode(degree int) Node {
	return Node{
		items:    make([]item, degree+1),
		children: make([]*Node, degree+2),
	}
}

func (n *Node) delete() {
	n.items = nil
	n.children = nil
	n.pointer = 0
	n.nextNodeL = nil
	n.nextNodeR = nil
}
