package decorator

import (
	"github.com/jbcpollak/greenstalk/core"
)

// UntilFailure updates its child until it returns Failure.
func UntilFailureNamed[Blackboard any](name string, child core.Node[Blackboard]) core.Node[Blackboard] {

	untilFailure := func(result core.ResultDetails) bool {
		return result.Status() == core.StatusFailure
	}

	return RepeatUntil(RepeatUntilParams{
		BaseParams: core.BaseParams(name),
		Until:      untilFailure,
	}, child)
}
func UntilFailure[Blackboard any](child core.Node[Blackboard]) core.Node[Blackboard] {
	return UntilFailureNamed("UntilFailure", child)
}
