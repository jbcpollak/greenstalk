package coro

import (
	"context"
	"iter"
	"reflect"
	"runtime"

	"github.com/jbcpollak/greenstalk/core"
)

// CompositeNodeFunc is like [NodeFunc], but for coroutines implementing
// composite nodes.
type CompositeNodeFunc[Blackboard any, P core.Params] func(
	ctx context.Context,
	params P,
	children []core.Node[Blackboard],
	next iter.Seq[Tick[Blackboard]],
) iter.Seq[core.ResultDetails]

type composite[Blackboard any, P core.Params] struct {
	core.Composite[Blackboard, P]
	common[Blackboard, P]
}

// Composite wraps a [CompositeNodeFunc] to implement a [core.Node] that
// supports Walk/Visit seeing the children.
func Composite[Blackboard any, P core.Params](
	f CompositeNodeFunc[Blackboard, P],
	params P,
	children ...core.Node[Blackboard],
) *composite[Blackboard, P] {
	return &composite[Blackboard, P]{
		Composite: core.NewComposite(params, children),
		common: wrap(
			func(ctx context.Context, params P, next iter.Seq[Tick[Blackboard]]) iter.Seq[core.ResultDetails] {
				return f(ctx, params, children, next)
			},
			params,
		),
	}
}

func SimpleComposite(
	f CompositeNodeFunc[core.EmptyBlackboard, core.BaseParams],
) *composite[core.EmptyBlackboard, core.BaseParams] {
	funcName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	return Composite(f, core.BaseParams("coro."+funcName))
}
