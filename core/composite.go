package core

// Composite is the base type for any specific composite node. Such a node
// may be domain-specific, but usually one of the common nodes will be used,
// such as Sequence or Selector.
type Composite[Blackboard any, P Params] struct {
	BaseNode[P]
	Children     []Node[Blackboard]
	CurrentChild int // TODO - move into instance nodes
}

// NewComposite creates a new composite base node.
func NewComposite[Blackboard any, P Params](params P, children []Node[Blackboard]) Composite[Blackboard, P] {
	for _, child := range children {
		child.SetNamePrefix(params.Name())
	}
	return Composite[Blackboard, P]{
		BaseNode: newBaseNode(CategoryComposite, params),
		Children: children,
	}
}

func (c *Composite[Blackboard, P]) Walk(walkFn WalkFunc[Blackboard], level int) {
	walkFn(c, level)
	for _, child := range c.Children {
		child.Walk(walkFn, level+1)
	}
}

// String returns a string representation of the composite node.
func (c *Composite[Blackboard, P]) String() string {
	return "+ " + c.Params.Name()
}

func (c *Composite[Blackboard, P]) SetNamePrefix(namePrefix string) {
	c.BaseNode.SetNamePrefix(namePrefix + c.Name())
	for _, child := range c.Children {
		child.SetNamePrefix(c.FullName())
	}
}
