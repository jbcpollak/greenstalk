package decorator

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
)

// Inverter ...
func Inverter[Blackboard any](child core.Node[Blackboard]) core.Node[Blackboard] {
	base := core.NewDecorator(core.BaseParams("Inverter"), child)
	return &inverter[Blackboard]{Decorator: base}
}

// inverter ...
type inverter[Blackboard any] struct {
	core.Decorator[Blackboard, core.BaseParams]
}

func (d *inverter[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	return d.Tick(ctx, bb, evt)
}

// Tick ...
func (d *inverter[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
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
func (d *inverter[Blackboard]) Leave(bb Blackboard) error {
	return nil
}
