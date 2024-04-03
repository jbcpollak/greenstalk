package action

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
	"github.com/rs/zerolog/log"
)

type SignallerParams[T any] struct {
	core.BaseParams

	Channel chan T
	Signal  T
}

// Sends a Signal on the provided channel
func Signaller[Blackboard any, T any](params SignallerParams[T]) *signaller[Blackboard, T] {
	base := core.NewLeaf[Blackboard](params, core.EmptyReturns{})
	return &signaller[Blackboard, T]{Leaf: base, params: params}
}

// succeed ...
type signaller[Blackboard any, T any] struct {
	core.Leaf[Blackboard, SignallerParams[T], core.EmptyReturns]

	params SignallerParams[T]
}

// Enter ...
func (a *signaller[Blackboard, T]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	log.Info().Msgf("%s: Signalling", a.Name())

	a.params.Channel <- a.params.Signal
	return core.StatusSuccess
}

func (a *signaller[Blackboard, T]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.NodeResult {
	return core.StatusError
}

// Leave ...
func (a *signaller[Blackboard, T]) Leave(bb Blackboard) {}
