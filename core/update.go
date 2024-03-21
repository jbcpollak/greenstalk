package core

import "context"

// Update updates a node by calling its Enter method if it is not running,
// then its Tick method, and finally Leave if it is not still running.
func Update[Blackboard any](ctx context.Context, node Node[Blackboard], bb Blackboard, evt Event) NodeResult {

	// var result NodeResult
	if node.Status() != StatusRunning {
		node.Enter(bb)
		// } else {
	}

	result := node.Tick(ctx, bb, evt)
	status := result.Status()
	node.SetStatus(status)

	if status != StatusRunning {
		node.Leave(bb)
	}

	return result
}
