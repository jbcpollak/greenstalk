package composite

import (
	"context"
	"math/rand"

	"github.com/jbcpollak/greenstalk/core"
)

// RandomSelector creates a new random selector node.
func RandomSelectorNamed[Blackboard any](name string, children ...core.Node[Blackboard]) core.Node[Blackboard] {
	base := core.NewComposite(core.BaseParams(name), children)
	return &randomSelector[Blackboard]{Composite: base}
}
func RandomSelector[Blackboard any](children ...core.Node[Blackboard]) core.Node[Blackboard] {
	return RandomSelectorNamed("RandomSelector", children...)
}

// randomSelector ...
type randomSelector[Blackboard any] struct {
	core.Composite[Blackboard, core.BaseParams]
}

// Activate ...
func (s *randomSelector[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	return s.Tick(ctx, bb, evt)
}

// Tick ...
func (s *randomSelector[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	index := rand.Intn(len(s.Children))
	child := s.Children[index]
	return core.Update(ctx, child, bb, evt)
}

// Leave ...
func (s *randomSelector[Blackboard]) Leave(bb Blackboard) error {
	return nil
}
