package composite

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
)

// Parallel updates all its children in parallel, i.e. every frame.
// It does not retry on nodes that have failed or succeeded.
//
// success/failReq is the minimum amount of nodes required to
// succeed/fail for the parallel sequence node itself to succeed/fail.
// A value of 0 for either node means that all nodes must succeed/fail.
func ParallelNamed[Blackboard any](name string, successReq, failReq int, children ...core.Node[Blackboard]) core.Node[Blackboard] {
	base := core.NewComposite(core.BaseParams(name), children)
	if successReq == 0 {
		successReq = len(children)
	}
	if failReq == 0 {
		failReq = len(children)
	}
	return &parallel[Blackboard]{
		base,
		successReq,
		failReq,
		0,
		0,
		make([]bool, len(children)),
	}
}

func Parallel[Blackboard any](successReq, failReq int, children ...core.Node[Blackboard]) core.Node[Blackboard] {
	return ParallelNamed("Parallel", successReq, failReq, children...)
}

type parallel[Blackboard any] struct {
	core.Composite[Blackboard, core.BaseParams]
	successReq int
	failReq    int
	succeeded  int
	failed     int
	completed  []bool
}

func (s *parallel[Blackboard]) Activate(ctx context.Context, bb Blackboard, evt core.Event) core.ResultDetails {
	s.succeeded = 0
	s.failed = 0

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
			s.succeeded++
			s.completed[i] = true
		case core.StatusFailure:
			s.failed++
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

	if s.succeeded >= s.successReq {
		return core.SuccessResult()
	}
	if s.failed >= s.failReq {
		return core.FailureResult()
	}

	if len(runningResultDetails) > 0 {
		return core.InitRunningResultsCollection(runningResultDetails)
	} else {
		return core.RunningResult()
	}
}

func (s *parallel[Blackboard]) Leave(context.Context, Blackboard) error {
	return nil
}

var _ core.Node[any] = (*parallel[any])(nil)
