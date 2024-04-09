package decorator

import (
	"github.com/jbcpollak/greenstalk/core"
)

// UntilFailure updates its child until it returns Failure.
func UntilFailure[Blackboard any](child core.Node[Blackboard]) core.Node[Blackboard] {

	untilFailure := func(status core.NodeResult) bool {
		return status == core.StatusFailure
	}

	return RepeatUntil(RepeatUntilParams{
		BaseParams: "UntilFailure",
		Until:      untilFailure,
	}, child)
}
