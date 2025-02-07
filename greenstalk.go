package greenstalk

import (
	"context"
	"errors"
	"fmt"

	"github.com/jbcpollak/greenstalk/core"
	"github.com/jbcpollak/greenstalk/internal"
	"github.com/jbcpollak/greenstalk/util"
)

// BehaviorTree ...
type behaviorTree[Blackboard any] struct {
	ctx        context.Context
	Root       core.Node[Blackboard]
	Blackboard Blackboard
	events     chan core.Event
	visitor    core.Visitor[Blackboard]
}

func NewBehaviorTree[Blackboard any](
	root core.Node[Blackboard],
	bb Blackboard,
	opts ...TreeOption[Blackboard],
) (*behaviorTree[Blackboard], error) {

	var eb internal.ErrorBuilder
	eb.SetMessage("NewBehaviorTree")
	if root == nil {
		eb.Write("Config.Root is nil")
	}

	if eb.Error() != nil {
		return nil, eb.Error()
	}

	tree := &behaviorTree[Blackboard]{
		ctx:        context.TODO(),
		Root:       root,
		Blackboard: bb,
		events:     make(chan core.Event, 100 /* arbitrary */),
		visitor:    func(n core.Walkable[Blackboard]) {},
	}

	// Apply all options to the tree.
	for _, opt := range opts {
		opt(tree)
	}

	return tree, nil
}

// Update propagates an update call down the behavior tree.
func (bt *behaviorTree[Blackboard]) Update(evt core.Event) core.ResultDetails {
	result := core.Update(bt.ctx, bt.Root, bt.Blackboard, evt)

	status := result.Status()
	if status == core.StatusError {
		if details, ok := result.(core.ErrorResultDetails); !ok {
			// Handle if we somehow get an error result that is not an ErrorResultDetails
			return core.ErrorResult(fmt.Errorf("erroneous status encountered %v", details))
		}
	}

	handleRunningResultDetails := func(running core.InitRunningResultDetails) {
		err := running.RunningFn(bt.ctx, func(evt core.Event) error {
			select {
			case <-bt.ctx.Done():
				return bt.ctx.Err()
			case bt.events <- evt:
				return nil
			}
		})
		// If we aren't shutting down, feed the error back through the event loop.
		if err != nil && !errors.Is(err, context.Canceled) {
			internal.Logger.Error("Error in running function", "err", err)

			select {
			case <-bt.ctx.Done():
				return
			case bt.events <- core.ErrorEvent{Err: err}:
				return
			}
		}
	}

	switch status {
	case core.StatusError:
	case core.StatusSuccess:
		// whatever
	case core.StatusFailure:
		// whatever
	case core.StatusRunning:
		if running, ok := result.(core.InitRunningResultDetails); ok {
			go handleRunningResultDetails(running)
		} else if runnings, ok := result.(core.InitRunningResultsDetailsCollection); ok {
			for _, running := range runnings.Results {
				go handleRunningResultDetails(running)
			}
		}
	default:
		return core.ErrorResult(fmt.Errorf("invalid status %v", status))
	}

	bt.visitor(bt.Root)

	return result
}

func (bt *behaviorTree[Blackboard]) EventLoop(evt core.Event) error {
	// Put the first event on the queue.
	bt.events <- evt

	for {
		select {
		case <-bt.ctx.Done():
			return nil
		case evt := <-bt.events:
			if errEvt, ok := evt.(core.ErrorEvent); ok {
				return errEvt.Err
			}
			internal.Logger.Info("Updating with Event", "event", evt)
			result := bt.Update(evt)
			if result.Status() == core.StatusError {
				if details, ok := result.(core.ErrorResultDetails); ok {
					return details.Err
				} else {
					// we should not be able to get here because currently Update ensures that an error status always
					// has ErrorResultDetails, but if that ever changes and we get here, we should still emit an error
					return fmt.Errorf("BT Update returned an error with no details %v", details)
				}
			}
		}
	}
}

// String creates a string representation of the behavior tree
// by traversing it and writing lexical elements to a string.
func (bt *behaviorTree[Blackboard]) String() string {
	return util.NodeToString(bt.Root)
}
