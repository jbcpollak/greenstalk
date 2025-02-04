package composite

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/jbcpollak/greenstalk"
	"github.com/jbcpollak/greenstalk/common/action"
	"github.com/jbcpollak/greenstalk/common/decorator"
	"github.com/jbcpollak/greenstalk/common/state"

	"github.com/jbcpollak/greenstalk/core"
)

type actionMessage struct {
	name    string
	message string
}

// verifies that two actions are executed in parallel
func TestParallelExecution(t *testing.T) {
	messageChan := make(chan actionMessage, 4)
	const action1Name = "action1"
	const action2Name = "action2"
	const start = "start"
	const end = "end"

	action1 := action.AsyncFunctionAction[core.EmptyBlackboard](action.AsyncFunctionActionParams{
		BaseParams: "Action1",
		Func: func(ctx context.Context) core.ResultDetails {
			messageChan <- actionMessage{action1Name, start}
			time.Sleep(200 * time.Millisecond)
			messageChan <- actionMessage{action1Name, end}
			return core.SuccessResult()
		},
	})

	action2 := action.AsyncFunctionAction[core.EmptyBlackboard](action.AsyncFunctionActionParams{
		BaseParams: "Action2",
		Func: func(ctx context.Context) core.ResultDetails {
			messageChan <- actionMessage{action2Name, start}
			time.Sleep(200 * time.Millisecond)
			messageChan <- actionMessage{action2Name, end}
			return core.SuccessResult()
		},
	})

	parallel := Parallel(2, 1, action1, action2)

	sigChan := make(chan bool)
	params := action.SignallerParams[bool]{
		BaseParams: "Signaller",

		Channel: sigChan,
		Signal:  true,
	}
	signaller := action.Signaller[core.EmptyBlackboard](params)

	ctx, cancel := context.WithCancel(context.Background())

	tree, err := greenstalk.NewBehaviorTree(
		Sequence(parallel, signaller),
		core.EmptyBlackboard{},
		greenstalk.WithContext[core.EmptyBlackboard](ctx),
	)
	if err != nil {
		t.Errorf("Unexpectedly got %v", err)
	}

	evt := core.DefaultEvent{}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err = tree.EventLoop(evt)
		if err != nil {
			t.Errorf("Unexpectedly got %v", err)
		}
		wg.Done()
	}()

	signal := <-sigChan
	if !signal {
		t.Errorf("Unexpectedly got signal %v", signal)
	}

	cancel()

	wg.Wait()
	close(messageChan)

	messages := []actionMessage{}
	for message := range messageChan {
		messages = append(messages, message)
	}

	if len(messages) != 4 {
		t.Errorf("Expected 4 messages, got %d", len(messages))
	}

	if messages[0].message != start || messages[1].message != start || messages[2].message != end || messages[3].message != end {
		t.Errorf("Unexpected sequence of starts and ends")
	}

	if !((messages[0].name == action1Name && messages[1].name == action2Name) || (messages[0].name == action2Name && messages[1].name == action1Name)) {
		t.Errorf("Unexpected order of actions")
	}

	if !((messages[2].name == action1Name && messages[3].name == action2Name) || (messages[2].name == action2Name && messages[3].name == action1Name)) {
		t.Errorf("Unexpected order of actions")
	}
}

func TestParallelCompletionReset(t *testing.T) {
	// we need to activate the same Parallel node 2+ times - before the fix
	// the second time around the node would hang without any child nodes being called

	counterParam := state.StateProvider[int]{}
	counterParam.Set(0)

	treeNode := decorator.RepeatUntil(decorator.RepeatUntilParams{
		BaseParams: "loop",
		Until: func(status core.ResultDetails) bool {
			if status.Status() != core.StatusSuccess {
				panic("Unexpected status")
			}

			count := counterParam.Get()
			return count > 1
		}},
		Sequence(
			Parallel(2, 1,
				action.Succeed[core.EmptyBlackboard](action.SucceedParams{BaseParams: "success1"}),
				action.Succeed[core.EmptyBlackboard](action.SucceedParams{BaseParams: "success2"}),
			),
			action.FunctionAction[core.EmptyBlackboard](action.FunctionActionParams{
				BaseParams: "increment",
				Func: func(ctx context.Context) core.ResultDetails {
					count := counterParam.Get()
					count++
					counterParam.Set(count)
					return core.SuccessResult()
				},
			}),
		),
	)

	sigChan := make(chan bool)
	params := action.SignallerParams[bool]{
		BaseParams: "Signaller",

		Channel: sigChan,
		Signal:  true,
	}
	signaller := action.Signaller[core.EmptyBlackboard](params)

	ctx, cancel := context.WithCancel(context.Background())

	tree, err := greenstalk.NewBehaviorTree(
		Sequence(treeNode, signaller),
		core.EmptyBlackboard{},
		greenstalk.WithContext[core.EmptyBlackboard](ctx),
	)
	if err != nil {
		t.Errorf("Unexpectedly got %v", err)
	}

	evt := core.DefaultEvent{}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err = tree.EventLoop(evt)
		if err != nil {
			t.Errorf("Unexpectedly got %v", err)
		}
		wg.Done()
	}()

	signal := <-sigChan
	if !signal {
		t.Errorf("Unexpectedly got signal %v", signal)
	}

	cancel()

	wg.Wait()

	// before the fix, the behavior tree would be stuck in running, and RepeatUntil verifies that the parallel node return success
	// so if the behavior tree exits cleanly, the test is comprehensive enough

}

func TestNestedParallels(t *testing.T) {
	// When a parallel node returns running functions, it wraps them in a collection
	// This tests that a top level parallel node can handle the running function collections
	// returned by child parallel nodes

	synchronizedCounter := state.SynchronizedStateProvider[int]{}

	makeAsyncIncrement := func() core.Node[core.EmptyBlackboard] {
		return action.AsyncFunctionAction[core.EmptyBlackboard](action.AsyncFunctionActionParams{
			BaseParams: "increment",
			Func: func(ctx context.Context) core.ResultDetails {
				synchronizedCounter.Lock()
				defer synchronizedCounter.Unlock()
				synchronizedCounter.Set(synchronizedCounter.Get() + 1)
				return core.SuccessResult()
			},
		})
	}

	treeNode := Parallel(2, 1,
		Parallel(2, 1, makeAsyncIncrement(), makeAsyncIncrement()),
		Parallel(2, 1, makeAsyncIncrement(), makeAsyncIncrement()),
	)

	sigChan := make(chan bool)
	params := action.SignallerParams[bool]{
		BaseParams: "Signaller",

		Channel: sigChan,
		Signal:  true,
	}
	signaller := action.Signaller[core.EmptyBlackboard](params)

	ctx, cancel := context.WithCancel(context.Background())

	tree, err := greenstalk.NewBehaviorTree(
		Sequence(treeNode, signaller),
		core.EmptyBlackboard{},
		greenstalk.WithContext[core.EmptyBlackboard](ctx),
	)
	if err != nil {
		t.Errorf("Unexpectedly got %v", err)
	}

	evt := core.DefaultEvent{}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err = tree.EventLoop(evt)
		if err != nil {
			t.Errorf("Unexpectedly got %v", err)
		}
		wg.Done()
	}()

	signal := <-sigChan
	if !signal {
		t.Errorf("Unexpectedly got signal %v", signal)
	}

	cancel()

	wg.Wait()

	if synchronizedCounter.Get() != 4 {
		t.Errorf("Expected 4, got %d", synchronizedCounter.Get())
	}
}
