package condition

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/jbcpollak/greenstalk/v2/core"
)

type ParamSwitchFunc[T comparable] func() T

// ParamSwitch activates the child at the map key returned by the switch function.
// Note that an If node is a rudimentary form of a ParamSwitch node with two children
// and the function returning 0 / 1 for true / false.
func ParamSwitchNamed[T comparable](name string, switchFunc ParamSwitchFunc[T], children map[T]core.Node) core.Node {
	base := core.NewComposite(core.BaseParams(name), slices.Collect(maps.Values(children)))
	return &paramSwitchNode[T]{Composite: base, switchFunc: switchFunc, children: children}
}

func ParamSwitch[T comparable](switchFunc ParamSwitchFunc[T], children map[T]core.Node) core.Node {
	return ParamSwitchNamed("ParamSwitch", switchFunc, children)
}

type paramSwitchNode[T comparable] struct {
	core.Composite[core.BaseParams]
	switchFunc      ParamSwitchFunc[T]
	children        map[T]core.Node
	currentChildKey T
}

func (i *paramSwitchNode[T]) Activate(ctx context.Context, evt core.Event) core.ResultDetails {
	switchKey := i.switchFunc()
	if _, ok := i.children[switchKey]; !ok {
		return core.ErrorResult(fmt.Errorf("Switch key does not exist: %v", switchKey))
	}

	i.currentChildKey = switchKey

	return i.Tick(ctx, evt)
}

func (s *paramSwitchNode[T]) Tick(ctx context.Context, evt core.Event) core.ResultDetails {
	child := s.children[s.currentChildKey]
	return core.Update(ctx, child, evt)
}

func (s *paramSwitchNode[T]) Leave(context.Context) error {
	return nil
}

var _ core.Node = (*paramSwitchNode[any])(nil)
