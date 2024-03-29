package decorator

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
)

// While node repeats the conditions and runs the action if the condition succeeds.
// The action is started after the first success of the condition.
// While succeeds if the action succeeds and fails when any child fails.
// While returns Running if any child is running.
//
// An alternative implementation of the While behavior is:
//
//	UntilFailure {
//	    Sequence {
//	        Condition (custom node)
//	        Action (custom node)
//	    }
//	}
//
// This also allows you to have multiple conditions and multiple actions (just put them after each other in the sequence).
//
// This implementation is taken from https://github.com/DanTulovsky/greenstalk/blob/master/common/decorator/while.go
// See also https://github.com/askft/greenstalk/pull/2
func While[Blackboard any](params core.DecoratorParams, cond, action core.Node[Blackboard]) core.Node[Blackboard] {

	base := core.NewDecorator(core.DecoratorParams{BaseParams: "While"}, cond)
	d := &while[Blackboard]{
		Decorator: base,
		action:    action, // action to run after condition succeeds
	}
	return d
}

type while[Blackboard any] struct {
	core.Decorator[Blackboard]
	action core.Node[Blackboard]
}

func (d *while[Blackboard]) Enter(bb Blackboard) {

}

func (d *while[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {

	// check the condition
	status := core.Update(ctx, d.Child, bb, evt)

	switch status {
	case core.StatusRunning:
		return core.StatusRunning
	case core.StatusFailure:
		return core.StatusFailure
	case core.StatusInvalid:
		return core.StatusInvalid
	}

	// here condition is successful
	actionStatus := core.Update(ctx, d.action, bb, evt)

	return actionStatus

}

func (d *while[Blackboard]) Leave(bb Blackboard) {}
