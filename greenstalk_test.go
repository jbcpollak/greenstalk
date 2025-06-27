package greenstalk

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/jbcpollak/greenstalk/core"
	"github.com/jbcpollak/greenstalk/internal"
	"github.com/jbcpollak/greenstalk/util"

	// Use dot imports to make a tree definition look nice.
	// Be careful when doing this! These packages export
	// common word identifiers such as "Fail" and "Sequence".
	. "github.com/jbcpollak/greenstalk/common/action"
	. "github.com/jbcpollak/greenstalk/common/composite"
	. "github.com/jbcpollak/greenstalk/common/decorator"
)

type TestBlackboard struct {
	id    int
	count uint
}

var n = 0

func untilTwo(status core.ResultDetails) bool {
	n++
	return n == 2
}

var synchronousRoot = Sequence(
	RepeatUntil(RepeatUntilParams{
		BaseParams: "RepeatUntilTwo",
		Until:      untilTwo,
	}, Fail[TestBlackboard](FailParams{})),
	Succeed[TestBlackboard](SucceedParams{}),
)

func TestUpdate(t *testing.T) {
	internal.Logger.Info("Testing synchronous tree...")

	// Synchronous, so does not need to be cancelled.
	ctx := context.Background()

	tree, err := NewBehaviorTree(
		synchronousRoot,
		TestBlackboard{id: 42},
		WithContext[TestBlackboard](ctx),
		WithVisitors(util.PrintTreeInColor[TestBlackboard]),
	)
	if err != nil {
		t.Errorf("Unexpectedly got %v", err)
	}

	for {
		evt := core.DefaultEvent{}
		result := tree.Update(evt)
		if result.Status() == core.StatusSuccess {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	internal.Logger.Info("Done!")
}

var countChan = make(chan uint)

var delay = 100
var asynchronousRoot = Sequence(
	// Repeater(core.Params{"n": 2}, Fail[TestBlackboard](nil, nil)),
	AsyncDelayer(
		AsyncDelayerParams{
			BaseParams: core.BaseParams("First"),
			Delay:      time.Duration(delay) * time.Millisecond,
		},
		Counter[TestBlackboard](CounterParams{
			BaseParams: "First Counter",
			Limit:      10,
			CountChan:  countChan,
		}),
	),
	AsyncDelayer(
		AsyncDelayerParams{
			BaseParams: core.BaseParams("Second"),
			Delay:      time.Duration(delay) * time.Millisecond,
		},
		Counter[TestBlackboard](CounterParams{
			BaseParams: "Second Counter",
			Limit:      10,
			CountChan:  countChan,
		}),
	),
)

func getCount(d time.Duration) (uint, bool) {
	select {
	case c := <-countChan:
		internal.Logger.Info("got count", "count", c)
		return c, true
	case <-time.After(d):
		internal.Logger.Info("Timeout after delaying", "delay", d)
		return 0, false
	}
}

func TestEventLoop(t *testing.T) {
	internal.Logger.Info("Testing asynchronous tree...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bb := TestBlackboard{id: 42, count: 0}
	tree, err := NewBehaviorTree(
		asynchronousRoot, bb,
		WithContext[TestBlackboard](ctx),
		WithVisitors(util.PrintTreeInColor[TestBlackboard]),
	)
	if err != nil {
		t.Errorf("Unexpectedly got %v", err)
	}

	evt := core.DefaultEvent{}
	go func() {
		err := tree.EventLoop(evt)
		if err != nil {
			t.Errorf("Unexpectedly got %v", err)
		}
	}()

	// Wait half the delay and verify no value sent
	first_halfway, ok := getCount(time.Duration(delay/2) * time.Millisecond)
	if ok {
		t.Errorf("Unexpectedly got count %d", first_halfway)
	} else {
		internal.Logger.Info("Halfway through first delay counter correctly is 0")
	}

	// Sleep a bit more
	first_after, ok := getCount(time.Duration(delay/2+10) * time.Millisecond)
	if !ok {
		t.Errorf("Expected to get count after delay but got a timeout")
	} else if first_after != 1 {
		t.Errorf("Expected count to be 1, got %d", first_after)
	} else {
		internal.Logger.Info("After first delay, counter is 1")
	}

	// Wait half the delay and verify value is 0
	second_halfway, ok := getCount(time.Duration(delay/2) * time.Millisecond)
	if ok {
		t.Errorf("Unexpectedly got count %d", second_halfway)
	} else {
		internal.Logger.Info("Halfway through second delay counter correctly is 1")
	}

	// Shut it _down_
	internal.Logger.Info("Shutting down...")
	cancel()

	after_cancel, ok := getCount(time.Duration(delay/2) * time.Millisecond)

	// Ensure we shut down before the second tick
	if ok {
		t.Errorf("Expected to shut down before second tick but got %d", after_cancel)
	}

	internal.Logger.Info("Done!")
}

type errorAsyncNode struct {
	core.Leaf[TestBlackboard, core.BaseParams]
	wg *sync.WaitGroup
}

func (a *errorAsyncNode) Activate(ctx context.Context, bb TestBlackboard, evt core.Event) core.ResultDetails {
	errorFunc := func(ctx context.Context, enqueue core.EnqueueFn) error {
		a.wg.Done()
		return fmt.Errorf("Expected error during tests")
	}
	return core.InitRunningResult(errorFunc)
}

func (a *errorAsyncNode) Tick(ctx context.Context, bb TestBlackboard, evt core.Event) core.ResultDetails {
	panic("Should never get ticked during tests")
}

func (a *errorAsyncNode) Leave(bb TestBlackboard) error {
	panic("Should never leave during tests")
}

func makeErrorAsyncNode(wg *sync.WaitGroup) *errorAsyncNode {
	base := core.NewLeaf[TestBlackboard](
		core.BaseParams("ErrorAsyncNode"),
	)
	return &errorAsyncNode{
		Leaf: base,
		wg:   wg,
	}
}

func TestAsyncErrorInTree(t *testing.T) {
	internal.Logger.Info("Testing handling errors returned by async functions...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bb := TestBlackboard{id: 42, count: 0}

	nodeWG := sync.WaitGroup{}
	nodeWG.Add(1)
	errorFuncNode := makeErrorAsyncNode(&nodeWG)
	tree, err := NewBehaviorTree(
		errorFuncNode, bb,
		WithContext[TestBlackboard](ctx),
		WithVisitors(util.PrintTreeInColor[TestBlackboard]),
	)
	if err != nil {
		t.Errorf("Unexpectedly got %v", err)
	}

	treeWG := sync.WaitGroup{}
	treeWG.Add(1)
	evt := core.DefaultEvent{}
	go func() {
		err := tree.EventLoop(evt)
		if err == nil || err.Error() != "Expected error during tests" {
			t.Errorf("Tree should have returned an expected error")
		}
		treeWG.Done()
	}()

	nodeWG.Wait()
	treeWG.Wait()

	cancel()

	t.Logf("Tree terminated cleanly")
}
