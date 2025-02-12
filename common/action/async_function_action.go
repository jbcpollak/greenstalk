package action

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jbcpollak/greenstalk/core"
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
func AsyncFunctionAction[Blackboard any](params AsyncFunctionActionParams) *asyncFunctionAction[Blackboard] {
	base := core.NewLeaf[Blackboard](params)
	return &asyncFunctionAction[Blackboard]{Leaf: base}
}

type asyncFunctionAction[Blackboard any] struct {
	core.Leaf[Blackboard, AsyncFunctionActionParams]

	fnResult core.ResultDetails
}

func (a *asyncFunctionAction[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	return core.InitRunningResult(a.performFunction)
}

func (a *asyncFunctionAction[Blackboard]) performFunction(ctx context.Context, enqueue core.EnqueueFn) error {
	a.fnResult = a.Params.Func(ctx)
	if a.fnResult.Status() == core.StatusRunning {
		a.fnResult = core.ErrorResult(fmt.Errorf("async function returned invalid status of StatusRunning"))
	}

	return enqueue(asyncFunctionFinishedEvent{
		targetNodeId: a.Id(),
	})
}

func (a *asyncFunctionAction[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	if afe, ok := evt.(asyncFunctionFinishedEvent); ok {
		if afe.TargetNodeId() == a.Id() {
			return a.fnResult
		}
	}

	return core.RunningResult()
}

func (a *asyncFunctionAction[Blackboard]) Leave(bb Blackboard) error {
	return nil
}
