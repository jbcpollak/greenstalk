package coro

import (
	"context"
	"errors"
	"iter"
	"sync"

	"github.com/jbcpollak/greenstalk/core"
)

// node wraps a coroutine style node function as a normal Node.
type node[Blackboard any, P core.Params] struct {
	core.Leaf[Blackboard, P]
	f     NodeFunc[Blackboard]
	fNext func() (core.ResultDetails, bool)
	fStop func()

	nMu sync.Mutex
	nB  *Blackboard
	nE  core.Event
}

type NodeFunc[Blackboard any] func(
	ctx context.Context,
	next iter.Seq2[Blackboard, core.Event],
) iter.Seq[core.ResultDetails]

func Node[Blackboard any, P core.Params](
	f NodeFunc[Blackboard],
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
)

// Activate implements core.Node.
func (n *node[Blackboard, P]) Activate(ctx context.Context, b Blackboard, e core.Event) core.ResultDetails {
	if n.fNext != nil {
		return core.ErrorResult(ErrAlreadyActivated)
	}
	// push b,e to the coroutine
	n.nMu.Lock()
	n.nB, n.nE = &b, e
	n.nMu.Unlock()

	// start the coroutine
	n.fNext, n.fStop = iter.Pull(n.f(ctx, n.next()))

	// and return its first result
	r, ok := n.fNext()
	if !ok {
		return core.ErrorResult(ErrNoResult)
	}
	return r
}

// Tick implements core.Node.
func (n *node[Blackboard, P]) Tick(ctx context.Context, b Blackboard, e core.Event) core.ResultDetails {
	// push b,e to the coroutine
	n.nMu.Lock()
	n.nB, n.nE = &b, e
	n.nMu.Unlock()
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

func (n *node[Blackboard, P]) next() iter.Seq2[Blackboard, core.Event] {
	return func(yield func(Blackboard, core.Event) bool) {
		for {
			n.nMu.Lock()
			if n.fNext == nil {
				n.nMu.Unlock()
				break
			}
			if n.nB == nil || n.nE == nil {
				n.nMu.Unlock()
				panic("called next too soon")
			}
			b, e := *n.nB, n.nE
			n.nB, n.nE = nil, nil
			n.nMu.Unlock()
			if !yield(b, e) {
				break
			}
		}
	}
}
