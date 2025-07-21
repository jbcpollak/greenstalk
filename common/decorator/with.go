package decorator

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
)

func WithNamed[Blackboard any](name string, createCloseable func() (func() error, error), child core.Node[Blackboard]) core.Node[Blackboard] {
	base := core.NewDecorator(core.BaseParams(name), child)
	return &with[Blackboard]{Decorator: base, createCloseable: createCloseable}
}
func With[Blackboard any](createCloseable func() (func() error, error), child core.Node[Blackboard]) core.Node[Blackboard] {
	return WithNamed("With", createCloseable, child)
}

type with[Blackboard any] struct {
	core.Decorator[Blackboard, core.BaseParams]
	createCloseable func() (func() error, error)
	closeFn         func() error
}

func (d *with[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	closeable, err := d.createCloseable()
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
func (d *with[Blackboard]) Leave(bb Blackboard) error {
	return d.closeFn()
}
