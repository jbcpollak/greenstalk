package action

import (
	"context"
	"testing"

	"github.com/jbcpollak/greenstalk"
	"github.com/rs/zerolog/log"

	"github.com/jbcpollak/greenstalk/common/composite"
	"github.com/jbcpollak/greenstalk/core"
	"github.com/jbcpollak/greenstalk/util"
)

type EmptyBlackboard struct{}

func TestFail(t *testing.T) {
	// Synchronous, so does not need to be cancelled.
	ctx := context.Background()

	fail := Fail[EmptyBlackboard](FailParams{})

	var failSequence = composite.Sequence[EmptyBlackboard](
		fail,
	)

	tree, err := greenstalk.NewBehaviorTree(ctx, failSequence, EmptyBlackboard{})
	if err != nil {
		panic(err)
	}

	evt := core.DefaultEvent{}
	status := tree.Update(evt)
	util.PrintTreeInColor(tree.Root)
	if status != core.StatusFailure {
		t.Errorf("Unexpectedly got %v", status)

	}

	log.Info().Msg("Done!")
}
