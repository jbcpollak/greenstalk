package composite

import (
	"context"
	"math/rand"

	"github.com/jbcpollak/greenstalk/core"
)

// RandomSequence works just like Sequence, except it shuffles
// the order of its children every time it is re-updated.
func RandomSequenceNamed[Blackboard any](name string, children ...core.Node[Blackboard]) core.Node[Blackboard] {
	base := core.NewComposite(core.BaseParams(name), children)
	return &randomSequence[Blackboard]{Composite: base}
}
func RandomSequence[Blackboard any](children ...core.Node[Blackboard]) core.Node[Blackboard] {
	return RandomSequenceNamed("RandomSequence", children...)
}

type randomSequence[Blackboard any] struct {
	core.Composite[Blackboard, core.BaseParams]
}

func (s *randomSequence[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	shuffle(s.Children)

	return s.Tick(ctx, bb, evt)
}

func (s *randomSequence[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	for s.CurrentChild < len(s.Children) {
		result := core.Update(ctx, s.Children[s.CurrentChild], bb, evt)
		if result.Status() != core.StatusSuccess {
			return result
		}
		s.CurrentChild++
	}
	return core.SuccessResult()
}

func (s *randomSequence[Blackboard]) Leave(bb Blackboard) error {
	s.CurrentChild = 0
	return nil
}

func shuffle[Blackboard any](nodes []core.Node[Blackboard]) {
	rand.Shuffle(len(nodes), func(i, j int) {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})
}
