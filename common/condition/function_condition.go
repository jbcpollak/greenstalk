package condition

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
)

type FunctionConditionParams struct {
	core.BaseParams
	Func func() bool
}

type FunctionConditionReturns struct{}

// FunctionCondition executes the provided function that returns a boolean, and returns Success/Failure based on that boolean value
func FunctionCondition[Blackboard any](params FunctionConditionParams, returns FunctionConditionReturns) *function_condition[Blackboard] {
	base := core.NewLeaf[Blackboard, FunctionConditionParams, FunctionConditionReturns](params, returns)
	return &function_condition[Blackboard]{Leaf: base}
}

type function_condition[Blackboard any] struct {
	core.Leaf[Blackboard, FunctionConditionParams, FunctionConditionReturns]
}

func (a *function_condition[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	funcResult := a.Params.Func()
	if funcResult {
		return core.StatusSuccess
	} else {
		return core.StatusFailure
	}
}

func (a *function_condition[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	// Should never get here
	return core.StatusError
}

func (a *function_condition[Blackboard]) Leave(bb Blackboard) {}
