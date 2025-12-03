package decorator

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
)

func WithNamed[Blackboard any](name string, createCloseable func() (closeFn func() error, err error), child core.Node[Blackboard]) core.Node[Blackboard] {
	base := core.NewDecorator(core.BaseParams(name), child)
	return &with[Blackboard]{Decorator: base, createCloseable: func(ctx context.Context) (func(context.Context) error, error) {
		closeFn, err := createCloseable()
		if err != nil {
			return nil, err
		}
		return func(context.Context) error {
			return closeFn()
		}, nil
	}}
}

func With[Blackboard any](createCloseable func() (closeFn func() error, err error), child core.Node[Blackboard]) core.Node[Blackboard] {
	return WithNamed("With", createCloseable, child)
}

func WithNamedContext[Blackboard any](name string, createCloseable func(context.Context) (closeFn func(context.Context) error, err error), child core.Node[Blackboard]) core.Node[Blackboard] {
	base := core.NewDecorator(core.BaseParams(name), child)
	return &with[Blackboard]{Decorator: base, createCloseable: createCloseable}
}

type with[Blackboard any] struct {
	core.Decorator[Blackboard, core.BaseParams]
	createCloseable func(context.Context) (func(context.Context) error, error)
	closeFn         func(context.Context) error
}

func (d *with[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	closeable, err := d.createCloseable(ctx)
	if err != nil {
		return core.ErrorResult(err)
	} else {
		d.closeFn = closeable
	}
	return d.Tick(ctx, bb, evt)
}

// Tick ...
func (d *with[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	return core.Update(ctx, d.Child, bb, evt)
}

// Leave ...
func (d *with[Blackboard]) Leave(ctx context.Context, _ Blackboard) error {
	return d.closeFn(ctx)
}
