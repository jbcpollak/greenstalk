package condition

import (
	"context"
	"testing"

	"github.com/jbcpollak/greenstalk"
	"github.com/jbcpollak/greenstalk/common/action"
	"github.com/jbcpollak/greenstalk/common/composite"
	"github.com/jbcpollak/greenstalk/common/decorator"
	"github.com/jbcpollak/greenstalk/core"
	"github.com/jbcpollak/greenstalk/util"
)

type EmptyBlackboard struct{}

func TestFunctionCondition(t *testing.T) {
	ctx := context.Background()

	var counter = 0
	var counterPtr = &counter
	var condition = FunctionCondition[EmptyBlackboard](FunctionConditionParams{
		Func: func() bool {
			return *counterPtr < 2
		},
	}, FunctionConditionReturns{})

	var action = action.FunctionAction[EmptyBlackboard](action.FunctionActionParams{
		Func: func() {
			*counterPtr++
		},
	}, action.FunctionActionReturns{})

	var root = decorator.UntilFailure(
		core.DecoratorParams{},
		composite.Sequence(condition, action),
	)

	var blackboard = EmptyBlackboard{}

	tree, err := greenstalk.NewBehaviorTree(ctx, root, blackboard)
	if err != nil {
		panic(err)
	}

	evt := core.DefaultEvent{}
	var status core.Status
	loopCounter := 0
	for ; loopCounter < 10; loopCounter++ {
		status = tree.Update(evt)
		if status != core.StatusRunning {
			if status != core.StatusSuccess {
				t.Errorf("Unexpectedly got %v", status)
			}
			break
		}
	}

	util.PrintTreeInColor(tree.Root)
	if counter != 2 {
		t.Errorf("Unexpected counter value %v", counter)
	}

	if loopCounter != 2 {
		t.Errorf("Unexpected loopCounter value %v", loopCounter)
	}
}
