package greenstalk

import (
	"github.com/jbcpollak/greenstalk/v2/core"
)

// TreeOption is used to set options when initializing a BehaviorTree.
// NewBehaviorTree() can accept a variable number of options.
//
// Example usage:
//
//	p := NewBehaviorTree(root, WithInput(someInput), WithOutput(someOutput))
type TreeOption func(*Tree)

// WithVisitor lets you specify a visitor which is called after every tick and visits every node.
func WithVisitors(v ...core.Visitor) TreeOption {
	return func(p *Tree) {
		p.visitors = v
	}
}
