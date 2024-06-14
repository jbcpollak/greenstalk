package decorator

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
)

func SwapResult[Blackboard any](swapFrom core.Status, swapTo core.Status, child core.Node[Blackboard]) core.Node[Blackboard] {
	if swapFrom != core.StatusSuccess && swapFrom != core.StatusFailure {
		panic("cannot swap from statuses other than Success or Failure")
	}

	if swapTo != core.StatusSuccess && swapTo != core.StatusFailure {
		panic("cannot swap to statuses other than Success or Failure")
	}

	base := core.NewDecorator(core.BaseParams("SwapResult"), child)
	return &swapResult[Blackboard]{base, swapFrom, swapTo}
}

type swapResult[Blackboard any] struct {
	core.Decorator[Blackboard, core.BaseParams]
	swapFrom core.Status
	swapTo   core.Status
}

func (d *swapResult[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	return d.Tick(ctx, bb, evt)
}

// Tick ...
func (d *swapResult[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	result := core.Update(ctx, d.Child, bb, evt)
	if result.Status() == d.swapFrom {
		if d.swapTo == core.StatusSuccess {
			return core.SuccessResult()
		}
		return core.FailureResult()
	}
	return result
}

// Leave ...
func (d *swapResult[Blackboard]) Leave(bb Blackboard) error {
	return nil
}
