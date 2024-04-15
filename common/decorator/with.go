package decorator

import (
	"context"
	"io"

	"github.com/jbcpollak/greenstalk/core"
)

func With[Blackboard any](child core.Node[Blackboard], createCloseable func() (io.Closer, error)) core.Node[Blackboard] {
	base := core.NewDecorator(core.BaseParams("With"), child)
	return &with[Blackboard]{Decorator: base, createCloseable: createCloseable}
}

type with[Blackboard any] struct {
	core.Decorator[Blackboard, core.BaseParams]
	createCloseable func() (io.Closer, error)
	closeable       io.Closer
}

func (d *with[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	closeable, err := d.createCloseable()
	if err != nil {
		return core.NodeRuntimeError{
			Err: err,
		}
	} else {
		d.closeable = closeable
	}
	return d.Tick(ctx, bb, evt)
}

// Tick ...
func (d *with[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	return core.Update(ctx, d.Child, bb, evt)
}

// Leave ...
func (d *with[Blackboard]) Leave(bb Blackboard) error {
	return d.closeable.Close()
}
