package condition

import (
	"context"
	"fmt"

	"github.com/jbcpollak/greenstalk/v2/core"
)

type SwitchFunc func() int

// Switch activates the child at the index returned by the switch function.
// Note that an If node is a rudimentary form of a Switch node with two children
// and the function returning 0 / 1 for true / false.
func SwitchNamed(name string, switchFunc SwitchFunc, children ...core.Node) core.Node {
	base := core.NewComposite(core.BaseParams(name), children)
	return &switchNode{Composite: base, switchFunc: switchFunc}
}

func Switch(switchFunc SwitchFunc, children ...core.Node) core.Node {
	return SwitchNamed("Switch", switchFunc, children...)
}

type switchNode struct {
	core.Composite[core.BaseParams]
	switchFunc SwitchFunc
}

func (i *switchNode) Activate(ctx context.Context, evt core.Event) core.ResultDetails {
	switchIx := i.switchFunc()
	if switchIx < 0 || switchIx >= len(i.Children) {
		return core.ErrorResult(fmt.Errorf("Switch index out of bounds: %d", switchIx))
	}

	i.CurrentChild = switchIx

	return i.Tick(ctx, evt)
}

func (s *switchNode) Tick(ctx context.Context, evt core.Event) core.ResultDetails {
	child := s.Children[s.CurrentChild]
	return core.Update(ctx, child, evt)
}

func (s *switchNode) Leave(context.Context) error {
	return nil
}

var _ core.Node = (*switchNode)(nil)
