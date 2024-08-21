package composite

import (
	"context"
	"fmt"

	"github.com/jbcpollak/greenstalk/common/state"
	"github.com/jbcpollak/greenstalk/core"
)

// IndexSequence functions the same way as Sequence, but it starts at a child given by the index state. Returns an error result
// if the provided index is out of bounds.
func IndexSequence[Blackboard any](index state.StateGetter[int], children ...core.Node[Blackboard]) core.Node[Blackboard] {
	base := core.NewComposite(core.BaseParams("IndexSequence"), children)
	return &indexSequence[Blackboard]{Composite: base, index: index}
}

type indexSequence[Blackboard any] struct {
	core.Composite[Blackboard, core.BaseParams]
	index state.StateGetter[int]
}

func (s *indexSequence[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	firstChildIndex := s.index.Get()
	if firstChildIndex < 0 || firstChildIndex >= len(s.Children) {
		return core.ErrorResult(fmt.Errorf("index %v out of bounds: children %v", firstChildIndex, len(s.Children)))
	}

	s.Composite.CurrentChild = firstChildIndex

	// Tick as expected
	return s.Tick(ctx, bb, evt)
}

func (s *indexSequence[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	for s.CurrentChild < len(s.Children) {
		result := core.Update(ctx, s.Children[s.CurrentChild], bb, evt)
		if result.Status() != core.StatusSuccess {
			return result
		}
		s.CurrentChild++
	}
	return core.SuccessResult()
}

func (s *indexSequence[Blackboard]) Leave(bb Blackboard) error {
	return nil
}
