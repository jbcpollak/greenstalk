package action

import (
	"context"
	"fmt"

	"github.com/jbcpollak/greenstalk/core"
)

type FunctionActionParams struct {
	core.BaseParams
	Func func(ctx context.Context) core.ResultDetails
}

// FunctionAction executes the provided function when activated and returns its result. Note that the function is executed
// synchronously so it must not block or the tree becomes unresponsive. Use AsyncFunctionAction for long running functions.
func FunctionAction[Blackboard any](params FunctionActionParams) *function_action[Blackboard] {
	base := core.NewLeaf[Blackboard](params)
	return &function_action[Blackboard]{Leaf: base}
}

type function_action[Blackboard any] struct {
	core.Leaf[Blackboard, FunctionActionParams]
}

func (a *function_action[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	return a.Params.Func(ctx)
}

func (a *function_action[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	// Should never get here
	return core.ErrorResult(
		fmt.Errorf("FunctionAction node should not be ticked"),
	)
}

func (a *function_action[Blackboard]) Leave(bb Blackboard) error {
	return nil
}
