package state

import (
	"github.com/jbcpollak/greenstalk/v2/common/action"
	"github.com/jbcpollak/greenstalk/v2/core"
)

// This node resets all provided states and returns SuccessStatus
func MakeStateResetAction(states ...StateResetter) core.Node {
	return action.FunctionAction(action.FunctionActionParams{
		BaseParams: "stateReset",
		Func: func() core.ResultDetails {
			for _, state := range states {
				state.Reset()
			}
			return core.SuccessResult()
		},
	})
}
