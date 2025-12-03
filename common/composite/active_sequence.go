package composite

import (
	"context"

	"github.com/jbcpollak/greenstalk/v2/core"
)

// ActiveSequence ticks each child in order. Returns success if
// all children succeed in one tick, else returns the status of
// the non-succeeding node. Restarts iteration the next tick.
func ActiveSequence(children ...core.Node) core.Node {
	return ActiveSequenceNamed("ActiveSequence", children...)
}

func ActiveSequenceNamed(name string, children ...core.Node) core.Node {
	base := core.NewComposite(core.BaseParams(name), children)
	return &activeSequence{Composite: base}
}

type activeSequence struct {
	core.Composite[core.BaseParams]
}

func (s *activeSequence) Activate(ctx context.Context, evt core.Event) core.ResultDetails {
	// No distinction between activation and ticking
	return s.Tick(ctx, evt)
}

func (s *activeSequence) Tick(ctx context.Context, evt core.Event) core.ResultDetails {
	for i := 0; i < len(s.Children); i++ {
		result := core.Update(ctx, s.Children[i], evt)
		if result.Status() != core.StatusSuccess {
			return result
		}
	}
	return core.SuccessResult()
}

func (s *activeSequence) Leave(context.Context) error {
	return nil
}

var _ core.Node = (*activeSequence)(nil)
