package decorator

import (
	"context"

	"github.com/jbcpollak/greenstalk/v2/common/action"
	"github.com/jbcpollak/greenstalk/v2/core"
)

func WithAsyncNamed(
	name string,
	enterFunc func(context.Context) error,
	exitFunc func(context.Context) error,
	child core.Node,
) core.Node {
	enterNode := action.AsyncFunctionAction(action.AsyncFunctionActionParams{
		BaseParams: core.BaseParams("enter"),
		Func: func(ctx context.Context) core.ResultDetails {
			err := enterFunc(ctx)
			if err != nil {
				return core.ErrorResult(err)
			}
			return core.SuccessResult()
		},
	})

	exitNode := action.AsyncFunctionAction(action.AsyncFunctionActionParams{
		BaseParams: core.BaseParams("exit"),
		Func: func(ctx context.Context) core.ResultDetails {
			err := exitFunc(ctx)
			if err != nil {
				return core.ErrorResult(err)
			}
			return core.SuccessResult()
		},
	})

	base := core.NewComposite(core.BaseParams(name), []core.Node{enterNode, child, exitNode})
	return &asyncWithSequence{Composite: base}

}

func WithAsync(
	enterFunc func(context.Context) error,
	exitFunc func(context.Context) error,
	child core.Node,
) core.Node {
	return WithAsyncNamed("WithAsync", enterFunc, exitFunc, child)
}

// private composite node for AsyncWith that has exactly 3 children, executes them in order and then
// returns the status of the second child
type asyncWithSequence struct {
	core.Composite[core.BaseParams]
	result core.ResultDetails
}

func (s *asyncWithSequence) Activate(ctx context.Context, evt core.Event) core.ResultDetails {
	if len(s.Children) != 3 {
		panic("asyncWithSequence must have exactly 3 children")
	}
	s.CurrentChild = 0

	return s.Tick(ctx, evt)
}

func (s *asyncWithSequence) Tick(ctx context.Context, evt core.Event) core.ResultDetails {
	for s.CurrentChild < len(s.Children) {
		result := core.Update(ctx, s.Children[s.CurrentChild], evt)
		if result.Status() == core.StatusRunning || result.Status() == core.StatusError {
			return result
		}

		if s.CurrentChild == 1 {
			s.result = result
		}
		s.CurrentChild++
	}
	return s.result
}

func (s *asyncWithSequence) Leave(context.Context) error {
	return nil
}

var _ core.Node = (*asyncWithSequence)(nil)
