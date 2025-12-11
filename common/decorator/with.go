package decorator

import (
	"context"

	"github.com/jbcpollak/greenstalk/v2/core"
)

func WithNamed(
	name string,
	createCloseable func(context.Context) (closeFn func(context.Context) error, err error),
	child core.Node,
) core.Node {
	base := core.NewDecorator(core.BaseParams(name), child)
	return &with{Decorator: base, createCloseable: createCloseable}
}

func With(
	createCloseable func(context.Context) (closeFn func(context.Context) error, err error),
	child core.Node,
) core.Node {
	return WithNamed("With", createCloseable, child)
}

type with struct {
	core.Decorator[core.BaseParams]
	createCloseable func(context.Context) (func(context.Context) error, error)
	closeFn         func(context.Context) error
}

func (d *with) Activate(ctx context.Context, evt core.Event) core.ResultDetails {
	closeable, err := d.createCloseable(ctx)
	if err != nil {
		return core.ErrorResult(err)
	} else {
		d.closeFn = closeable
	}
	return d.Tick(ctx, evt)
}

func (d *with) Tick(ctx context.Context, evt core.Event) core.ResultDetails {
	return core.Update(ctx, d.Child, evt)
}

func (d *with) Leave(ctx context.Context) error {
	return d.closeFn(ctx)
}
