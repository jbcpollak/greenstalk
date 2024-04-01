package composite

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
)

// Sequence updates each child in order, returning success only if
// all children succeed. If a child returns Running, the sequence node
// will resume execution from that child the next tick.
func Sequence[Blackboard any](children ...core.Node[Blackboard]) core.Node[Blackboard] {
	base := core.NewComposite(core.BaseParams("Sequence"), children)
	return &sequence[Blackboard]{Composite: base}
}

type sequence[Blackboard any] struct {
	core.Composite[Blackboard]
}

func (s *sequence[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	s.Composite.CurrentChild = 0

	// Tick as expected
	return s.Tick(ctx, bb, evt)
}

func (s *sequence[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	for s.CurrentChild < len(s.Children) {
		status := core.Update(ctx, s.Children[s.CurrentChild], bb, evt)
		if status != core.StatusSuccess {
			return status
		}
		s.CurrentChild++
	}
	return core.StatusSuccess
}

func (s *sequence[Blackboard]) Leave(bb Blackboard) {}
