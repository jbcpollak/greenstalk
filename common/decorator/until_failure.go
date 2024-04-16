package decorator

import (
	"github.com/jbcpollak/greenstalk/core"
)

// UntilFailure updates its child until it returns Failure.
func UntilFailure[Blackboard any](child core.Node[Blackboard]) core.Node[Blackboard] {

	untilFailure := func(result core.ResultDetails) bool {
		return result.Status() == core.StatusFailure
	}

	return RepeatUntil(RepeatUntilParams{
		BaseParams: "UntilFailure",
		Until:      untilFailure,
	}, child)
}
