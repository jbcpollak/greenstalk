package decorator

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
)

// Inverter ...
func InverterNamed[Blackboard any](name string, child core.Node[Blackboard]) core.Node[Blackboard] {
	base := core.NewDecorator(core.BaseParams(name), child)
	return &inverter[Blackboard]{Decorator: base}
}
func Inverter[Blackboard any](child core.Node[Blackboard]) core.Node[Blackboard] {
	return InverterNamed("Inverter", child)
}

// inverter ...
type inverter[Blackboard any] struct {
	core.Decorator[Blackboard, core.BaseParams]
}

func (d *inverter[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	return d.Tick(ctx, bb, evt)
}

// Tick ...
func (d *inverter[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	switch result := core.Update(ctx, d.Child, bb, evt); result.Status() {
	case core.StatusSuccess:
		return core.FailureResult()
	case core.StatusFailure:
		return core.SuccessResult()
	default:
		return result
	}
}

// Leave ...
func (d *inverter[Blackboard]) Leave(bb Blackboard) error {
	return nil
}
