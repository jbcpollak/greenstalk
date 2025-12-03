package decorator

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
	"github.com/jbcpollak/greenstalk/internal"
)

type (
	UntilCondition    func(status core.ResultDetails) bool
	RepeatUntilParams struct {
		core.BaseParams

		Until UntilCondition
	}
)

// RepeatUntil updates its child n times, at which point the repeater
// returns Success. The repeater runs forever if n == 0.
func RepeatUntil(params RepeatUntilParams, child core.Node) core.Node {
	base := core.NewDecorator(params, child)
	return &repeatUntil{
		Decorator: base,
	}
}

type repeatUntil struct {
	core.Decorator[RepeatUntilParams]
}

func (d *repeatUntil) repeat(_ context.Context, enqueue core.EnqueueFn) error {
	internal.Logger.Info("Repeating", "name", d.Name())
	return enqueue(core.TargetNodeEvent(d.Id()))
}

func (d *repeatUntil) Activate(ctx context.Context, evt core.Event) core.ResultDetails {
	return d.Tick(ctx, evt)
}

func (d *repeatUntil) Tick(ctx context.Context, evt core.Event) core.ResultDetails {
	internal.Logger.Info("Repeater: Calling child")
	result := core.Update(ctx, d.Child, evt)
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

func (d *repeatUntil) Leave(context.Context) error {
	return nil
}

var _ core.Node = (*repeatUntil)(nil)
