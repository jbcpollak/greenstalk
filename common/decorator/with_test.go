package decorator

import (
	"context"
	"fmt"
	"io"
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

type testCloser struct {
	closeCalled *bool
}

func (t testCloser) Close() error {
	*t.closeCalled = true
	return nil
}

func TestWith(t *testing.T) {
	var wg sync.WaitGroup

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

	tree, err := greenstalk.NewBehaviorTree(
		testSequence,
		core.EmptyBlackboard{},
		greenstalk.WithContext[core.EmptyBlackboard](ctx),
		greenstalk.WithVisitor(util.PrintTreeInColor[core.EmptyBlackboard]),
	)
	if err != nil {
		panic(err)
	}

	evt := core.DefaultEvent{}
	wg.Add(1)
	go func() {
		tree.EventLoop(evt)
		wg.Done()
	}()

	d := time.Duration(100) * time.Millisecond

	internal.WaitForSignalOrTimeout(sigChan, d)

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


type errorCloser struct {
	closeCalled *bool
}

func (t errorCloser) Close() error {
	*t.closeCalled = true
	return fmt.Errorf("This is an explected error")
}

func TestWithCloserError(t *testing.T) {
	var wg sync.WaitGroup

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
	closer := errorCloser{
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

	tree, err := greenstalk.NewBehaviorTree(
		testSequence,
		core.EmptyBlackboard{},
		greenstalk.WithContext[core.EmptyBlackboard](ctx),
		greenstalk.WithVisitor(util.PrintTreeInColor[core.EmptyBlackboard]),
	)
	if err != nil {
		panic(err)
	}

	evt := core.DefaultEvent{}
	wg.Add(1)
	go func() {
		tree.EventLoop(evt)
		wg.Done()
	}()

	d := time.Duration(100) * time.Millisecond

	internal.WaitForSignalOrTimeout(sigChan, d)

	cancel()
	wg.Wait()
	status := tree.Root.Result().Status()
	if status != core.StatusError {
		t.Errorf("Unexpectedly got %v", status)
	}

	if !*childCalled {
		t.Errorf("Child was not called")
	}

	if !*closeCalled {
		t.Errorf("Close was not called")
	}
}

func TestWithInitError(t *testing.T) {
	var wg sync.WaitGroup

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

	with := With(func() (io.Closer, error) {
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
		greenstalk.WithVisitor(util.PrintTreeInColor[core.EmptyBlackboard]),
	)
	if err != nil {
		panic(err)
	}

	evt := core.DefaultEvent{}
	wg.Add(1)
	go func() {
		err = tree.EventLoop(evt)
		wg.Done()
	}()

	d := time.Duration(100) * time.Millisecond

	internal.WaitForSignalOrTimeout(sigChan, d)

	cancel()
	wg.Wait()
	status := tree.Root.Result().Status()
	if status != core.StatusError {
		t.Errorf("Unexpectedly got %v", status)
	}

	if err.Error() != "This is an error" {
		t.Errorf("Error does not have correct contents: %v", err)
	}

	if *childCalled {
		t.Errorf("Child was unexpectedly called")
	}

	if *closeCalled {
		t.Errorf("Close was unexpectedly called")
	}
}
