package decorator

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
)

func DynamicDecoratorNamed[Blackboard any](name string, childFn func() (core.Node[Blackboard], error)) core.Node[Blackboard] {
	base := core.NewDynamicDecorator(core.BaseParams(name), childFn)
	return &dynamicDecorator[Blackboard]{DynamicDecorator: base}
}
func DynamicDecorator[Blackboard any](childFn func() (core.Node[Blackboard], error)) core.Node[Blackboard] {
	return DynamicDecoratorNamed("DynamicDecorator", childFn)
}

type dynamicDecorator[Blackboard any] struct {
	core.DynamicDecorator[Blackboard, core.BaseParams]
}

func (d *dynamicDecorator[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	child, err := d.ChildFn()
	if err != nil {
		return core.ErrorResult(err)
	}
	d.Child = child

	return d.Tick(ctx, bb, evt)
}

func (d *dynamicDecorator[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	return core.Update(ctx, d.Child, bb, evt)
}

func (d *dynamicDecorator[Blackboard]) Leave(bb Blackboard) error {
	d.Child = nil
	return nil
}
