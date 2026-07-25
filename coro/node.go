package coro

import (
	"context"
	"iter"
	"reflect"
	"runtime"
	"strings"

	"github.com/jbcpollak/greenstalk/core"
)

// node wraps a coroutine style node function as a normal Node.
type node[Blackboard any, P core.Params] struct {
	core.Leaf[Blackboard, P]
	common[Blackboard, P]
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
	return &node[Blackboard, P]{
		Leaf:   core.NewLeaf[Blackboard](params),
		common: wrap(f, params),
	}
}

// SimpleNode wraps [Node] for the common case of a coroutine that doesn't use a
// blackboard (i.e. uses [core.EmptyBlackboard]) and can use the function name
// as the node name.
func SimpleNode(
	f NodeFunc[core.EmptyBlackboard, core.BaseParams],
) *node[core.EmptyBlackboard, core.BaseParams] {
	funcName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	// strip it down to just the leaf package name
	if _, leaf, ok := strings.Cut(funcName, "/"); ok {
		funcName = leaf
	}
	// replace problematic chars
	funcName = strings.ReplaceAll(funcName, ".", "_")
	return Node(f, core.BaseParams("coro_"+funcName))
}
