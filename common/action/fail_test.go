package action

import (
	"testing"

	"github.com/jbcpollak/greenstalk"

	"github.com/jbcpollak/greenstalk/common/composite"
	"github.com/jbcpollak/greenstalk/core"
	"github.com/jbcpollak/greenstalk/internal"
	"github.com/jbcpollak/greenstalk/util"
)

func TestFail(t *testing.T) {
	fail := Fail[core.EmptyBlackboard](FailParams{})

	var failSequence = composite.Sequence(
		fail,
	)

	tree, err := greenstalk.NewBehaviorTree(
		failSequence,
		core.EmptyBlackboard{},
		greenstalk.WithVisitors(util.PrintTreeInColor[core.EmptyBlackboard]),
	)
	if err != nil {
		t.Errorf("Unexpectedly got %v", err)
	}

	evt := core.DefaultEvent{}
	result := tree.Update(evt)
	if result.Status() != core.StatusFailure {
		t.Errorf("Unexpectedly got %v", result)

	}

	internal.Logger.Info("Done!")
}
