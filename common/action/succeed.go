package action

import (
	"context"
	"fmt"

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
func (a *succeed[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	return core.SuccessResult()
}

// Tick ...
func (a *succeed[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	return core.ErrorResult(
		fmt.Errorf("Succeed node should not be ticked"),
	)
}

// Leave ...
func (a *succeed[Blackboard]) Leave(bb Blackboard) error {
	return nil
}
