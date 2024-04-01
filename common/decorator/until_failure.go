package decorator

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
)

// UntilFailure updates its child until it returns Failure.
func UntilFailure[Blackboard any](params core.DecoratorParams, child core.Node[Blackboard]) core.Node[Blackboard] {
	base := core.NewDecorator[Blackboard](core.DecoratorParams{BaseParams: "UntilFailure"}, child)
	return &untilFailure[Blackboard]{Decorator: base}
}

type untilFailure[Blackboard any] struct {
	core.Decorator[Blackboard]
}

func (d *untilFailure[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	return d.Tick(ctx, bb, evt)
}

func (d *untilFailure[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	status := core.Update(ctx, d.Child, bb, evt)
	if status == core.StatusFailure {
		return core.StatusSuccess
	}
	return core.StatusRunning
}

func (d *untilFailure[Blackboard]) Leave(bb Blackboard) {}
