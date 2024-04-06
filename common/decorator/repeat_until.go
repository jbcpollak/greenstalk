package decorator

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
	"github.com/rs/zerolog/log"
)

type UntilCondition func(status core.NodeResult) bool
type RepeatUntilParams struct {
	core.BaseParams

	Until UntilCondition
}

// RepeatUntil updates its child n times, at which point the repeater
// returns Success. The repeater runs forever if n == 0.
func RepeatUntil[Blackboard any](params RepeatUntilParams, child core.Node[Blackboard]) core.Node[Blackboard] {
	base := core.NewDecorator(params, child)
	return &repeatUntil[Blackboard]{
		Decorator: base,
		until:     params.Until,
	}
}

type repeatUntil[Blackboard any] struct {
	core.Decorator[Blackboard, RepeatUntilParams]

	until UntilCondition
}

func (d *repeatUntil[Blackboard]) repeat(_ context.Context, enqueue core.EnqueueFn) error {
	log.Info().Msg(d.Name() + ": repeating")
	enqueue(core.TargetNodeEvent(d.Id()))
	return nil
}

func (d *repeatUntil[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {

	return d.Tick(ctx, bb, evt)
}

func (d *repeatUntil[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	log.Info().Msg("Repeater: Calling child")
	status := core.Update(ctx, d.Child, bb, evt)

	if status == core.StatusError || status == core.StatusRunning || status == core.StatusInvalid {
		return status
	}

	if d.until(status) {
		return core.StatusSuccess
	}

	return core.NodeAsyncRunning(d.repeat)
}

func (d *repeatUntil[Blackboard]) Leave(bb Blackboard) {}
