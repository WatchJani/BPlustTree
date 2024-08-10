package BPTree

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
