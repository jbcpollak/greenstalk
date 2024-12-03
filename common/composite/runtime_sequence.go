package composite

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
)

func RuntimeSequence[Blackboard any](childrenFn func() ([]core.Node[Blackboard], error)) core.Node[Blackboard] {
	base := core.NewRuntimeComposite(core.BaseParams("RuntimeSequence"), childrenFn)
	return &runtimeSequence[Blackboard]{RuntimeComposite: base}
}

type runtimeSequence[Blackboard any] struct {
	core.RuntimeComposite[Blackboard, core.BaseParams]
}

func (s *runtimeSequence[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	children, err := s.ChildrenFn()

	if err != nil {
		return core.ErrorResult(err)
	}

	s.Children = children
	s.CurrentChild = 0

	// Tick as expected
	return s.Tick(ctx, bb, evt)
}

func (s *runtimeSequence[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	for s.CurrentChild < len(s.Children) {
		result := core.Update(ctx, s.Children[s.CurrentChild], bb, evt)
		if result.Status() != core.StatusSuccess {
			return result
		}
		s.CurrentChild++
	}
	return core.SuccessResult()
}

func (s *runtimeSequence[Blackboard]) Leave(bb Blackboard) error {
	s.Children = []core.Node[Blackboard]{}
	return nil
}
