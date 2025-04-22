package coro_test

import (
	"context"
	"fmt"
	"iter"

	"github.com/jbcpollak/greenstalk/core"
	"github.com/jbcpollak/greenstalk/coro"
)

func sequence[Blackboard any](
	_ context.Context,
	_ core.BaseParams,
	children []core.Node[Blackboard],
	next iter.Seq[coro.Tick[Blackboard]],
) iter.Seq[core.ResultDetails] {
	return func(yield func(core.ResultDetails) bool) {
		i := 0
	EVENTS:
		for args := range next {
			for ; i < len(children); i++ {
				// TODO: ctx is wrong here
				r := core.Update(args.Ctx, children[i], args.BB, args.Event)
				if s := r.Status(); s != core.StatusSuccess {
					if !yield(r) {
						return
					}
					continue EVENTS // next event
				}
			}
			// all children succeeded, we're done
			break
		}
		yield(core.SuccessResult())
	}
}

func NewSequence[Blackboard any](
	children ...core.Node[Blackboard],
) core.Node[Blackboard] {
	return coro.Composite(sequence, core.BaseParams("Example"), children...)
}

type CounterParams struct {
	core.BaseParams
	Limit int
}

// RunningCounter is like [core.Counter], except it returns Running until it
// hits its limit and then returns Success.
type RunningCounter[Blackboard any] struct {
	core.Leaf[Blackboard, CounterParams]
	Current int
}

// Activate implements core.Node.
func (r *RunningCounter[Blackboard]) Activate(ctx context.Context, bb Blackboard, e core.Event) core.ResultDetails {
	r.Current = 0
	return r.Tick(ctx, bb, e)
}

// Tick implements core.Node.
func (r *RunningCounter[Blackboard]) Tick(ctx context.Context, bb Blackboard, e core.Event) core.ResultDetails {
	r.Current++
	if r.Current < r.Params.Limit {
		return core.RunningResult()
	}
	return core.SuccessResult()
}

// Leave implements core.Node.
func (r *RunningCounter[Blackboard]) Leave(Blackboard) error {
	return nil
}

func ExampleRunningCounter() {
	c1 := &RunningCounter[core.EmptyBlackboard]{
		Leaf: core.NewLeaf[core.EmptyBlackboard](CounterParams{
			BaseParams: core.BaseParams("c1"),
			Limit:      5,
		}),
	}
	c2 := &RunningCounter[core.EmptyBlackboard]{
		Leaf: core.NewLeaf[core.EmptyBlackboard](CounterParams{
			BaseParams: core.BaseParams("c2"),
			Limit:      6,
		}),
	}
	s := NewSequence(c1, c2)

	bb := core.EmptyBlackboard{}
	// 10 events, because the 5th event both finishes c1 and starts c2
	for i := range 10 {
		r := core.Update(context.TODO(), s, bb, core.DefaultEvent{})
		fmt.Printf("event %d, result=%v, counts=%d,%d\n", i, r.Status(), c1.Current, c2.Current)
	}

	// make sure the composite structure works correctly
	s.Walk(func(node core.Walkable[core.EmptyBlackboard], level int) {
		fmt.Printf("node: %s at %d\n", node.Name(), level)
	}, 0)

	// Output: event 0, result=3, counts=1,0
	// event 1, result=3, counts=2,0
	// event 2, result=3, counts=3,0
	// event 3, result=3, counts=4,0
	// event 4, result=3, counts=5,1
	// event 5, result=3, counts=5,2
	// event 6, result=3, counts=5,3
	// event 7, result=3, counts=5,4
	// event 8, result=3, counts=5,5
	// event 9, result=1, counts=5,6
	// node: Example at 0
	// node: c1 at 1
	// node: c2 at 1
}
