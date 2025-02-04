package state

import (
	"context"

	"github.com/jbcpollak/greenstalk/common/action"
	"github.com/jbcpollak/greenstalk/core"
)

// This node resets all provided states and returns SuccessStatus
func MakeStateResetAction[Blackboard any](states ...StateResetter) core.Node[Blackboard] {
	return action.FunctionAction[Blackboard](action.FunctionActionParams{
		BaseParams: "stateReset",
		Func: func(ctx context.Context) core.ResultDetails {
			for _, state := range states {
				state.Reset()
			}
			return core.SuccessResult()
		},
	})
}
