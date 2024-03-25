package decorator

import (
	"context"
	"time"

	"github.com/jbcpollak/greenstalk/core"
)

type DelayerParams struct {
	core.DecoratorParams

	Delay time.Duration
}

// Delayer ...
func Delayer[Blackboard any](params DelayerParams, child core.Node[Blackboard]) core.Node[Blackboard] {
	base := core.NewDecorator(params.DecoratorParams, child)

	d := &delayer[Blackboard]{
		Decorator: base,
		delay:     params.Delay,
	}
	return d
}

// delayer ...
type delayer[Blackboard any] struct {
	core.Decorator[Blackboard]
	delay time.Duration // delay in milliseconds
	start time.Time
}

// Enter ...
func (d *delayer[Blackboard]) Enter(bb Blackboard) {
	d.start = time.Now()
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
