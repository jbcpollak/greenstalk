package action

import (
	"context"

	"github.com/jbcpollak/greenstalk/v2/core"
	"github.com/jbcpollak/greenstalk/v2/internal"
)

type CounterParams struct {
	core.BaseParams

	Limit     uint
	CountChan chan uint
}

// Counter increments a counter on the blackboard until it reaches a certain value.
func Counter(params CounterParams) core.Node {
	base := core.NewLeaf(params)
	return &counter{Leaf: base, currentValue: 0}
}

type counter struct {
	core.Leaf[CounterParams]

	currentValue uint
}

// Activate ...
func (a *counter) Activate(ctx context.Context, evt core.Event) core.ResultDetails {
	return a.Tick(ctx, evt)
}

func (a *counter) Tick(ctx context.Context, evt core.Event) core.ResultDetails {
	internal.Logger.Info("Incrementing count", "name", a.Name())
	a.currentValue++
	a.Params.CountChan <- a.currentValue

	if a.currentValue < a.Params.Limit {
		return core.SuccessResult()
	}
	return core.FailureResult()
}

// Leave ...
func (a *counter) Leave(context.Context) error {
	return nil
}

var _ core.Node = (*counter)(nil)
