package decorator

import (
	"context"
	"fmt"
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

func TestWith(t *testing.T) {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	sigChan := make(chan bool)

	childCalled := false
	child := action.FunctionAction(action.FunctionActionParams{
		Func: func() core.ResultDetails {
			childCalled = true
			return core.SuccessResult()
		},
	})

	closeCalled := false
	closeFn := func(context.Context) error {
		closeCalled = true
		return nil
	}
	with := With(func(context.Context) (func(context.Context) error, error) {
		return closeFn, nil
	}, child)

	params := action.SignallerParams[bool]{
		BaseParams: "Signaller",
		Channel:    sigChan,
		Signal:     true,
	}
	signaller := action.Signaller(params)

	testSequence := composite.Sequence(
		with,
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
		err := tree.EventLoop(ctx, evt)
		if err != nil {
			t.Errorf("Unexpectedly got %v", err)
		}
	})

	d := time.Duration(100) * time.Millisecond

	signal, err := internal.WaitForSignalOrTimeout(sigChan, d)
	if (err != nil) || !signal {
		t.Errorf("Unexpectedly got %v", signal)
	}

	cancel()
	wg.Wait()
	if status := testSequence.Result().Status(); status != core.StatusSuccess {
		t.Errorf("Unexpectedly got %v", status)
	}

	if !childCalled {
		t.Errorf("Child was not called")
	}

	if !closeCalled {
		t.Errorf("Close was not called")
	}
}

func TestWithCloserError(t *testing.T) {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	sigChan := make(chan bool)

	childCalled := false
	child := action.FunctionAction(action.FunctionActionParams{
		Func: func() core.ResultDetails {
			childCalled = true
			return core.SuccessResult()
		},
	})

	closeCalled := false
	closeFn := func(context.Context) error {
		closeCalled = true
		return fmt.Errorf("This is an expected error")
	}

	with := With(func(context.Context) (func(context.Context) error, error) {
		return closeFn, nil
	}, child)

	params := action.SignallerParams[bool]{
		BaseParams: "Signaller",
		Channel:    sigChan,
		Signal:     true,
	}
	signaller := action.Signaller(params)

	testSequence := composite.Sequence(
		with,
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
		err := tree.EventLoop(ctx, evt)
		if err == nil {
			t.Errorf("We are expecting an error here")
		}
	})

	d := time.Duration(100) * time.Millisecond

	signal, err := internal.WaitForSignalOrTimeout(sigChan, d)
	if err == nil {
		t.Errorf("Was expecting to timeout here but got %v", signal)
	}

	cancel()
	wg.Wait()
	if status := testSequence.Result().Status(); status != core.StatusError {
		t.Errorf("Unexpectedly got %v", status)
	}

	if !childCalled {
		t.Errorf("Child was not called")
	}

	if !closeCalled {
		t.Errorf("Close was not called")
	}
}

func TestWithInitError(t *testing.T) {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	sigChan := make(chan bool)

	childCalled := false
	child := action.FunctionAction(action.FunctionActionParams{
		Func: func() core.ResultDetails {
			childCalled = true
			return core.SuccessResult()
		},
	})

	closeCalled := false

	with := With(func(context.Context) (func(context.Context) error, error) {
		return nil, fmt.Errorf("This is an error")
	}, child)

	params := action.SignallerParams[bool]{
		BaseParams: "Signaller",
		Channel:    sigChan,
		Signal:     true,
	}
	signaller := action.Signaller(params)

	testSequence := composite.Sequence(
		with,
		signaller,
	)

	tree, err := greenstalk.NewBehaviorTree(
		testSequence,
		greenstalk.WithVisitors(util.PrintTreeInColor),
	)
	if err != nil {
		t.Errorf("Should net error here %v", err)
	}

	evt := core.DefaultEvent{}
	wg.Go(func() {
		err = tree.EventLoop(ctx, evt)
		if err.Error() != "This is an error" {
			t.Errorf("Error does not have correct contents: %v", err)
		}
	})

	d := time.Duration(100) * time.Millisecond

	signal, err := internal.WaitForSignalOrTimeout(sigChan, d)
	if err == nil {
		t.Errorf("Was expecting to timeout here but got %v", signal)
	}

	cancel()
	wg.Wait()
	if status := testSequence.Result().Status(); status != core.StatusError {
		t.Errorf("Unexpectedly got %v", status)
	}

	if childCalled {
		t.Errorf("Child was unexpectedly called")
	}

	if closeCalled {
		t.Errorf("Close was unexpectedly called")
	}
}
