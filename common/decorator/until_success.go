package decorator

import (
	"github.com/jbcpollak/greenstalk/v2/core"
)

// UntilSuccess updates its child until it returns Success.
func UntilSuccessNamed(name string, child core.Node) core.Node {
	untilSuccess := func(result core.ResultDetails) bool {
		return result.Status() == core.StatusSuccess
	}

	return RepeatUntil(RepeatUntilParams{
		BaseParams: core.BaseParams(name),
		Until:      untilSuccess,
	}, child)
}

func UntilSuccess(child core.Node) core.Node {
	return UntilSuccessNamed("UntilSuccess", child)
}
