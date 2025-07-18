package core

import (
	"fmt"
)

// Decorator is the base type for any specific decorator node. Such a node
// may be domain-specific, but usually one of the common nodes will be used,
// such as Inverter or Repeater. Each decorator node has Params: a key-value
// map used for setting variables for a specific decorator node, for instance
// Params{"n": 5} for a Repeater node or Params{"ms": 500} for a
// Delayer node.
type Decorator[Blackboard any, P Params] struct {
	BaseNode[P]
	Child Node[Blackboard]
}

// NewDecorator creates a new decorator base node.
func NewDecorator[Blackboard any, P Params](params P, child Node[Blackboard]) Decorator[Blackboard, P] {
	child.SetNamePrefix(params.Name())
	return Decorator[Blackboard, P]{
		BaseNode: newBaseNode(CategoryDecorator, params),
		Child:    child,
	}
}

func (c *Decorator[Blackboard, P]) Walk(walkFn WalkFunc[Blackboard], level int) {
	walkFn(c, level)
	c.Child.Walk(walkFn, level+1)
}

// String returns a string representation of the decorator node.
func (d *Decorator[Blackboard, P]) String() string {
	return fmt.Sprintf("* %s (%v)", d.Params.Name(), d.Params)
}

func (d *Decorator[Blackboard, P]) SetNamePrefix(namePrefix string) {
	d.BaseNode.SetNamePrefix(namePrefix)
	d.Child.SetNamePrefix(d.FullName())
}
