package condition

import (
	"context"
	"fmt"

	"github.com/jbcpollak/greenstalk/core"
)

type SwitchFunc func() int

func Switch[Blackboard any](switchFunc SwitchFunc, children ...core.Node[Blackboard]) core.Node[Blackboard] {
	base := core.NewComposite(core.BaseParams("Switch"), children)
	return &switchNode[Blackboard]{Composite: base, switchFunc: switchFunc}
}

type switchNode[Blackboard any] struct {
	core.Composite[Blackboard, core.BaseParams]
	switchFunc SwitchFunc
}

func (i *switchNode[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	switchIx := i.switchFunc()
	if switchIx < 0 || switchIx >= len(i.Children) {
		return core.ErrorResult(fmt.Errorf("Switch index out of bounds: %d", switchIx))
	}

	i.CurrentChild = switchIx

	return i.Tick(ctx, bb, evt)
}

func (s *switchNode[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	child := s.Children[s.CurrentChild]
	return core.Update(ctx, child, bb, evt)
}

func (s *switchNode[Blackboard]) Leave(bb Blackboard) error {
	return nil
}
