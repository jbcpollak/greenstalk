package decorator

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
)

// Inverter ...
func InverterNamed(name string, child core.Node) core.Node {
	base := core.NewDecorator(core.BaseParams(name), child)
	return &inverter{Decorator: base}
}

func Inverter(child core.Node) core.Node {
	return InverterNamed("Inverter", child)
}

// inverter ...
type inverter struct {
	core.Decorator[core.BaseParams]
}

func (d *inverter) Activate(ctx context.Context, evt core.Event) core.ResultDetails {
	return d.Tick(ctx, evt)
}

// Tick ...
func (d *inverter) Tick(ctx context.Context, evt core.Event) core.ResultDetails {
	switch result := core.Update(ctx, d.Child, evt); result.Status() {
	case core.StatusSuccess:
		return core.FailureResult()
	case core.StatusFailure:
		return core.SuccessResult()
	default:
		return result
	}
}

// Leave ...
func (d *inverter) Leave(context.Context) error {
	return nil
}

var _ core.Node = (*inverter)(nil)
