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

	ticks, completions := 0, 0
	counter := func(
		_ context.Context,
		_ core.BaseParams,
		next iter.Seq[coro.Tick[core.EmptyBlackboard]],
	) iter.Seq[core.ResultDetails] {
		return func(yield func(core.ResultDetails) bool) {
			for args := range next {
				t.Logf("got tick: %v %T", args.BB, args.Event)
				ticks++
				if _, ok := args.Event.(completionEvent); ok {
					t.Log("ending on completion event")
					completions++
					yield(core.SuccessResult())
					// if we don't break here, `next` will still end
					break
				} else {
					yield(core.RunningResult())
				}
			}
			t.Log("coro node complete")
		}
	}
	tree := coro.SimpleNode(counter)

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

	if ticks != 10 {
		t.Errorf("expect ticks=10, got %d", ticks)
	}
	if completions != 1 {
		t.Errorf("Expect completions=1, got %d", completions)
	}
}

type completionEvent struct {
	core.DefaultEvent
}
