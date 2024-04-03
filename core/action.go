package core

import (
	"fmt"
)

// Leaf is the base type for any specific leaf node (domain-specific).
// Each leaf node has Params: data keys that the implementation imports
// and Returns: data keys that the implementation exports.
type Leaf[Blackboard any, P Params] struct {
	BaseNode
	Params P
}

// NewLeaf creates a new leaf base node.
// TODO: change Params to interface and save it?
func NewLeaf[Blackboard any, P Params](params P) Leaf[Blackboard, P] {
	return Leaf[Blackboard, P]{
		BaseNode: newBaseNode(CategoryLeaf, params),
		Params:   params,
	}
}

func (c *Leaf[Blackboard, Params]) Walk(walkFn WalkFunc[Blackboard], level int) {
	walkFn(c, level)
}

// String returns a string representation of the leaf node.
func (a *Leaf[Blackboard, Params]) String() string {
	return fmt.Sprintf("! %s (%v)",
		a.name,
		a.Params,
	)
}
