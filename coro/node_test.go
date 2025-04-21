package coro_test

import (
	"context"
	"iter"
	"testing"

	"github.com/jbcpollak/greenstalk/core"
	"github.com/jbcpollak/greenstalk/coro"
)

func TestNode(t *testing.T) {
	ctx := t.Context()

	tree := coro.Node(
		func(
			ctx context.Context,
			next iter.Seq2[core.EmptyBlackboard, core.Event],
		) iter.Seq[core.ResultDetails] {
			return func(yield func(core.ResultDetails) bool) {
				for b, e := range next {
					t.Logf("got tick: %v %T", b, e)
					if _, ok := e.(completionEvent); ok {
						t.Log("ending on completion event")
						yield(core.SuccessResult())
						break
					}
					yield(core.RunningResult())
				}
			}
		},
		core.BaseParams("coro.TestNode"),
	)

	bb := core.EmptyBlackboard{}
	for i := range 10 {
		var e core.Event = core.DefaultEvent{}
		if i >= 9 {
			e = completionEvent{}
		}
		r := core.Update(ctx, tree, bb, e)
		t.Logf("got result %v", r)
		if s := r.Status(); s == core.StatusFailure || s == core.StatusSuccess {
			break
		}
	}
}

type completionEvent struct {
	core.DefaultEvent
}
