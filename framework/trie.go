package framework

type Tree struct {
	root *node
}
type node struct {
	isLast  bool
	semgent string
	handler ControllerHandler
	childs  []*node
}

func newNode() *node {
	return &node{
		isLast:  false,
		semgent: "",
		childs:  []*node{},
	}
}

func NewTree() *Tree {
	root := newNode()
	return &Tree{root}
}
