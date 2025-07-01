package decorator

import (
	"github.com/jbcpollak/greenstalk/core"
)

// UntilSuccess updates its child until it returns Success.
func UntilSuccessNamed[Blackboard any](name string, child core.Node[Blackboard]) core.Node[Blackboard] {

	untilSuccess := func(result core.ResultDetails) bool {
		return result.Status() == core.StatusSuccess
	}

	return RepeatUntil(RepeatUntilParams{
		BaseParams: core.BaseParams(name),
		Until:      untilSuccess,
	}, child)
}
func UntilSuccess[Blackboard any](child core.Node[Blackboard]) core.Node[Blackboard] {
	return UntilSuccessNamed("UntilSuccess", child)
}
