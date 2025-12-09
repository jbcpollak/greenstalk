package decorator

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/jbcpollak/greenstalk/v2"
	"github.com/jbcpollak/greenstalk/v2/common/action"
	"github.com/jbcpollak/greenstalk/v2/common/composite"
	"github.com/jbcpollak/greenstalk/v2/core"
	"github.com/jbcpollak/greenstalk/v2/internal"
	"github.com/jbcpollak/greenstalk/v2/util"
)

func TestUntilFailure(t *testing.T) {
	var wg sync.WaitGroup

	// Synchronous, so does not need to be cancelled.
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	sigChan := make(chan bool)
	defer close(sigChan)

	countChan := make(chan uint)
	defer close(countChan)

	child := action.Counter(action.CounterParams{
		BaseParams: "Counter",
		Limit:      3,
		CountChan:  countChan,
	})

	untilFailure := UntilFailure(child)

	params := action.SignallerParams[bool]{
		BaseParams: "Signaller",
		Channel:    sigChan,
		Signal:     true,
	}
	signaller := action.Signaller(params)

	testSequence := composite.Sequence(
		untilFailure,
		action.Succeed(action.SucceedParams{
			BaseParams: "Success",
		}),
		signaller,
	)

	tree, err := greenstalk.NewBehaviorTree(
		testSequence,
		greenstalk.WithVisitors(util.PrintTreeInColor),
	)
	if err != nil {
		t.Errorf("Unexpectedly got %v", err)
	}

	evt := core.DefaultEvent{}
	wg.Go(func() {
		err = tree.EventLoop(ctx, evt)
		if err != nil {
			t.Errorf("Unexpectedly got %v", err)
		}
	})

	go func() {
		for range countChan {
			// drain it
		}
	}()

	d := time.Duration(250) * time.Millisecond
	signal, timeout_err := internal.WaitForSignalOrTimeout(sigChan, d)
	if (timeout_err != nil) || !signal {
		t.Errorf("Unexpectedly got %v", signal)
	}

	cancel()
	wg.Wait()
	status := tree.Root.Result().Status()
	if status != core.StatusSuccess {
		t.Errorf("Unexpectedly got %v", status)
	}
}
