package action

import (
	"context"
	"fmt"

	"github.com/jbcpollak/greenstalk/core"
)

type FunctionActionParams struct {
	core.BaseParams
	Func func() core.ResultDetails
}

// FunctionAction executes the provided function when activated and returns its result. Note that the function is executed
// synchronously so it must not block or the tree becomes unresponsive. Use AsyncFunctionAction for long running functions.
func FunctionAction(params FunctionActionParams) *function_action {
	base := core.NewLeaf(params)
	return &function_action{Leaf: base}
}

type function_action struct {
	core.Leaf[FunctionActionParams]
}

func (a *function_action) Activate(ctx context.Context, evt core.Event) core.ResultDetails {
	return a.Params.Func()
}

func (a *function_action) Tick(ctx context.Context, evt core.Event) core.ResultDetails {
	// Should never get here
	return core.ErrorResult(
		fmt.Errorf("FunctionAction node should not be ticked"),
	)
}

func (a *function_action) Leave(context.Context) error {
	return nil
}

var _ core.Node = (*function_action)(nil)
