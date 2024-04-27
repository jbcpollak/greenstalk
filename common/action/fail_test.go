package action

import (
	"testing"

	"github.com/jbcpollak/greenstalk"
	"github.com/rs/zerolog/log"

	"github.com/jbcpollak/greenstalk/common/composite"
	"github.com/jbcpollak/greenstalk/core"
	"github.com/jbcpollak/greenstalk/util"
)

func TestFail(t *testing.T) {
	fail := Fail[core.EmptyBlackboard](FailParams{})

	var failSequence = composite.Sequence[core.EmptyBlackboard](
		fail,
	)

	tree, err := greenstalk.NewBehaviorTree(
		failSequence,
		core.EmptyBlackboard{},
		greenstalk.WithVisitor(util.PrintTreeInColor[core.EmptyBlackboard]),
	)
	if err != nil {
		panic(err)
	}

	evt := core.DefaultEvent{}
	status := tree.Update(evt)
	if status != core.StatusFailure {
		t.Errorf("Unexpectedly got %v", status)

	}

	log.Info().Msg("Done!")
}
