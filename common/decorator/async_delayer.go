package decorator

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jbcpollak/greenstalk/core"
	"github.com/jbcpollak/greenstalk/internal"
)

type AsyncDelayerParams struct {
	core.BaseParams

	Delay time.Duration
}

// AsyncDelayer ...
func AsyncDelayer[Blackboard any](params AsyncDelayerParams, child core.Node[Blackboard]) core.Node[Blackboard] {

	base := core.NewDecorator(params, child)

	d := &asyncdelayer[Blackboard]{
		Decorator: base,
		delay:     params.Delay,
	}
	return d
}

// delayer ...
type asyncdelayer[Blackboard any] struct {
	core.Decorator[Blackboard, AsyncDelayerParams]
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

func (d *asyncdelayer[Blackboard]) doDelay(ctx context.Context, enqueue core.EnqueueFn) error {
	t := time.NewTimer(d.delay)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return fmt.Errorf("async delay interrupted: %w", ctx.Err())
	case <-t.C:
		internal.Logger.Info("Delay Duration", "duration", time.Since(d.start))
		return enqueue(DelayerFinishedEvent{d.Id(), d.start})
	}
}

// Activate ...
func (d *asyncdelayer[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	d.start = time.Now()

	internal.Logger.Info("Returning AsyncRunning", "name", d.Name())

	return core.InitRunningResult(d.doDelay)
}

// Tick ...
func (d *asyncdelayer[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	internal.Logger.Info("Tick", "name", d.Name())

	if dfe, ok := evt.(DelayerFinishedEvent); ok {
		if dfe.TargetNodeId() == d.Id() {
			internal.Logger.Info("DelayerFinishedEvent", "name", d.Name())
			return core.Update(ctx, d.Child, bb, evt)
		}
	}

	return core.RunningResult()
}

// Leave ...
func (d *asyncdelayer[Blackboard]) Leave(bb Blackboard) error {
	return nil
}
