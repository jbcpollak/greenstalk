package core

import "fmt"

type DynamicDecorator[P Params] struct {
	BaseNode[P]
	Child   Node
	ChildFn func() (Node, error)
}

func NewDynamicDecorator[P Params](params P, childFn func() (Node, error)) DynamicDecorator[P] {
	return DynamicDecorator[P]{
		BaseNode: newBaseNode(CategoryDecorator, params),
		ChildFn:  childFn,
	}
}

func (c *DynamicDecorator[P]) Walk(walkFn WalkFunc, level int) {
	walkFn(c, level)
	if c.Child != nil {
		c.Child.Walk(walkFn, level+1)
	}
}

func (d *DynamicDecorator[P]) String() string {
	return fmt.Sprintf("*d %s (%v)", d.Params.Name(), d.Params)
}

func (d *DynamicDecorator[P]) SetNamePrefix(namePrefix string) {
	d.BaseNode.SetNamePrefix(namePrefix)
	if d.Child != nil {
		d.Child.SetNamePrefix(d.FullName())
	}
}
