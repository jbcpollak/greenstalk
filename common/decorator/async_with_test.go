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

func TestWithAsync(t *testing.T) {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	sigChan := make(chan bool)

	enterCalled := false
	enterFunc := func(context.Context) error {
		enterCalled = true
		return nil
	}

	childCalled := false
	child := action.AsyncFunctionAction(action.AsyncFunctionActionParams{
		BaseParams: core.BaseParams("child"),
		Func: func(ctx context.Context) core.ResultDetails {
			childCalled = true
			return core.SuccessResult()
		},
	})

	exitCalled := false
	exitFunc := func(context.Context) error {
		exitCalled = true
		return nil
	}

	with := WithAsync(enterFunc, exitFunc, child)

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

	if !enterCalled {
		t.Errorf("Enter was not called")
	}

	if !childCalled {
		t.Errorf("Child was not called")
	}

	if !exitCalled {
		t.Errorf("Exit was not called")
	}
}

func TestWithAsyncExitError(t *testing.T) {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	sigChan := make(chan bool)

	enterCalled := false
	enterFunc := func(context.Context) error {
		enterCalled = true
		return nil
	}

	childCalled := false
	child := action.AsyncFunctionAction(action.AsyncFunctionActionParams{
		BaseParams: core.BaseParams("child"),
		Func: func(ctx context.Context) core.ResultDetails {
			childCalled = true
			return core.SuccessResult()
		},
	})

	exitCalled := false
	exitFunc := func(context.Context) error {
		exitCalled = true
		return fmt.Errorf("This is an expected error")
	}

	with := WithAsync(enterFunc, exitFunc, child)

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

	if !enterCalled {
		t.Errorf("Enter was not called")
	}

	if !childCalled {
		t.Errorf("Child was not called")
	}

	if !exitCalled {
		t.Errorf("Exit was not called")
	}
}

func TestWithAsyncEnterError(t *testing.T) {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	sigChan := make(chan bool)

	enterCalled := false
	enterFunc := func(context.Context) error {
		enterCalled = true
		return fmt.Errorf("This is an error")
	}

	childCalled := false
	child := action.AsyncFunctionAction(action.AsyncFunctionActionParams{
		BaseParams: core.BaseParams("child"),
		Func: func(ctx context.Context) core.ResultDetails {
			childCalled = true
			return core.SuccessResult()
		},
	})

	exitCalled := false
	exitFunc := func(context.Context) error {
		exitCalled = true
		return nil
	}

	with := WithAsync(enterFunc, exitFunc, child)

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
		t.Errorf("Should not error here %v", err)
	}

	evt := core.DefaultEvent{}
	wg.Go(func() {
		err = tree.EventLoop(ctx, evt)
		if err == nil {
			t.Error("Should have errored here")
		} else if err.Error() != "This is an error" {
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

	if !enterCalled {
		t.Errorf("Enter was not called")
	}

	if childCalled {
		t.Errorf("Child was unexpectedly called")
	}

	if exitCalled {
		t.Errorf("Exit was unexpectedly called")
	}
}

func TestWithAsyncChildFailure(t *testing.T) {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	sigChan := make(chan bool)

	enterCalled := false
	enterFunc := func(context.Context) error {
		enterCalled = true
		return nil
	}

	childCalled := false
	child := action.AsyncFunctionAction(action.AsyncFunctionActionParams{
		BaseParams: core.BaseParams("child"),
		Func: func(ctx context.Context) core.ResultDetails {
			childCalled = true
			return core.FailureResult()
		},
	})

	exitCalled := false
	exitFunc := func(context.Context) error {
		exitCalled = true
		return nil
	}

	with := WithAsync(enterFunc, exitFunc, child)

	params := action.SignallerParams[bool]{
		BaseParams: "Signaller",
		Channel:    sigChan,
		Signal:     true,
	}
	signaller := action.Signaller(params)

	testSequence := composite.Selector(
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
		t.Errorf("Unexpectedly got signal=%v, err=%v", signal, err)
	}

	cancel()
	wg.Wait()
	if status := testSequence.Result().Status(); status != core.StatusSuccess {
		t.Errorf("Unexpectedly got %v", status)
	}

	// The with node should return failure (from child)
	if status := with.Result().Status(); status != core.StatusFailure {
		t.Errorf("Expected failure but got %v", status)
	}

	if !enterCalled {
		t.Errorf("Enter was not called")
	}

	if !childCalled {
		t.Errorf("Child was not called")
	}

	if !exitCalled {
		t.Errorf("Exit was not called")
	}
}
