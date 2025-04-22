package coro

import (
	"context"
	"errors"
	"iter"
	"sync/atomic"

	"github.com/jbcpollak/greenstalk/core"
)

// node wraps a coroutine style node function as a normal Node.
type node[Blackboard any, P core.Params] struct {
	core.Leaf[Blackboard, P]
	f     NodeFunc[Blackboard, P]
	fNext func() (core.ResultDetails, bool)
	fStop func()

	nextArgs atomic.Pointer[Tick[Blackboard]]
}

// Tick represents the arguments to a tick. In a normal node these would be
// passed as arguments. In a coroutine node they are yielded from the `next`
// iterator.
type Tick[Blackboard any] struct {
	Ctx   context.Context
	BB    Blackboard
	Event core.Event
}

// A NodeFunc is used to implement a [core.Node] as a coroutine.
//
// On activation, the function is called to start the iteration. It MUST support
// multiple iteration so that it can be re-activated after completing
// previously. The `next` parameter provides an iterator that will yield the
// blackboard & event objects from each update of the tree.
//
// The function must strictly alternate between retrieving (blackboard, event)
// pairs from `next` and yielding results. Once it yields a Success or Failure
// result, it will be presumed done and will be allowed to complete any
// "trailer" logic in the function, and iteration from `next` will end.
//
// If the function attempts to pull more than one value from `next` without
// yielding a result in between, it will get an empty Blackboard and a
// [core.ErrorEvent] wrapping [ErrNextTooSoon].
//
// If the function ends without yielding its final result, the node will end
// with an Error result wrapping [ErrNoResult].
//
// The `ctx` parameter is only valid during activation, i.e. up until the
// NodeFunc yields its first result. After that the Ctx yielded from `next` must
// be used for the duration of each tick.
type NodeFunc[Blackboard any, P core.Params] func(
	ctx context.Context,
	params P,
	next iter.Seq[Tick[Blackboard]],
) iter.Seq[core.ResultDetails]

// Node wraps a [NodeFunc] to implement a [core.Node]. See [NodeFunc] for
// details.
func Node[Blackboard any, P core.Params](
	f NodeFunc[Blackboard, P],
	params P,
) *node[Blackboard, P] {
	n := &node[Blackboard, P]{
		Leaf: core.NewLeaf[Blackboard](params),
		f:    f,
	}
	return n
}

var (
	ErrAlreadyActivated = errors.New("already activated")
	ErrNoResult         = errors.New("no result")
	ErrNextTooSoon      = errors.New("called next too soon")
)

// Activate implements core.Node.
func (n *node[Blackboard, P]) Activate(ctx context.Context, b Blackboard, e core.Event) core.ResultDetails {
	if n.fNext != nil {
		return core.ErrorResult(ErrAlreadyActivated)
	}
	// push args to the coroutine
	n.nextArgs.Store(&Tick[Blackboard]{ctx, b, e})

	// start the coroutine
	n.fNext, n.fStop = iter.Pull(n.f(ctx, n.Params, n.next()))

	// and return its first result
	r, ok := n.fNext()
	if !ok {
		return core.ErrorResult(ErrNoResult)
	}
	return r
}

// Tick implements core.Node.
func (n *node[Blackboard, P]) Tick(ctx context.Context, b Blackboard, e core.Event) core.ResultDetails {
	// push args to the coroutine
	n.nextArgs.Store(&Tick[Blackboard]{ctx, b, e})
	r, ok := n.fNext()
	if !ok {
		return core.ErrorResult(ErrNoResult)
	}
	return r
}

// Leave implements core.Node.
func (n *node[Blackboard, P]) Leave(Blackboard) error {
	stop := n.fStop
	n.fNext, n.fStop = nil, nil
	stop()
	return nil
}

func (n *node[Blackboard, P]) next() iter.Seq[Tick[Blackboard]] {
	return func(yield func(Tick[Blackboard]) bool) {
		var last Tick[Blackboard]
		// loop until node deactivated, which sets fNext nil
		for n.fNext != nil {
			args := n.nextArgs.Swap(nil)
			if args == nil {
				last.Event = core.ErrorEvent{Err: ErrNextTooSoon}
				yield(last)
				break
			}
			last = *args
			if !yield(last) {
				break
			}
		}
	}
}
