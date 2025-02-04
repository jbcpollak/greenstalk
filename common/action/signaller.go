package action

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
)

type SignallerParams[T any] struct {
	core.BaseParams

	Channel chan T
	Signal  T
}

// Sends a Signal on the provided channel
func Signaller[Blackboard any, T any](params SignallerParams[T]) *function_action[Blackboard] {
	fap := FunctionActionParams{
		Func: func(ctx context.Context) core.ResultDetails {
			// TODO: FunctionAction should pass some information to the function
			// internal.Logger.Info("Signalling", "name", a.Name())

			params.Channel <- params.Signal
			return core.SuccessResult()
		},
	}
	base := core.NewLeaf[Blackboard](fap)
	return &function_action[Blackboard]{Leaf: base}
}
