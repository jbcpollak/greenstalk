package condition

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
)

type FunctionConditionParams[Blackboard any] struct {
	core.BaseParams
	Func func(bb Blackboard) bool
}

type FunctionConditionReturns struct{}

// FunctionCondition executes the provided function that returns a boolean, and returns Success/Failure based on that boolean value
func FunctionCondition[Blackboard any](params FunctionConditionParams[Blackboard], returns FunctionConditionReturns) *function_condition[Blackboard] {
	base := core.NewLeaf[Blackboard, FunctionConditionParams[Blackboard], FunctionConditionReturns](params, returns)
	return &function_condition[Blackboard]{Leaf: base}
}

type function_condition[Blackboard any] struct {
	core.Leaf[Blackboard, FunctionConditionParams[Blackboard], FunctionConditionReturns]
}

func (a *function_condition[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	funcResult := a.Params.Func(bb)
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
