package condition

import (
	"context"
	"testing"

	"github.com/jbcpollak/greenstalk"
	"github.com/jbcpollak/greenstalk/common/action"
	"github.com/jbcpollak/greenstalk/common/decorator"
	"github.com/jbcpollak/greenstalk/core"
	"github.com/jbcpollak/greenstalk/util"
)

type BlackboardWithCounter struct {
	counter int
}

func TestFunctionCondition(t *testing.T) {
	ctx := context.Background()

	var condition = FunctionCondition(FunctionConditionParams[BlackboardWithCounter]{
		Func: func(bb BlackboardWithCounter) bool {
			return bb.counter < 2
		},
	}, FunctionConditionReturns{})

	var action = action.SideEffect(action.SideEffectParams[BlackboardWithCounter]{
		Func: func(bb BlackboardWithCounter) {
			bb.counter++
		},
	}, action.SideEffectReturns{})

	var root = decorator.While(core.DecoratorParams{}, condition, action)
	var blackboard = BlackboardWithCounter{counter: 0}

	tree, err := greenstalk.NewBehaviorTree(ctx, root, blackboard)
	if err != nil {
		panic(err)
	}

	evt := core.DefaultEvent{}
	status := tree.Update(evt)
	util.PrintTreeInColor(tree.Root)
	if status != core.StatusSuccess {
		t.Errorf("Unexpectedly got %v", status)
	}

	if blackboard.counter != 2 {
		t.Errorf("Unexpected counter value %v", blackboard.counter)
	}
}
