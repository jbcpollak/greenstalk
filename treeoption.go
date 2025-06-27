package greenstalk

import (
	"context"

	"github.com/jbcpollak/greenstalk/core"
)

// TreeOption is used to set options when initializing a BehaviorTree.
// NewBehaviorTree() can accept a variable number of options.
//
// Example usage:
//
//	p := NewBehaviorTree(root, WithInput(someInput), WithOutput(someOutput))
type TreeOption[Blackboard any] func(*behaviorTree[Blackboard])

// WithVisitor lets you specify a visitor which is called after every tick and visits every node.
func WithVisitors[Blackboard any](v ...core.Visitor[Blackboard]) TreeOption[Blackboard] {
	return func(p *behaviorTree[Blackboard]) {
		p.visitors = v
	}
}

func WithContext[Blackboard any](ctx context.Context) TreeOption[Blackboard] {
	return func(p *behaviorTree[Blackboard]) {
		p.ctx = ctx
	}
}
