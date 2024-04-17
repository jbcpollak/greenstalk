package decorator

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
	"github.com/rs/zerolog/log"
)

type UntilCondition func(status core.ResultDetails) bool
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
	}
}

type repeatUntil[Blackboard any] struct {
	core.Decorator[Blackboard, RepeatUntilParams]
}

func (d *repeatUntil[Blackboard]) repeat(_ context.Context, _ Blackboard, enqueue core.EnqueueFn) error {
	log.Info().Msg(d.Name() + ": repeating")
	enqueue(core.TargetNodeEvent(d.Id()))
	return nil
}

func (d *repeatUntil[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {

	return d.Tick(ctx, bb, evt)
}

func (d *repeatUntil[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	log.Info().Msg("Repeater: Calling child")
	result := core.Update(ctx, d.Child, bb, evt)
	status := result.Status()

	if status == core.StatusError ||
		status == core.StatusInvalid ||
		status == core.StatusRunning {
		return result
	}

	if d.Params.Until(result) {
		return core.SuccessResult()
	}

	return core.InitRunningResult(d.repeat)
}

func (d *repeatUntil[Blackboard]) Leave(bb Blackboard) error {
	return nil
}
