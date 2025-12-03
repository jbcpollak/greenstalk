package core

import (
	"fmt"
)

// Leaf is the base type for any specific leaf node (domain-specific).
// Each leaf node has Params: data keys that the implementation imports
// and Returns: data keys that the implementation exports.
type Leaf[P Params] struct {
	BaseNode[P]
}

// NewLeaf creates a new leaf base node.
// TODO: change Params to interface and save it?
func NewLeaf[P Params](params P) Leaf[P] {
	return Leaf[P]{
		BaseNode: newBaseNode(CategoryLeaf, params),
	}
}

func (c *Leaf[Params]) Walk(walkFn WalkFunc, level int) {
	walkFn(c, level)
}

// String returns a string representation of the leaf node.
func (a *Leaf[Params]) String() string {
	return fmt.Sprintf("! %s (%v)",
		a.Params.Name(),
		a.Params,
	)
}
