package composite

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
)

// Parallel updates all its children in parallel, i.e. every frame.
// It does not retry on nodes that have failed or succeeded.
//
// succ/failReq is the minimum amount of nodes required to
// succeed/fail for the parallel sequence node itself to succeed/fail.
// A value of 0 for either node means that all nodes must succeed/fail.
func ParallelNamed[Blackboard any](name string, succReq, failReq int, children ...core.Node[Blackboard]) core.Node[Blackboard] {
	base := core.NewComposite(core.BaseParams(name), children)
	if succReq == 0 {
		succReq = len(children)
	}
	if failReq == 0 {
		failReq = len(children)
	}
	return &parallel[Blackboard]{
		base,
		succReq,
		failReq,
		0,
		0,
		make([]bool, len(children)),
	}
}
func Parallel[Blackboard any](succReq, failReq int, children ...core.Node[Blackboard]) core.Node[Blackboard] {
	return ParallelNamed("Parallel", succReq, failReq, children...)
}

type parallel[Blackboard any] struct {
	core.Composite[Blackboard, core.BaseParams]
	succReq   int
	failReq   int
	succ      int
	fail      int
	completed []bool
}

func (s *parallel[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	s.succ = 0
	s.fail = 0

	for i := 0; i < len(s.Children); i++ {
		s.completed[i] = false
	}

	return s.Tick(ctx, bb, evt)
}

func (s *parallel[Blackboard]) Tick(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	runningResultDetails := []core.InitRunningResultDetails{}

	// Update every child that has not completed yet every tick.
	for i := 0; i < len(s.Children); i++ {

		// Ignore a child if has already succeeded or failed.
		if s.completed[i] {
			continue
		}

		// Update a child and count whether it succeeded or failed,
		// and mark it as completed in either of those two cases.
		result := core.Update(ctx, s.Children[i], bb, evt)
		status := result.Status()
		switch status {
		case core.StatusSuccess:
			s.succ++
			s.completed[i] = true
		case core.StatusFailure:
			s.fail++
			s.completed[i] = true
		case core.StatusRunning:
			if initRunningResult, ok := result.(core.InitRunningResultDetails); ok {
				runningResultDetails = append(runningResultDetails, initRunningResult)
			} else if initRunningResultsCollection, ok := result.(core.InitRunningResultsDetailsCollection); ok {
				runningResultDetails = append(runningResultDetails, initRunningResultsCollection.Results...)
			}
		case core.StatusError:
			// any errors are returned immediately so the whole tree can error out
			return result
		}
	}

	if s.succ >= s.succReq {
		return core.SuccessResult()
	}
	if s.fail >= s.failReq {
		return core.FailureResult()
	}

	if len(runningResultDetails) > 0 {
		return core.InitRunningResultsCollection(runningResultDetails)
	} else {
		return core.RunningResult()
	}
}

func (s *parallel[Blackboard]) Leave(bb Blackboard) error {
	return nil
}
