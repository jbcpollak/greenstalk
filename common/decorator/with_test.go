package decorator

import (
	"context"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/jbcpollak/greenstalk"
	"github.com/jbcpollak/greenstalk/common/action"
	"github.com/jbcpollak/greenstalk/common/composite"
	"github.com/jbcpollak/greenstalk/core"
	"github.com/jbcpollak/greenstalk/util"
	"github.com/rs/zerolog/log"
)

type testCloser struct {
	closeCalled *bool
}

func (t testCloser) Close() error {
	*t.closeCalled = true
	return nil
}

func TestWith(t *testing.T) {
	var wg sync.WaitGroup

	// Synchronous, so does not need to be cancelled.
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan bool)

	childCalled := new(bool)
	*childCalled = false
	child := action.FunctionAction[core.EmptyBlackboard](action.FunctionActionParams{
		Func: func() core.ResultDetails {
			*childCalled = true
			return core.SuccessResult()
		},
	})

	closeCalled := new(bool)
	*closeCalled = false
	closer := testCloser{
		closeCalled: closeCalled,
	}
	with := With(func() (io.Closer, error) {
		return closer, nil
	}, child)
	// TODO: testCloser and the whole closeCalled thing above just creates a struct with an attached method
	// Is there a more succinct way to write this?

	params := action.SignallerParams[bool]{
		BaseParams: "Signaller",

		Channel: sigChan,
		Signal:  true,
	}
	signaller := action.Signaller[core.EmptyBlackboard](params)

	var testSequence = composite.Sequence(
		with,
		signaller,
	)

	tree, err := greenstalk.NewBehaviorTree(ctx, testSequence, core.EmptyBlackboard{})
	if err != nil {
		panic(err)
	}

	evt := core.DefaultEvent{}
	wg.Add(1)
	go func() {
		tree.EventLoop(evt)
		wg.Done()
	}()

	util.PrintTreeInColor(tree.Root)

	d := time.Duration(100) * time.Millisecond

LOOP:
	for {
		select {
		case c := <-sigChan:
			log.Info().Msgf("loop is finished %v", c)

			break LOOP
		case <-time.After(d):
			t.Errorf("Timeout after delaying %v", d)
		}
	}

	cancel()
	wg.Wait()
	status := tree.Root.Result().Status()
	if status != core.StatusSuccess {
		t.Errorf("Unexpectedly got %v", status)
	}

	if !*childCalled {
		t.Errorf("Child was not called")
	}

	if !*closeCalled {
		t.Errorf("Close was not called")
	}
}
