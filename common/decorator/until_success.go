package decorator

import (
	"github.com/jbcpollak/greenstalk/core"
)

// UntilSuccess updates its child until it returns Success.
func UntilSuccess[Blackboard any](child core.Node[Blackboard]) core.Node[Blackboard] {

	untilSuccess := func(status core.NodeResult) bool {
		return status == core.StatusSuccess
	}

	return RepeatUntil(RepeatUntilParams{
		BaseParams: "UntilSuccess",
		Until:      untilSuccess,
	}, child)
}
