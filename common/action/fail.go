package action

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
)

type FailParams struct {
	core.BaseParams
}

func (p FailParams) Name() string {
	return "Fail" + p.BaseParams.Name()
}

type FailReturns struct{}

// Fail returns a new fail node, which always fails in one tick.
func Fail[Blackboard any](params FailParams, returns FailReturns) *fail[Blackboard] {
	base := core.NewLeaf[Blackboard](params, returns)
	return &fail[Blackboard]{Leaf: base}
}

// fail ...
type fail[Blackboard any] struct {
	core.Leaf[Blackboard, FailParams, FailReturns]
}

// Activate ...
func (a *fail[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	return core.StatusFailure
}

// Tick ...
func (a *fail[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	// Should never get here
	return core.StatusError
}

// Leave ...
func (a *fail[Blackboard]) Leave(bb Blackboard) {}
