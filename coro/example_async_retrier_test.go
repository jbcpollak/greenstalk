package coro_test

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"
	"github.com/jbcpollak/greenstalk"
	"github.com/jbcpollak/greenstalk/core"
	"github.com/jbcpollak/greenstalk/coro"
)

type AsyncParams[T any] struct {
	core.BaseParams
	F func() (T, error)
	// a real implementation wouldn't stop the BT loop on success
	Done context.CancelFunc
}

func AsyncRetrier[Blackboard any, T any](
	f func() (T, error),
	p core.BaseParams,
	done context.CancelFunc,
) core.Node[Blackboard] {
	return coro.Node[Blackboard](
		runAsyncRetrier,
		AsyncParams[T]{
			BaseParams: p,
			F:          f,
			Done:       done,
		},
	)
}

type AsyncResultEvent[T any] struct {
	// targetNodeId uuid.UUID
	Result T
	Err    error
}

func (e AsyncResultEvent[T]) TargetNodeId() uuid.UUID {
	// return e.targetNodeId
	return uuid.Nil
}

type DelayResultEvent struct{ core.DefaultEvent }

func runAsyncRetrier[Blackboard any, T any](
	ctx context.Context,
	params AsyncParams[T],
	events iter.Seq2[Blackboard, core.Event],
) iter.Seq[core.ResultDetails] {
	return func(yield func(core.ResultDetails) bool) {
		// this could also be written as a for loop over `events`, this is just
		// an example of how things can be done with pull wrappers when that
		// makes things easier.
		next, stop := iter.Pull2(events)
		defer stop()

		_, _, ok := next()
		if !ok {
			panic("WAT")
		}
		var do core.RunningFn = func(ctx context.Context, enqueue core.EnqueueFn) error {
			var evt core.Event
			if res, err := params.F(); err != nil {
				evt = AsyncResultEvent[T]{Err: err}
			} else {
				evt = AsyncResultEvent[T]{Result: res}
			}
			return enqueue(evt)
		}
		var delay core.RunningFn = func(ctx context.Context, enqueue core.EnqueueFn) error {
			time.Sleep(time.Millisecond)
			return enqueue(DelayResultEvent{})
		}
		if !yield(core.InitRunningResult(do)) {
			panic("WAT")
		}

		for {
			_, evt, ok := next()
			if !ok {
				break
			}
			switch evt := evt.(type) {
			case DelayResultEvent:
				fmt.Printf("retry delay elapsed, starting again\n")
				if !yield(core.InitRunningResult(do)) {
					panic("WAT")
				}
			case AsyncResultEvent[T]:
				if evt.Err != nil {
					fmt.Printf("async routine failed, delaying: %v\n", evt.Err)
					if !yield(core.InitRunningResult(delay)) {
						panic("WAT")
					}
				} else {
					fmt.Printf("async routine completed: %v\n", evt.Result)
					params.Done()
					yield(core.SuccessResult())
					return
				}
			}
		}
	}
}

var predictableRand = rand.New(rand.NewPCG(0, 2))

func asyncFaker() (int, error) {
	time.Sleep(time.Millisecond)
	if predictableRand.IntN(2) == 0 {
		return predictableRand.Int(), nil
	} else {
		return -1, errors.New("fake error")
	}
}

func ExampleAsyncRetrier() {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	retrier := AsyncRetrier[core.EmptyBlackboard](
		asyncFaker,
		core.BaseParams("AsyncFaker"),
		cancel,
	)
	bb := core.EmptyBlackboard{}
	bt, err := greenstalk.NewBehaviorTree(
		retrier,
		bb,
		greenstalk.WithContext[core.EmptyBlackboard](ctx),
	)
	if err != nil {
		panic(err)
	}

	if err := bt.EventLoop(core.DefaultEvent{}); err != nil {
		panic(err)
	}

	// expected output is based on the `predictableRand` behavior

	// Output: async routine failed, delaying: fake error
	// retry delay elapsed, starting again
	// async routine failed, delaying: fake error
	// retry delay elapsed, starting again
	// async routine failed, delaying: fake error
	// retry delay elapsed, starting again
	// async routine failed, delaying: fake error
	// retry delay elapsed, starting again
	// async routine failed, delaying: fake error
	// retry delay elapsed, starting again
	// async routine completed: 2976040640374945586
}
