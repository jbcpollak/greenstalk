package decorator

import (
	"github.com/jbcpollak/greenstalk/core"
)

// UntilFailure updates its child until it returns Failure.
func UntilFailureNamed(name string, child core.Node) core.Node {
	untilFailure := func(result core.ResultDetails) bool {
		return result.Status() == core.StatusFailure
	}

	return RepeatUntil(RepeatUntilParams{
		BaseParams: core.BaseParams(name),
		Until:      untilFailure,
	}, child)
}

func UntilFailure(child core.Node) core.Node {
	return UntilFailureNamed("UntilFailure", child)
}
