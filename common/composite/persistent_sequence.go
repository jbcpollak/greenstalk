package composite

import (
	"context"

	"github.com/jbcpollak/greenstalk/v2/core"
)

// PersistentSequence updates each child in order. If a child
// returns Failure or Running, this node returns the same value,
// and resumes execution from the same child node the next tick.
func PersistentSequenceNamed(name string, children ...core.Node) core.Node {
	base := core.NewComposite(core.BaseParams(name), children)
	return &persistentSequence{Composite: base}
}

func PersistentSequence(children ...core.Node) core.Node {
	return PersistentSequenceNamed("PersistentSequence", children...)
}

type persistentSequence struct {
	core.Composite[core.BaseParams]
}

func (s *persistentSequence) Activate(ctx context.Context, evt core.Event) core.ResultDetails {
	return s.Tick(ctx, evt)
}

func (s *persistentSequence) Tick(ctx context.Context, evt core.Event) core.ResultDetails {
	for s.CurrentChild < len(s.Children) {
		result := core.Update(ctx, s.Children[s.CurrentChild], evt)
		if result.Status() != core.StatusSuccess {
			return result
		}
		s.CurrentChild++
	}
	return core.SuccessResult()
}

func (s *persistentSequence) Leave(context.Context) error {
	s.CurrentChild = 0
	return nil
}

var _ core.Node = (*persistentSequence)(nil)
