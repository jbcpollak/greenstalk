package coro

import (
	"context"
	"fmt"
	"iter"
	"sync/atomic"

	"github.com/jbcpollak/greenstalk/core"
)

type common[Blackboard any, P core.Params] struct {
	coro     NodeFunc[Blackboard, P]
	params   P
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

func wrap[Blackboard any, P core.Params](
	f NodeFunc[Blackboard, P],
	params P,
) common[Blackboard, P] {
	return common[Blackboard, P]{coro: f, params: params}
}

// Activate implements core.Node.
func (n *common[Blackboard, P]) Activate(ctx context.Context, b Blackboard, e core.Event) core.ResultDetails {
	if n.coroNext != nil {
		return core.ErrorResult(fmt.Errorf("%s: %w", n.params.Name(), ErrAlreadyActivated))
	}
	// push args to the coroutine
	n.nextTick.Store(&Tick[Blackboard]{ctx, b, e})

	// start the coroutine
	n.coroNext, n.coroStop = iter.Pull(n.coro(ctx, n.params, n.next()))

	// and return its first result
	r, ok := n.coroNext()
	if !ok {
		return core.ErrorResult(ErrNoResult)
	}
	return r
}

// Tick implements core.Node.
func (n *common[Blackboard, P]) Tick(ctx context.Context, b Blackboard, e core.Event) core.ResultDetails {
	// push args to the coroutine
	n.nextTick.Store(&Tick[Blackboard]{ctx, b, e})
	if n.coroNext == nil {
		return core.ErrorResult(fmt.Errorf("%s: %w", n.params.Name(), ErrNotActivated))
	}
	r, ok := n.coroNext()
	if !ok {
		return core.ErrorResult(ErrNoResult)
	}
	return r
}

// Leave implements core.Node.
func (n *common[Blackboard, P]) Leave(Blackboard) error {
	stop := n.coroStop
	if stop == nil {
		return fmt.Errorf("%s: %w", n.params.Name(), ErrNotActivated)
	}
	n.coroNext, n.coroStop = nil, nil
	stop()
	return nil
}

func (n *common[Blackboard, P]) next() iter.Seq[Tick[Blackboard]] {
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
