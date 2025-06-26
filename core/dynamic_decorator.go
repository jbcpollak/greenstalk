package core

import "fmt"

type DynamicDecorator[Blackboard any, P Params] struct {
	BaseNode[P]
	Child   Node[Blackboard]
	ChildFn func() (Node[Blackboard], error)
}

func NewDynamicDecorator[Blackboard any, P Params](params P, childFn func() (Node[Blackboard], error)) DynamicDecorator[Blackboard, P] {
	return DynamicDecorator[Blackboard, P]{
		BaseNode: newBaseNode(CategoryDecorator, params),
		ChildFn:  childFn,
	}
}

func (c *DynamicDecorator[Blackboard, P]) Walk(walkFn WalkFunc[Blackboard], level int) {
	walkFn(c, level)
	if c.Child != nil {
		c.Child.Walk(walkFn, level+1)
	}
}

func (d *DynamicDecorator[Blackboard, P]) String() string {
	return fmt.Sprintf("*d %s (%v)", d.Params.Name(), d.Params)
}

func (d *DynamicDecorator[Blackboard, P]) SetName(newName string) Walkable[Blackboard] {
	d.Params = d.Params.SetName(newName).(P)
	return d
}
