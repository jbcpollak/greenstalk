package core

import (
	"context"
)

// Update updates a node by calling its Enter method if it is not running,
// then its Tick method, and finally Leave if it is not still running.
func Update(ctx context.Context, node Node, evt Event) ResultDetails {
	var result ResultDetails

	if node.Result().Status() != StatusRunning {
		result = node.Activate(ctx, evt)
	} else {
		result = node.Tick(ctx, evt)
	}

	node.SetResult(result)

	if s := result.Status(); s == StatusError || s == StatusRunning {
		return result
	}

	err := node.Leave(ctx)
	if err != nil {
		result = ErrorResult(err)
		node.SetResult(result)
	}

	return result
}
