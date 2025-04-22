package coro

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"reflect"
	"runtime"
	"sync/atomic"

	"github.com/jbcpollak/greenstalk/core"
)

// node wraps a coroutine style node function as a normal Node.
type node[Blackboard any, P core.Params] struct {
	core.Leaf[Blackboard, P]
	coro     NodeFunc[Blackboard, P]
	coroNext func() (core.ResultDetails, bool)
	coroStop func()

	nextTick atomic.Pointer[Tick[Blackboard]]
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
		coro: f,
	}
	return n
}

// SimpleNode wraps [Node] for the common case of a coroutine that doesn't use a
// blackboard (i.e. uses [core.EmptyBlackboard]) and can use the function name
// as the node name.
func SimpleNode(f NodeFunc[core.EmptyBlackboard, core.BaseParams]) *node[core.EmptyBlackboard, core.BaseParams] {
	funcName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	return Node(f, core.BaseParams("coro."+funcName))
}

var (
	ErrAlreadyActivated = errors.New("already activated")
	ErrNotActivated     = errors.New("not activated")
	ErrNoResult         = errors.New("no result")
	ErrNextTooSoon      = errors.New("called next too soon")
)

// Activate implements core.Node.
func (n *node[Blackboard, P]) Activate(ctx context.Context, b Blackboard, e core.Event) core.ResultDetails {
	if n.coroNext != nil {
		return core.ErrorResult(fmt.Errorf("%s: %w", n.Params.Name(), ErrAlreadyActivated))
	}
	// push args to the coroutine
	n.nextTick.Store(&Tick[Blackboard]{ctx, b, e})

	// start the coroutine
	n.coroNext, n.coroStop = iter.Pull(n.coro(ctx, n.Params, n.next()))

	// and return its first result
	r, ok := n.coroNext()
	if !ok {
		return core.ErrorResult(ErrNoResult)
	}
	return r
}

// Tick implements core.Node.
func (n *node[Blackboard, P]) Tick(ctx context.Context, b Blackboard, e core.Event) core.ResultDetails {
	// push args to the coroutine
	n.nextTick.Store(&Tick[Blackboard]{ctx, b, e})
	if n.coroNext == nil {
		return core.ErrorResult(fmt.Errorf("%s: %w", n.Params.Name(), ErrNotActivated))
	}
	r, ok := n.coroNext()
	if !ok {
		return core.ErrorResult(ErrNoResult)
	}
	return r
}

// Leave implements core.Node.
func (n *node[Blackboard, P]) Leave(Blackboard) error {
	stop := n.coroStop
	if stop == nil {
		return fmt.Errorf("%s: %w", n.Params.Name(), ErrNotActivated)
	}
	n.coroNext, n.coroStop = nil, nil
	stop()
	return nil
}

func (n *node[Blackboard, P]) next() iter.Seq[Tick[Blackboard]] {
	return func(yield func(Tick[Blackboard]) bool) {
		var last Tick[Blackboard]
		// loop until node deactivated, which sets fNext nil
		for n.coroNext != nil {
			args := n.nextTick.Swap(nil)
			if args == nil {
				// keep the context & blackboard from the last tick if we have them
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
