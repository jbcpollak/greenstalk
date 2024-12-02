package core

type RuntimeComposite[Blackboard any, P Params] struct {
	BaseNode[P]
	ChildrenFn   func() []Node[Blackboard]
	Children     []Node[Blackboard]
	CurrentChild int
}

func NewRuntimeComposite[Blackboard any, P Params](params P, childrenFn func() []Node[Blackboard]) RuntimeComposite[Blackboard, P] {
	return RuntimeComposite[Blackboard, P]{
		BaseNode:   newBaseNode(CategoryComposite, params),
		ChildrenFn: childrenFn,
	}
}

func (c *RuntimeComposite[Blackboard, P]) Walk(walkFn WalkFunc[Blackboard], level int) {
	walkFn(c, level)
	for _, child := range c.Children {
		child.Walk(walkFn, level+1)
	}
}

// String returns a string representation of the composite node.
func (c *RuntimeComposite[Blackboard, P]) String() string {
	return "+ " + c.Params.Name()
}
