package action

import (
	"context"
	"fmt"

	"github.com/jbcpollak/greenstalk/core"
)

type FailParams struct {
	core.BaseParams
}

func (p FailParams) Name() string {
	return "Fail" + p.BaseParams.Name()
}

// Fail returns a new fail node, which always fails in one tick.
func Fail[Blackboard any](params FailParams) *fail[Blackboard] {
	base := core.NewLeaf[Blackboard](params)
	return &fail[Blackboard]{Leaf: base}
}

// fail ...
type fail[Blackboard any] struct {
	core.Leaf[Blackboard, FailParams]
}

// Activate ...
func (a *fail[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	return core.FailureResult()
}

// Tick ...
func (a *fail[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	// Should never get here
	return core.ErrorResult(fmt.Errorf("Fail node should not be ticked"))
}

// Leave ...
func (a *fail[Blackboard]) Leave(bb Blackboard) error {
	return nil
}
