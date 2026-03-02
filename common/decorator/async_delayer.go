package decorator

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jbcpollak/greenstalk/v2/core"
	"github.com/jbcpollak/greenstalk/v2/internal"
)

type AsyncDelayerParams struct {
	core.BaseParams

	Delay time.Duration
}

// AsyncDelayer ...
func AsyncDelayer(params AsyncDelayerParams, child core.Node) core.Node {
	base := core.NewDecorator(params, child)

	d := &asyncdelayer{
		Decorator: base,
		delay:     params.Delay,
	}
	return d
}

// delayer ...
type asyncdelayer struct {
	core.Decorator[AsyncDelayerParams]
	delay time.Duration // delay in milliseconds
	start time.Time
}

type DelayerFinishedEvent struct {
	targetNodeId uuid.UUID
	start        time.Time
}

func (e DelayerFinishedEvent) TargetNodeId() uuid.UUID {
	return e.targetNodeId
}

func (d *asyncdelayer) doDelay(ctx context.Context, enqueue core.EnqueueFn) error {
	t := time.NewTimer(d.delay)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return fmt.Errorf("async delay interrupted: %w", ctx.Err())
	case <-t.C:
		internal.Logger.DebugContext(ctx, "Delay Duration", "duration", time.Since(d.start))
		return enqueue(DelayerFinishedEvent{d.Id(), d.start})
	}
}

// Activate ...
func (d *asyncdelayer) Activate(ctx context.Context, evt core.Event) core.ResultDetails {
	d.start = time.Now()

	internal.Logger.DebugContext(ctx, "Returning AsyncRunning", "name", d.Name())

	return core.InitRunningResult(d.doDelay)
}

// Tick ...
func (d *asyncdelayer) Tick(ctx context.Context, evt core.Event) core.ResultDetails {
	internal.Logger.DebugContext(ctx, "Tick", "name", d.Name())

	if dfe, ok := evt.(DelayerFinishedEvent); ok {
		if dfe.TargetNodeId() == d.Id() {
			internal.Logger.DebugContext(ctx, "DelayerFinishedEvent", "name", d.Name())
			return core.Update(ctx, d.Child, evt)
		}
	}

	return core.RunningResult()
}

// Leave ...
func (d *asyncdelayer) Leave(context.Context) error {
	return nil
}

var _ core.Node = (*asyncdelayer)(nil)
