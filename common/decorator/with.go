package decorator

import (
	"context"
	"io"

	"github.com/jbcpollak/greenstalk/core"
)

func With[Blackboard any](child core.Node[Blackboard]) core.Node[Blackboard] {
	base := core.NewDecorator(core.BaseParams("With"), child)
	return &with[Blackboard]{Decorator: base}
}

type with[Blackboard any] struct {
	core.Decorator[Blackboard, core.BaseParams]
	CreateCloseable func() (io.Closer, error)
	closeable       io.Closer
}

func (d *with[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	closeable, err := d.CreateCloseable()
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
	switch result := core.Update(ctx, d.Child, bb, evt); result {
	case core.StatusSuccess:
		return core.StatusFailure
	case core.StatusFailure:
		return core.StatusSuccess
	default:
		return result
	}
}

// Leave ...
func (d *with[Blackboard]) Leave(bb Blackboard) error {
	return d.closeable.Close()
}