package core

import (
	"context"
)

// Update updates a node by calling its Enter method if it is not running,
// then its Tick method, and finally Leave if it is not still running.
func Update[Blackboard any](ctx context.Context, node Node[Blackboard], bb Blackboard, evt Event) ResultDetails {

	var result ResultDetails

	if node.Result().Status() != StatusRunning {
		result = node.Activate(ctx, bb, evt)
	} else {
		result = node.Tick(ctx, bb, evt)
	}

	node.SetResult(result)

	if result.Status() == StatusError ||
		result.Status() == StatusRunning {
		return result
	}

	err := node.Leave(bb)
	if err != nil {
		result = ErrorResult(err)
		node.SetResult(result)
	}

	return result
}
