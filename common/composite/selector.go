package composite

import (
	"context"

	"github.com/jbcpollak/greenstalk/v2/core"
)

// Selector updates each child in order, returning success as soon as
// a child succeeds. If a child returns Running, the selector node
// will resume execution from that child the next tick.
func SelectorNamed(name string, children ...core.Node) core.Node {
	base := core.NewComposite(core.BaseParams(name), children)
	return &selector{Composite: base}
}

func Selector(children ...core.Node) core.Node {
	return SelectorNamed("Selector", children...)
}

type selector struct {
	core.Composite[core.BaseParams]
}

func (s *selector) Activate(ctx context.Context, evt core.Event) core.ResultDetails {
	s.CurrentChild = 0

	// Tick as expected
	return s.Tick(ctx, evt)
}

func (s *selector) Tick(ctx context.Context, evt core.Event) core.ResultDetails {
	for s.CurrentChild < len(s.Children) {
		result := core.Update(ctx, s.Children[s.CurrentChild], evt)
		if result.Status() != core.StatusFailure {
			return result
		}
		s.CurrentChild++
	}
	return core.FailureResult()
}

func (s *selector) Leave(context.Context) error {
	return nil
}

var _ core.Node = (*selector)(nil)
