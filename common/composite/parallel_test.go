package composite

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/jbcpollak/greenstalk"
	"github.com/jbcpollak/greenstalk/common/action"

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
		Func: func() core.ResultDetails {
			messageChan <- actionMessage{action1Name, start}
			time.Sleep(200 * time.Millisecond)
			messageChan <- actionMessage{action1Name, end}
			return core.SuccessResult()
		},
	})

	action2 := action.AsyncFunctionAction[core.EmptyBlackboard](action.AsyncFunctionActionParams{
		BaseParams: "Action2",
		Func: func() core.ResultDetails {
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
