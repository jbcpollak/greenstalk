package action

import (
	"testing"

	"github.com/jbcpollak/greenstalk/v2"
	"github.com/jbcpollak/greenstalk/v2/common/composite"
	"github.com/jbcpollak/greenstalk/v2/core"
	"github.com/jbcpollak/greenstalk/v2/internal"
	"github.com/jbcpollak/greenstalk/v2/util"
)

func TestFail(t *testing.T) {
	fail := Fail(FailParams{})

	failSequence := composite.Sequence(
		fail,
	)

	tree, err := greenstalk.NewBehaviorTree(
		failSequence,
		greenstalk.WithVisitors(util.PrintTreeInColor),
	)
	if err != nil {
		t.Errorf("Unexpectedly got %v", err)
	}

	evt := core.DefaultEvent{}
	result := tree.Update(t.Context(), evt)
	if result.Status() != core.StatusFailure {
		t.Errorf("Unexpectedly got %v", result)
	}

	internal.Logger.Info("Done!")
}
