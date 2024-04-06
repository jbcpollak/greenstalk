package decorator

import (
	"context"
	"time"

	"github.com/jbcpollak/greenstalk/core"
)

type DelayerParams struct {
	core.BaseParams

	Delay time.Duration
}

// Delayer ...
func Delayer[Blackboard any](params DelayerParams, child core.Node[Blackboard]) core.Node[Blackboard] {
	base := core.NewDecorator(params, child)

	d := &delayer[Blackboard]{
		Decorator: base,
		delay:     params.Delay,
	}
	return d
}

// delayer ...
type delayer[Blackboard any] struct {
	core.Decorator[Blackboard, DelayerParams]
	delay time.Duration // delay in milliseconds
	start time.Time
}

// Activate ...
func (d *delayer[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	d.start = time.Now()

	return d.Tick(ctx, bb, evt)
}

// Tick ...
func (d *delayer[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	if time.Since(d.start) > d.delay {
		return core.Update(ctx, d.Child, bb, evt)
	}
	return core.StatusRunning
}

// Leave ...
func (d *delayer[Blackboard]) Leave(bb Blackboard) {}
