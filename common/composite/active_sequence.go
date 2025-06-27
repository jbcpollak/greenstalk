package composite

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
)

// ActiveSequence ticks each child in order. Returns success if
// all children succeed in one tick, else returns the status of
// the non-succeeding node. Restarts iteration the next tick.
func ActiveSequence[Blackboard any](children ...core.Node[Blackboard]) core.Node[Blackboard] {
	return ActiveSequenceNamed("ActiveSequence", children...)
}
func ActiveSequenceNamed[Blackboard any](name string, children ...core.Node[Blackboard]) core.Node[Blackboard] {
	base := core.NewComposite(core.BaseParams(name), children)
	return &activeSequence[Blackboard]{Composite: base}
}

type activeSequence[Blackboard any] struct {
	core.Composite[Blackboard, core.BaseParams]
}

func (s *activeSequence[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	// No distinction between activation and ticking
	return s.Tick(ctx, bb, evt)
}

func (s *activeSequence[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	for i := 0; i < len(s.Children); i++ {
		result := core.Update(ctx, s.Children[i], bb, evt)
		if result.Status() != core.StatusSuccess {
			return result
		}
	}
	return core.SuccessResult()
}

func (s *activeSequence[Blackboard]) Leave(bb Blackboard) error {
	return nil
}
