package decorator

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/jbcpollak/greenstalk"
	"github.com/jbcpollak/greenstalk/common/action"
	"github.com/jbcpollak/greenstalk/common/composite"
	"github.com/jbcpollak/greenstalk/core"
	"github.com/jbcpollak/greenstalk/internal"
	"github.com/jbcpollak/greenstalk/util"
)

func TestUntilSuccess(t *testing.T) {
	var wg sync.WaitGroup

	// Synchronous, so does not need to be cancelled.
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan bool)

	countChan := make(chan uint)

	child := action.Counter[core.EmptyBlackboard](action.CounterParams{
		BaseParams: "Counter",
		Limit:      3,
		CountChan:  countChan,
	})

	untilFailure := UntilSuccess(Inverter(child))

	params := action.SignallerParams[bool]{
		BaseParams: "Signaller",

		Channel: sigChan,
		Signal:  true,
	}
	signaller := action.Signaller[core.EmptyBlackboard](params)

	var testSequence = composite.Sequence(
		untilFailure,
		action.Succeed[core.EmptyBlackboard](action.SucceedParams{
			BaseParams: "Success",
		}),
		signaller,
	)

	tree, err := greenstalk.NewBehaviorTree(
		testSequence,
		core.EmptyBlackboard{},
		greenstalk.WithContext[core.EmptyBlackboard](ctx),
		greenstalk.WithVisitors(util.PrintTreeInColor[core.EmptyBlackboard]),
	)
	if err != nil {
		t.Errorf("Unexpectedly got %v", err)
	}

	evt := core.DefaultEvent{}
	wg.Add(1)
	go func() {
		err := tree.EventLoop(evt)
		if err != nil {
			t.Errorf("Unexpectedly got %v", err)
		}
		wg.Done()
	}()

	// Drain the countChan
	go func() {
		for {
			<-countChan
		}
	}()

	d := time.Duration(200) * time.Millisecond
	signal, err := internal.WaitForSignalOrTimeout(sigChan, d)
	if (err != nil) || !signal {
		t.Errorf("Unexpectedly got %v", signal)
	}

	cancel()
	wg.Wait()
	status := tree.Root.Result().Status()
	if status != core.StatusSuccess {
		t.Errorf("Unexpectedly got %v", status)
	}
}
