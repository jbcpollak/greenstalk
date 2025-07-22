package decorator

import (
	"context"
	"fmt"
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

func TestWith(t *testing.T) {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan bool)

	childCalled := false
	child := action.FunctionAction[core.EmptyBlackboard](action.FunctionActionParams{
		Func: func() core.ResultDetails {
			childCalled = true
			return core.SuccessResult()
		},
	})

	closeCalled := false
	closeFn := func() error {
		closeCalled = true
		return nil
	}
	with := With(func() (func() error, error) {
		return closeFn, nil
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

	d := time.Duration(100) * time.Millisecond

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

	if !childCalled {
		t.Errorf("Child was not called")
	}

	if !closeCalled {
		t.Errorf("Close was not called")
	}
}

func TestWithCloserError(t *testing.T) {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan bool)

	childCalled := false
	child := action.FunctionAction[core.EmptyBlackboard](action.FunctionActionParams{
		Func: func() core.ResultDetails {
			childCalled = true
			return core.SuccessResult()
		},
	})

	closeCalled := false
	closeFn := func() error {
		closeCalled = true
		return fmt.Errorf("This is an expected error")
	}

	with := With(func() (func() error, error) {
		return closeFn, nil
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
		if err == nil {
			t.Errorf("We are expecting an error here")
		}
		wg.Done()
	}()

	d := time.Duration(100) * time.Millisecond

	signal, err := internal.WaitForSignalOrTimeout(sigChan, d)
	if err == nil {
		t.Errorf("Was expecting to timeout here but got %v", signal)
	}

	cancel()
	wg.Wait()
	status := tree.Root.Result().Status()
	if status != core.StatusError {
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

	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan bool)

	childCalled := false
	child := action.FunctionAction[core.EmptyBlackboard](action.FunctionActionParams{
		Func: func() core.ResultDetails {
			childCalled = true
			return core.SuccessResult()
		},
	})

	closeCalled := false

	with := With(func() (func() error, error) {
		return nil, fmt.Errorf("This is an error")
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

	tree, err := greenstalk.NewBehaviorTree(
		testSequence,
		core.EmptyBlackboard{},
		greenstalk.WithContext[core.EmptyBlackboard](ctx),
		greenstalk.WithVisitors(util.PrintTreeInColor[core.EmptyBlackboard]),
	)
	if err != nil {
		t.Errorf("Should net error here %v", err)
	}

	evt := core.DefaultEvent{}
	wg.Add(1)
	go func() {
		err = tree.EventLoop(evt)
		if err.Error() != "This is an error" {
			t.Errorf("Error does not have correct contents: %v", err)
		}
		wg.Done()
	}()

	d := time.Duration(100) * time.Millisecond

	signal, err := internal.WaitForSignalOrTimeout(sigChan, d)
	if err == nil {
		t.Errorf("Was expecting to timeout here but got %v", signal)
	}

	cancel()
	wg.Wait()
	status := tree.Root.Result().Status()
	if status != core.StatusError {
		t.Errorf("Unexpectedly got %v", status)
	}

	if childCalled {
		t.Errorf("Child was unexpectedly called")
	}

	if closeCalled {
		t.Errorf("Close was unexpectedly called")
	}
}
