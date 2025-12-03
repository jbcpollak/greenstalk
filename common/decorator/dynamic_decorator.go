package decorator

import (
	"context"

	"github.com/jbcpollak/greenstalk/v2/core"
)

func DynamicDecoratorNamed(name string, childFn func() (core.Node, error)) core.Node {
	base := core.NewDynamicDecorator(core.BaseParams(name), childFn)
	return &dynamicDecorator{DynamicDecorator: base}
}

func DynamicDecorator(childFn func() (core.Node, error)) core.Node {
	return DynamicDecoratorNamed("DynamicDecorator", childFn)
}

type dynamicDecorator struct {
	core.DynamicDecorator[core.BaseParams]
}

func (d *dynamicDecorator) Activate(ctx context.Context, evt core.Event) core.ResultDetails {
	child, err := d.ChildFn()
	if err != nil {
		return core.ErrorResult(err)
	}
	child.SetNamePrefix(d.FullName())
	d.Child = child

	return d.Tick(ctx, evt)
}

func (d *dynamicDecorator) Tick(ctx context.Context, evt core.Event) core.ResultDetails {
	return core.Update(ctx, d.Child, evt)
}

func (d *dynamicDecorator) Leave(context.Context) error {
	return nil
}

var _ core.Node = (*dynamicDecorator)(nil)
