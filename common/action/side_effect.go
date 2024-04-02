package action

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
)

type SideEffectParams[Blackboard any] struct {
	core.BaseParams
	Func func(bb Blackboard)
}

type SideEffectReturns struct{}

// SideEffect executes the provided function when activated and returns Success
func SideEffect[Blackboard any](params SideEffectParams[Blackboard], returns SideEffectReturns) *side_effect[Blackboard] {
	base := core.NewLeaf[Blackboard](params, returns)
	return &side_effect[Blackboard]{Leaf: base}
}

type side_effect[Blackboard any] struct {
	core.Leaf[Blackboard, SideEffectParams[Blackboard], SideEffectReturns]
}

func (a *side_effect[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	a.Params.Func(bb)
	return core.StatusSuccess
}

func (a *side_effect[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	// Should never get here
	return core.StatusError
}

func (a *side_effect[Blackboard]) Leave(bb Blackboard) {}
