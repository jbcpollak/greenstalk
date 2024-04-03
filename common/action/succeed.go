package action

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
)

type SucceedParams struct {
	core.BaseParams
}

func (p SucceedParams) Name() string {
	return "Succeed" + p.BaseParams.Name()
}

// Succeed returns a new succeed node, which always succeeds in one tick.
func Succeed[Blackboard any](params SucceedParams) core.Node[Blackboard] {
	base := core.NewLeaf[Blackboard](params)
	return &succeed[Blackboard]{Leaf: base}
}

// succeed ...
type succeed[Blackboard any] struct {
	core.Leaf[Blackboard, SucceedParams]
}

// Activate ...
func (a *succeed[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	return core.StatusSuccess
}

// Tick ...
func (a *succeed[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	// Should never get here
	return core.StatusError
}

// Leave ...
func (a *succeed[Blackboard]) Leave(bb Blackboard) {}
