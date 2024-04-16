package greenstalk

import (
	"context"
	"fmt"

	"github.com/jbcpollak/greenstalk/core"
	"github.com/jbcpollak/greenstalk/internal"
	"github.com/jbcpollak/greenstalk/util"
	"github.com/rs/zerolog/log"
)

// BehaviorTree ...
type BehaviorTree[Blackboard any] struct {
	ctx        context.Context
	Root       core.Node[Blackboard]
	Blackboard Blackboard
	events     chan core.Event
}

func NewBehaviorTree[Blackboard any](ctx context.Context, root core.Node[Blackboard], bb Blackboard) (*BehaviorTree[Blackboard], error) {
	var eb internal.ErrorBuilder
	eb.SetMessage("NewBehaviorTree")
	if root == nil {
		eb.Write("Config.Root is nil")
	}

	if eb.Error() != nil {
		return nil, eb.Error()
	}
	tree := &BehaviorTree[Blackboard]{
		ctx:        ctx,
		Root:       root,
		Blackboard: bb,
		events:     make(chan core.Event, 100 /* arbitrary */),
	}
	return tree, nil
}

// Update propagates an update call down the behavior tree.
func (bt *BehaviorTree[Blackboard]) Update(evt core.Event) core.Status {
	result := core.Update(bt.ctx, bt.Root, bt.Blackboard, evt)

	status := result.Status()
	if status == core.StatusError {
		if details, ok := result.(core.ErrorResultDetails); ok {
			panic(details.Err)
		} else {
			// Handle if we somehow get an error result that is not an ErrorResultDetails
			panic(fmt.Errorf("erroneous status encountered %v", details))
		}
	}

	switch status {
	case core.StatusSuccess:
		// whatever
	case core.StatusFailure:
		// whatever
	case core.StatusRunning:
		if running, ok := result.(core.InitRunningResultDetails); ok {
			go running.RunningFn(bt.ctx, func(evt core.Event) error {
				bt.events <- evt
				return nil
			})
		}
	default:
		panic(fmt.Errorf("invalid status %v", status))
	}

	return status
}

func (bt *BehaviorTree[Blackboard]) EventLoop(evt core.Event) {
	defer close(bt.events)

	// Put the first event on the queue.
	bt.events <- evt

	for {
		select {
		case <-bt.ctx.Done():
			return
		case evt := <-bt.events:
			log.Info().Msgf("Event: %v", evt)
			bt.Update(evt)

			// TODO: Change to visitor pattern.
			util.PrintTreeInColor(bt.Root)
		}
	}
}

// String creates a string representation of the behavior tree
// by traversing it and writing lexical elements to a string.
func (bt *BehaviorTree[Blackboard]) String() string {
	return util.NodeToString(bt.Root)
}
