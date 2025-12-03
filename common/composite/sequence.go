package composite

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
)

// Sequence updates each child in order, returning success only if
// all children succeed. If a child returns Running, the sequence node
// will resume execution from that child the next tick.
func SequenceNamed(name string, children ...core.Node) core.Node {
	base := core.NewComposite(core.BaseParams(name), children)
	return &sequence{Composite: base}
}

func Sequence(children ...core.Node) core.Node {
	return SequenceNamed("Sequence", children...)
}

type sequence struct {
	core.Composite[core.BaseParams]
}

func (s *sequence) Activate(ctx context.Context, evt core.Event) core.ResultDetails {
	s.CurrentChild = 0

	// Tick as expected
	return s.Tick(ctx, evt)
}

func (s *sequence) Tick(ctx context.Context, evt core.Event) core.ResultDetails {
	for s.CurrentChild < len(s.Children) {
		result := core.Update(ctx, s.Children[s.CurrentChild], evt)
		if result.Status() != core.StatusSuccess {
			return result
		}
		s.CurrentChild++
	}
	return core.SuccessResult()
}

func (s *sequence) Leave(context.Context) error {
	return nil
}

var _ core.Node = (*sequence)(nil)
