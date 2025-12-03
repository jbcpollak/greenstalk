package action

import (
	"context"
	"fmt"

	"github.com/jbcpollak/greenstalk/v2/core"
)

type FailParams struct {
	core.BaseParams
}

func (p FailParams) Name() string {
	return "Fail" + p.BaseParams.Name()
}

// Fail returns a new fail node, which always fails in one tick.
func Fail(params FailParams) *fail {
	base := core.NewLeaf(params)
	return &fail{Leaf: base}
}

// fail ...
type fail struct {
	core.Leaf[FailParams]
}

// Activate ...
func (a *fail) Activate(ctx context.Context, evt core.Event) core.ResultDetails {
	return core.FailureResult()
}

// Tick ...
func (a *fail) Tick(ctx context.Context, evt core.Event) core.ResultDetails {
	// Should never get here
	return core.ErrorResult(fmt.Errorf("Fail node should not be ticked"))
}

// Leave ...
func (a *fail) Leave(context.Context) error {
	return nil
}
