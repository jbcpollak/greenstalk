package composite

import (
	"context"
	"math/rand"

	"github.com/jbcpollak/greenstalk/core"
)

// RandomSelector creates a new random selector node.
func RandomSelectorNamed(name string, children ...core.Node) core.Node {
	base := core.NewComposite(core.BaseParams(name), children)
	return &randomSelector{Composite: base}
}

func RandomSelector(children ...core.Node) core.Node {
	return RandomSelectorNamed("RandomSelector", children...)
}

// randomSelector ...
type randomSelector struct {
	core.Composite[core.BaseParams]
}

// Activate ...
func (s *randomSelector) Activate(ctx context.Context, evt core.Event) core.ResultDetails {
	return s.Tick(ctx, evt)
}

// Tick ...
func (s *randomSelector) Tick(ctx context.Context, evt core.Event) core.ResultDetails {
	index := rand.Intn(len(s.Children))
	child := s.Children[index]
	return core.Update(ctx, child, evt)
}

// Leave ...
func (s *randomSelector) Leave(context.Context) error {
	return nil
}

var _ core.Node = (*randomSelector)(nil)
