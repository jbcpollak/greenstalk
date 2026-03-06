package condition

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/jbcpollak/greenstalk/v2/core"
)

type SwitchFunc[T comparable] func() T

// SwitchMap activates the child at the map key returned by the switch function.
// Note that an If node is a rudimentary form of a SwitchMap node with two children
// and the function returning 0 / 1 for true / false.
func SwitchMapNamed[T comparable](name string, switchFunc SwitchFunc[T], children map[T]core.Node) core.Node {
	base := core.NewComposite(core.BaseParams(name), slices.Collect(maps.Values(children)))
	return &switchMapNode[T]{Composite: base, switchFunc: switchFunc, children: children}
}

func SwitchMap[T comparable](switchFunc SwitchFunc[T], children map[T]core.Node) core.Node {
	return SwitchMapNamed("SwitchMap", switchFunc, children)
}

type switchMapNode[T comparable] struct {
	core.Composite[core.BaseParams]
	switchFunc      SwitchFunc[T]
	children        map[T]core.Node
	currentChildKey T
}

func (i *switchMapNode[T]) Activate(ctx context.Context, evt core.Event) core.ResultDetails {
	switchKey := i.switchFunc()
	if _, ok := i.children[switchKey]; !ok {
		return core.ErrorResult(fmt.Errorf("Switch key does not exist: %v", switchKey))
	}

	i.currentChildKey = switchKey

	return i.Tick(ctx, evt)
}

func (s *switchMapNode[T]) Tick(ctx context.Context, evt core.Event) core.ResultDetails {
	child := s.children[s.currentChildKey]
	return core.Update(ctx, child, evt)
}

func (s *switchMapNode[T]) Leave(context.Context) error {
	return nil
}

var _ core.Node = (*switchMapNode[any])(nil)
