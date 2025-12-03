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
func Delayer(params DelayerParams, child core.Node) core.Node {
	base := core.NewDecorator(params, child)

	d := &delayer{
		Decorator: base,
		delay:     params.Delay,
	}
	return d
}

// delayer ...
type delayer struct {
	core.Decorator[DelayerParams]
	delay time.Duration // delay in milliseconds
	start time.Time
}

// Activate ...
func (d *delayer) Activate(ctx context.Context, evt core.Event) core.ResultDetails {
	d.start = time.Now()

	return d.Tick(ctx, evt)
}

// Tick ...
func (d *delayer) Tick(ctx context.Context, evt core.Event) core.ResultDetails {
	if time.Since(d.start) > d.delay {
		return core.Update(ctx, d.Child, evt)
	}
	return core.RunningResult()
}

// Leave ...
func (d *delayer) Leave(context.Context) error {
	return nil
}

var _ core.Node = (*delayer)(nil)
