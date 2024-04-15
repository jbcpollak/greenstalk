package action

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
	"github.com/rs/zerolog/log"
)

type CounterParams struct {
	core.BaseParams

	Limit     uint
	CountChan chan uint
}

// Counter increments a counter on the blackboard until it reaches a certain value.
func Counter[Blackboard any](params CounterParams) core.Node[Blackboard] {
	base := core.NewLeaf[Blackboard](params)
	return &counter[Blackboard]{Leaf: base, currentValue: 0}
}

type counter[Blackboard any] struct {
	core.Leaf[Blackboard, CounterParams]

	currentValue uint
}

// Activate ...
func (a *counter[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	return a.Tick(ctx, bb, evt)
}

func (a *counter[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	log.Info().Msgf("%s: Incrementing count", a.Name())
	a.currentValue++
	a.Params.CountChan <- a.currentValue

	if a.currentValue < a.Params.Limit {
		return core.StatusSuccess
	}
	return core.StatusFailure
}

// Leave ...
func (a *counter[Blackboard]) Leave(bb Blackboard) error {
	return nil
}
