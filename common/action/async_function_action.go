package action

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jbcpollak/greenstalk/v2/core"
)

type AsyncFunctionActionParams struct {
	core.BaseParams
	Func func(ctx context.Context) core.ResultDetails
}

type asyncFunctionFinishedEvent struct {
	targetNodeId uuid.UUID
}

func (e asyncFunctionFinishedEvent) TargetNodeId() uuid.UUID {
	return e.targetNodeId
}

// Same as FunctionAction but the function is executed in a separate goroutine. Returns the same status as the function,
// except StatusRunning which return an ErrorResult.
func AsyncFunctionAction(params AsyncFunctionActionParams) *asyncFunctionAction {
	base := core.NewLeaf(params)
	return &asyncFunctionAction{Leaf: base}
}

type asyncFunctionAction struct {
	core.Leaf[AsyncFunctionActionParams]

	fnResult core.ResultDetails
}

func (a *asyncFunctionAction) Activate(ctx context.Context, evt core.Event) core.ResultDetails {
	return core.InitRunningResult(a.performFunction)
}

func (a *asyncFunctionAction) performFunction(ctx context.Context, enqueue core.EnqueueFn) error {
	a.fnResult = a.Params.Func(ctx)
	if a.fnResult.Status() == core.StatusRunning {
		a.fnResult = core.ErrorResult(fmt.Errorf("async function returned invalid status of StatusRunning"))
	}

	return enqueue(asyncFunctionFinishedEvent{
		targetNodeId: a.Id(),
	})
}

func (a *asyncFunctionAction) Tick(ctx context.Context, evt core.Event) core.ResultDetails {
	if afe, ok := evt.(asyncFunctionFinishedEvent); ok {
		if afe.TargetNodeId() == a.Id() {
			return a.fnResult
		}
	}

	return core.RunningResult()
}

func (a *asyncFunctionAction) Leave(context.Context) error {
	return nil
}

var _ core.Node = (*asyncFunctionAction)(nil)
