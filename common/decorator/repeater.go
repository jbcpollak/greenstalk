package decorator

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
	"github.com/rs/zerolog/log"
)

type RepeaterParams struct {
	N int
}

// Repeater updates its child n times, at which point the repeater
// returns Success. The repeater runs forever if n == 0.
func Repeater[Blackboard any](params RepeaterParams, child core.Node[Blackboard]) core.Node[Blackboard] {
	base := core.NewDecorator(core.DecoratorParams{BaseParams: "Repeater"}, child)
	d := &repeater[Blackboard]{Decorator: base}

	d.n = params.N
	return d
}

type repeater[Blackboard any] struct {
	core.Decorator[Blackboard]
	n int
	i int
}

func (d *repeater[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	d.i = 0

	return d.Tick(ctx, bb, evt)
}

func (d *repeater[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	log.Info().Msg("Repeater: Calling child")
	status := core.Update(ctx, d.Child, bb, evt)

	if status == core.StatusRunning {
		return core.StatusRunning
	}

	// Run forever if n == 0.
	if d.n == 0 {
		return core.StatusRunning
	}

	d.i++
	if d.i < d.n {
		return core.StatusRunning
	}

	// At this point, the repeater has updated its child n times.
	return core.StatusSuccess
}

func (d *repeater[Blackboard]) Leave(bb Blackboard) {}
