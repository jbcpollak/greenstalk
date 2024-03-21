package decorator

import (
	"context"
	"fmt"
	"time"

	"github.com/jbcpollak/greenstalk/core"
)

// Delayer ...
func Delayer[Blackboard any](params core.Params, child core.Node[Blackboard]) core.Node[Blackboard] {
	base := core.NewDecorator("Delayer", params, child)

	v, err := params.Get("delay")
	if err != nil {
		panic(err)
	}
	delay, ok := v.(time.Duration)
	if !ok {
		panic(fmt.Errorf("delay must be a time.Duration"))
	}

	d := &delayer[Blackboard]{
		Decorator: base,
		delay:     delay,
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
