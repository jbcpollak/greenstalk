package decorator

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
	"github.com/rs/zerolog/log"
)

// UntilFailure updates its child until it returns Failure.
func UntilFailure[Blackboard any](child core.Node[Blackboard]) core.Node[Blackboard] {
	base := core.NewDecorator[Blackboard](core.BaseParams("UntilFailure"), child)
	return &untilFailure[Blackboard]{Decorator: base}
}

type untilFailure[Blackboard any] struct {
	core.Decorator[Blackboard, core.BaseParams]
}

func (d *untilFailure[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	return d.Tick(ctx, bb, evt)
}

func (d *untilFailure[Blackboard]) repeat(_ context.Context, enqueue core.EnqueueFn) error {
	log.Info().Msg("UntilFailure: repeating")
	enqueue(core.TargetNodeEvent(d.Id()))
	return nil
}

func (d *untilFailure[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	status := core.Update(ctx, d.Child, bb, evt)

	if status == core.StatusError || status == core.StatusRunning || status == core.StatusInvalid {
		return status
	}

	if status == core.StatusFailure {
		return core.StatusSuccess
	}

	return core.NodeAsyncRunning(d.repeat)
}

func (d *untilFailure[Blackboard]) Leave(bb Blackboard) {}
