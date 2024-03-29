package greenstalk

import (
	"context"
	"testing"
	"time"

	"github.com/jbcpollak/greenstalk/core"
	"github.com/jbcpollak/greenstalk/util"
	"github.com/rs/zerolog/log"

	// Use dot imports to make a tree definition look nice.
	// Be careful when doing this! These packages export
	// common word identifiers such as "Fail" and "Sequence".
	. "github.com/jbcpollak/greenstalk/common/action"
	. "github.com/jbcpollak/greenstalk/common/composite"
	. "github.com/jbcpollak/greenstalk/common/decorator"
)

type TestBlackboard struct {
	id    int
	count uint
}

// Counter increments a counter on the blackboard.
func Counter(params core.BaseParams) core.Node[TestBlackboard] {
	name := params.Name()

	base := core.NewLeaf[TestBlackboard](core.BaseParams("Counter "+name), struct{}{})
	return &counter{Leaf: base}
}

// succeed ...
type counter struct {
	core.Leaf[TestBlackboard, core.BaseParams, struct{}]
}

// Enter ...
func (a *counter) Enter(bb TestBlackboard) {}

// Tick ...
var countChan = make(chan uint, 0)

func (a *counter) Tick(ctx context.Context, bb TestBlackboard, evt core.Event) core.NodeResult {
	log.Info().Msgf("%s: Incrementing count", a.Name())
	bb.count++
	countChan <- bb.count
	return core.StatusSuccess
}

// Leave ...
func (a *counter) Leave(bb TestBlackboard) {}

var synchronousRoot = Sequence[TestBlackboard](
	Repeater(RepeaterParams{N: 2}, Fail[TestBlackboard](FailParams{}, struct{}{})),
	Succeed[TestBlackboard](SucceedParams{}, struct{}{}),
)

func TestUpdate(t *testing.T) {
	log.Info().Msg("Testing synchronous tree...")

	// Synchronous, so does not need to be cancelled.
	ctx := context.Background()

	tree, err := NewBehaviorTree(ctx, synchronousRoot, TestBlackboard{id: 42})
	if err != nil {
		panic(err)
	}

	for {
		evt := core.DefaultEvent{}
		status := tree.Update(evt)
		util.PrintTreeInColor(tree.Root)
		if status == core.StatusSuccess {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	log.Info().Msg("Done!")
}

var delay = 100
var asynchronousRoot = Sequence[TestBlackboard](
	// Repeater(core.Params{"n": 2}, Fail[TestBlackboard](nil, nil)),
	AsyncDelayer[TestBlackboard](
		AsyncDelayerParams{
			DecoratorParams: core.DecoratorParams{
				BaseParams: core.BaseParams("First"),
			},
			Delay: time.Duration(delay) * time.Millisecond,
		},
		Counter("First"),
	),
	AsyncDelayer[TestBlackboard](
		AsyncDelayerParams{
			DecoratorParams: core.DecoratorParams{
				BaseParams: core.BaseParams("Second"),
			},
			Delay: time.Duration(delay) * time.Millisecond,
		},
		Counter(core.BaseParams("Second")),
	),
)

func getCount(d time.Duration) (uint, bool) {
	select {
	case c := <-countChan:
		log.Info().Msgf("got count %v", c)
		return c, true
	case <-time.After(d):
		log.Info().Msgf("Timeout after delaying %v", d)
		return 0, false
	}
}

func TestEventLoop(t *testing.T) {
	log.Info().Msg("Testing asynchronous tree...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bb := TestBlackboard{id: 42, count: 0}

	tree, err := NewBehaviorTree(ctx, asynchronousRoot, bb)
	if err != nil {
		panic(err)
	}

	evt := core.DefaultEvent{}
	go tree.EventLoop(evt)

	// Wait half the delay and verify no value sent
	first_halfway, ok := getCount(time.Duration(delay/2) * time.Millisecond)
	if ok {
		t.Errorf("Unexpectedly got count %d", first_halfway)
	} else {
		log.Info().Msg("Halfway through first delay counter correctly is 0")
	}

	// Sleep a bit more
	first_after, ok := getCount(time.Duration(delay/2+10) * time.Millisecond)
	if ok && first_after != 1 {
		t.Errorf("Expected count to be 1, got %d", first_after)
	} else {
		log.Info().Msg("After first delay, counter is 1")
	}

	// Wait half the delay and verify value is 0
	second_halfway, ok := getCount(time.Duration(delay/2) * time.Millisecond)
	if ok {
		t.Errorf("Unexpectedly got count %d", second_halfway)
	} else {
		log.Info().Msg("Halfway through second delay counter correctly is 1")
	}

	// Shut it _down_
	log.Info().Msg("Shutting down...")
	cancel()

	after_cancel, ok := getCount(time.Duration(delay/2) * time.Millisecond)

	// Ensure we shut down before the second tick
	if ok {
		t.Errorf("Expected to shut down before second tick but got %d", after_cancel)
	}

	log.Info().Msg("Done!")
}
