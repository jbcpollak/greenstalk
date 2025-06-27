package decorator

import (
	"context"
	"io"

	"github.com/jbcpollak/greenstalk/core"
)

func WithNamed[Blackboard any](name string, createCloseable func() (io.Closer, error), child core.Node[Blackboard]) core.Node[Blackboard] {
	base := core.NewDecorator(core.BaseParams(name), child)
	return &with[Blackboard]{Decorator: base, createCloseable: createCloseable}
}
func With[Blackboard any](createCloseable func() (io.Closer, error), child core.Node[Blackboard]) core.Node[Blackboard] {
	return WithNamed("With", createCloseable, child)
}

type with[Blackboard any] struct {
	core.Decorator[Blackboard, core.BaseParams]
	createCloseable func() (io.Closer, error)
	closeable       io.Closer
}

func (d *with[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	closeable, err := d.createCloseable()
	if err != nil {
		return core.ErrorResult(err)
	} else {
		d.closeable = closeable
	}
	return d.Tick(ctx, bb, evt)
}

// Tick ...
func (d *with[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	return core.Update(ctx, d.Child, bb, evt)
}

// Leave ...
func (d *with[Blackboard]) Leave(bb Blackboard) error {
	return d.closeable.Close()
}
