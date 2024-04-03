package action

import (
	"context"
	"testing"
	"time"

	"github.com/jbcpollak/greenstalk"
	"github.com/rs/zerolog/log"

	"github.com/jbcpollak/greenstalk/common/composite"
	"github.com/jbcpollak/greenstalk/core"
	"github.com/jbcpollak/greenstalk/util"
)

func TestSignaller(t *testing.T) {
	// Synchronous, so does not need to be cancelled.
	ctx := context.Background()

	sigChan := make(chan bool, 1)

	params := SignallerParams[bool]{
		BaseParams: "Signaller",

		Channel: sigChan,
		Signal:  true,
	}
	signaller := Signaller[EmptyBlackboard](params)

	var signalSequence = composite.Sequence(
		signaller,
	)

	tree, err := greenstalk.NewBehaviorTree(ctx, signalSequence, EmptyBlackboard{})
	if err != nil {
		panic(err)
	}

	evt := core.DefaultEvent{}
	status := tree.Update(evt)
	util.PrintTreeInColor(tree.Root)

	d := time.Duration(100) * time.Millisecond

	select {
	case c := <-sigChan:
		log.Info().Msgf("got signal %v", c)
		if !c {
			t.Errorf("Expected true, got %v", c)
		}
	case <-time.After(d):
		t.Errorf("Timeout after delaying %v", d)
	}

	if status != core.StatusSuccess {
		t.Errorf("Unexpectedly got %v", status)

	}
}

func TestAsyncSignaller(t *testing.T) {
	// Synchronous, so does not need to be cancelled.
	ctx := context.Background()

	sigChan := make(chan bool)

	params := SignallerParams[bool]{
		BaseParams: "Signaller",

		Channel: sigChan,
		Signal:  true,
	}
	signaller := Signaller[EmptyBlackboard](params)

	var signalSequence = composite.Sequence(
		signaller,
	)

	tree, err := greenstalk.NewBehaviorTree(ctx, signalSequence, EmptyBlackboard{})
	if err != nil {
		panic(err)
	}

	evt := core.DefaultEvent{}
	go tree.EventLoop(evt)
	util.PrintTreeInColor(tree.Root)

	d := time.Duration(100) * time.Millisecond

	select {
	case c := <-sigChan:
		log.Info().Msgf("got count %v", c)
	case <-time.After(d):
		t.Errorf("Timeout after delaying %v", d)
	}

	// if status != core.StatusSuccess {
	// 	t.Errorf("Unexpectedly got %v", status)

	// }
}
