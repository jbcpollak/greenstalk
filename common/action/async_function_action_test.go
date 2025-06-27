package action

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/jbcpollak/greenstalk"
	"github.com/jbcpollak/greenstalk/core"
)

func TestAsyncFunctionAction(t *testing.T) {

	asyncFunctionExpectedResultsMap := map[core.Status]func(ctx context.Context) core.ResultDetails{}

	asyncFunctionExpectedResultsMap[core.StatusSuccess] = func(ctx context.Context) core.ResultDetails {
		return core.SuccessResult()
	}

	asyncFunctionExpectedResultsMap[core.StatusFailure] = func(ctx context.Context) core.ResultDetails {
		return core.FailureResult()
	}

	const error_msg = "error here"
	asyncFunctionExpectedResultsMap[core.StatusError] = func(ctx context.Context) core.ResultDetails {
		return core.ErrorResult(errors.New(error_msg))
	}

	for status, fn := range asyncFunctionExpectedResultsMap {
		testAsyncFunctionAction(t, status, fn, error_msg)
	}

	// invalid usage - should never return RunningResult from the async function
	testAsyncFunctionAction(t, core.StatusError, func(ctx context.Context) core.ResultDetails {
		return core.RunningResult()
	}, "async function returned invalid status of StatusRunning")
}

func testAsyncFunctionAction(t *testing.T, expectedStatus core.Status, fn func(ctx context.Context) core.ResultDetails, errMsg string) {
	asyncFunctionAction := AsyncFunctionAction[core.EmptyBlackboard](AsyncFunctionActionParams{
		BaseParams: "asyncFunctionNode",
		Func:       fn,
	})

	ctx, cancel := context.WithCancel(context.Background())

	nodeWG := sync.WaitGroup{}
	nodeWG.Add(1)
	asyncNodeStatuses := []core.Status{}
	visitor := func(node core.Walkable[core.EmptyBlackboard]) {
		if node.Id() == asyncFunctionAction.Id() {
			status := node.Result().Status()
			asyncNodeStatuses = append(asyncNodeStatuses, status)
			if status == expectedStatus {
				nodeWG.Done()
			}
		}
	}

	tree, err := greenstalk.NewBehaviorTree(
		asyncFunctionAction,
		core.EmptyBlackboard{},
		greenstalk.WithContext[core.EmptyBlackboard](ctx),
		greenstalk.WithVisitors(visitor),
	)

	if err != nil {
		cancel()
		t.Errorf("Unexpectedly got %v", err)
	}

	evt := core.DefaultEvent{}

	treeWG := sync.WaitGroup{}
	treeWG.Add(1)
	go func() {
		err := tree.EventLoop(evt)
		if err != nil && expectedStatus != core.StatusError && err.Error() != errMsg {
			t.Errorf("Unexpectedly got %v", err)
		}
		treeWG.Done()
	}()

	// first, wait for the function to finish
	nodeWG.Wait()
	// then cancel the tree context
	cancel()
	// and finally wait for the event loop to exit
	treeWG.Wait()

	if !reflect.DeepEqual(asyncNodeStatuses, []core.Status{core.StatusRunning, expectedStatus}) {
		t.Errorf("Expected %v got %v", []core.Status{core.StatusRunning, expectedStatus}, asyncNodeStatuses)
	}

	if tree.Root.Result().Status() != expectedStatus {
		t.Errorf("Expected %v got %v", expectedStatus, tree.Root.Result().Status())
	}

}
