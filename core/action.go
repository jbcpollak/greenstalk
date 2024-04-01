package core

import (
	"fmt"
)

// Leaf is the base type for any specific leaf node (domain-specific).
// Each leaf node has Params: data keys that the implementation imports
// and Returns: data keys that the implementation exports.
type Leaf[Blackboard any, P Params, Returns any] struct {
	BaseNode
	Params  P
	Returns Returns
}

// NewLeaf creates a new leaf base node.
// TODO: change Params to interface and save it?
func NewLeaf[Blackboard any, P Params, Returns any](params P, returns Returns) Leaf[Blackboard, P, Returns] {
	return Leaf[Blackboard, P, Returns]{
		BaseNode: newBaseNode(CategoryLeaf, params),
		Params:   params,
		Returns:  returns,
	}
}

func (c *Leaf[Blackboard, Params, Returns]) Walk(walkFn WalkFunc[Blackboard], level int) {
	walkFn(c, level)
}

// String returns a string representation of the leaf node.
func (a *Leaf[Blackboard, Params, Returns]) String() string {
	return fmt.Sprintf("! %s (%v : %v)",
		a.name,
		a.Params,
		a.Returns,
	)
}
