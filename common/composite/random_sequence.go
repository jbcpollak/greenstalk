package composite

import (
	"context"
	"math/rand"

	"github.com/jbcpollak/greenstalk/v2/core"
)

// RandomSequence works just like Sequence, except it shuffles
// the order of its children every time it is re-updated.
func RandomSequenceNamed(name string, children ...core.Node) core.Node {
	base := core.NewComposite(core.BaseParams(name), children)
	return &randomSequence{Composite: base}
}

func RandomSequence(children ...core.Node) core.Node {
	return RandomSequenceNamed("RandomSequence", children...)
}

type randomSequence struct {
	core.Composite[core.BaseParams]
}

func (s *randomSequence) Activate(ctx context.Context, evt core.Event) core.ResultDetails {
	shuffle(s.Children)

	return s.Tick(ctx, evt)
}

func (s *randomSequence) Tick(ctx context.Context, evt core.Event) core.ResultDetails {
	for s.CurrentChild < len(s.Children) {
		result := core.Update(ctx, s.Children[s.CurrentChild], evt)
		if result.Status() != core.StatusSuccess {
			return result
		}
		s.CurrentChild++
	}
	return core.SuccessResult()
}

func (s *randomSequence) Leave(context.Context) error {
	s.CurrentChild = 0
	return nil
}

func shuffle(nodes []core.Node) {
	rand.Shuffle(len(nodes), func(i, j int) {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})
}

var _ core.Node = (*randomSequence)(nil)
