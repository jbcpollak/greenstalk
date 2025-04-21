package composite

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
)

// Selector updates each child in order, returning success as soon as
// a child succeeds. If a child returns Running, the selector node
// will resume execution from that child the next tick.
func Selector[Blackboard any](children ...core.Node[Blackboard]) core.Node[Blackboard] {
	base := core.NewComposite(core.BaseParams("Selector"), children)
	return &selector[Blackboard]{Composite: base}
}

type selector[Blackboard any] struct {
	core.Composite[Blackboard, core.BaseParams]
}

func (s *selector[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	s.CurrentChild = 0

	// Tick as expected
	return s.Tick(ctx, bb, evt)
}

func (s *selector[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	for s.CurrentChild < len(s.Children) {
		result := core.Update(ctx, s.Children[s.CurrentChild], bb, evt)
		if result.Status() != core.StatusFailure {
			return result
		}
		s.CurrentChild++
	}
	return core.FailureResult()
}

func (s *selector[Blackboard]) Leave(bb Blackboard) error {
	return nil
}
