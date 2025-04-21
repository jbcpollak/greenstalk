package composite

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
)

// PersistentSequence updates each child in order. If a child
// returns Failure or Running, this node returns the same value,
// and resumes execution from the same child node the next tick.
func PersistentSequence[Blackboard any](children ...core.Node[Blackboard]) core.Node[Blackboard] {
	base := core.NewComposite(core.BaseParams("PersistentSequence"), children)
	return &persistentSequence[Blackboard]{Composite: base}
}

type persistentSequence[Blackboard any] struct {
	core.Composite[Blackboard, core.BaseParams]
}

func (s *persistentSequence[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	return s.Tick(ctx, bb, evt)
}

func (s *persistentSequence[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	for s.CurrentChild < len(s.Children) {
		result := core.Update(ctx, s.Children[s.CurrentChild], bb, evt)
		if result.Status() != core.StatusSuccess {
			return result
		}
		s.CurrentChild++
	}
	return core.SuccessResult()
}

func (s *persistentSequence[Blackboard]) Leave(bb Blackboard) error {
	s.CurrentChild = 0
	return nil
}
