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
func Succeed(params SucceedParams) core.Node {
	base := core.NewLeaf(params)
	return &succeed{Leaf: base}
}

// succeed ...
type succeed struct {
	core.Leaf[SucceedParams]
}

// Activate ...
func (a *succeed) Activate(ctx context.Context, evt core.Event) core.ResultDetails {
	return core.SuccessResult()
}

// Tick ...
func (a *succeed) Tick(ctx context.Context, evt core.Event) core.ResultDetails {
	return core.ErrorResult(
		fmt.Errorf("Succeed node should not be ticked"),
	)
}

// Leave ...
func (a *succeed) Leave(context.Context) error {
	return nil
}

var _ core.Node = (*succeed)(nil)
