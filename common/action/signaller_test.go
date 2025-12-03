package action

import (
	"testing"
	"time"

	"github.com/jbcpollak/greenstalk/v2"
	"github.com/jbcpollak/greenstalk/v2/common/composite"
	"github.com/jbcpollak/greenstalk/v2/core"
	"github.com/jbcpollak/greenstalk/v2/internal"
	"github.com/jbcpollak/greenstalk/v2/util"
)

func TestSignaller(t *testing.T) {
	sigChan := make(chan bool, 1)

	params := SignallerParams[bool]{
		BaseParams: "Signaller",
		Channel:    sigChan,
		Signal:     true,
	}
	signaller := Signaller(params)

	signalSequence := composite.Sequence(
		signaller,
	)

	tree, err := greenstalk.NewBehaviorTree(
		signalSequence,
		greenstalk.WithVisitors(util.PrintTreeInColor),
	)
	if err != nil {
		t.Errorf("Unexpectedly got %v", err)
	}

	evt := core.DefaultEvent{}
	result := tree.Update(evt)

	d := time.Duration(100) * time.Millisecond

	select {
	case c := <-sigChan:
		internal.Logger.Info("got signal", "signal", c)
		if !c {
			t.Errorf("Expected true, got %v", c)
		}
	case <-time.After(d):
		t.Errorf("Timeout after delaying %v", d)
	}

	if result.Status() != core.StatusSuccess {
		t.Errorf("Unexpectedly got %v", result)
	}
}

func TestAsyncSignaller(t *testing.T) {
	sigChan := make(chan bool)

	params := SignallerParams[bool]{
		BaseParams: "Signaller",
		Channel:    sigChan,
		Signal:     true,
	}
	signaller := Signaller(params)

	signalSequence := composite.Sequence(
		signaller,
	)

	tree, err := greenstalk.NewBehaviorTree(
		signalSequence,
		greenstalk.WithVisitors(util.PrintTreeInColor),
	)
	if err != nil {
		t.Errorf("Unexpectedly got %v", err)
	}

	evt := core.DefaultEvent{}
	go func() {
		err := tree.EventLoop(evt)
		if err != nil {
			t.Errorf("Unexpectedly got %v", err)
		}
	}()

	d := time.Duration(100) * time.Millisecond

	select {
	case c := <-sigChan:
		internal.Logger.Info("got signal", "signal", c)
	case <-time.After(d):
		t.Errorf("Timeout after delaying %v", d)
	}

	// if status != core.StatusSuccess {
	// 	t.Errorf("Unexpectedly got %v", status)

	// }
}
