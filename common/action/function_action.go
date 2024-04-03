package action

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
)

type FunctionActionParams struct {
	core.BaseParams
	Func func()
}

// FunctionAction executes the provided function when activated and returns Success
func FunctionAction[Blackboard any](params FunctionActionParams) *function_action[Blackboard] {
	base := core.NewLeaf[Blackboard](params)
	return &function_action[Blackboard]{Leaf: base}
}

type function_action[Blackboard any] struct {
	core.Leaf[Blackboard, FunctionActionParams]
}

func (a *function_action[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	a.Params.Func()
	return core.StatusSuccess
}

func (a *function_action[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	// Should never get here
	return core.StatusError
}

func (a *function_action[Blackboard]) Leave(bb Blackboard) {}
