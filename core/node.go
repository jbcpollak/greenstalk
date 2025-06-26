package core

import (
	"context"

	"github.com/google/uuid"
)

type Event interface {
	// UUID of the node that generated the event
	// Use uuid.Nil for events that are applicable to many nodes
	TargetNodeId() uuid.UUID
}

type DefaultEvent struct {
}

func (e DefaultEvent) TargetNodeId() uuid.UUID {
	return uuid.Nil
}

func TargetNodeEvent(id uuid.UUID) targetNodeEvent {
	return targetNodeEvent{targetNodeId: id}
}

type targetNodeEvent struct {
	targetNodeId uuid.UUID
}

func (e targetNodeEvent) TargetNodeId() uuid.UUID {
	return e.targetNodeId
}

type ErrorEvent struct {
	Err error
}

func (e ErrorEvent) TargetNodeId() uuid.UUID {
	return uuid.Nil
}

// Preliminary interface to work around intermediate types like
// composite, decorator, etc not implementing Enter/Tick/Leave
type Walkable[Blackboard any] interface {
	// Automatically implemented by embedding a pointer to a
	// Composite, Decorator or Leaf node in the custom node.
	Result() ResultDetails
	SetResult(ResultDetails)
	Id() uuid.UUID
	Name() string
	Category() Category
	String() string

	Walk(WalkFunc[Blackboard], int)
}

type Visitor[Blackboard any] func(Walkable[Blackboard])
type WalkFunc[Blackboard any] func(node Walkable[Blackboard], level int)

// The Node interface must be satisfied by any custom node.
type Node[Blackboard any] interface {
	Walkable[Blackboard]

	// Must be implemented by the custom node.
	Activate(context.Context, Blackboard, Event) ResultDetails
	Tick(context.Context, Blackboard, Event) ResultDetails
	Leave(Blackboard) error
}

type Params interface {
	Name() string
	SetName(string) Params
}

// BaseNode contains properties shared by all categories of node.
// Do not use this type directly.
type BaseNode[P Params] struct {
	id       uuid.UUID
	category Category
	result   ResultDetails
	Params   P
}

func newBaseNode[P Params](category Category, params P) BaseNode[P] {
	return BaseNode[P]{
		id:       uuid.New(),
		category: category,
		result:   InvalidResult(),
		Params:   params,
	}
}

// Status returns the status of this node.
func (n *BaseNode[P]) Id() uuid.UUID {
	return n.id
}

// Name returns the name of this node.
func (n *BaseNode[P]) Name() string {
	return n.Params.Name()
}

// SetName sets the name of this node
func (n *BaseNode[P]) SetName(newName string) *BaseNode[P] {
	n.Params = n.Params.SetName(newName).(P)
	return n
}

// Status returns the status of this node.
func (n *BaseNode[P]) Result() ResultDetails {
	return n.result
}

// SetResult sets the current result of this node.
func (n *BaseNode[P]) SetResult(result ResultDetails) {
	n.result = result
}

// GetCategory returns the category of this node.
func (n *BaseNode[P]) Category() Category {
	return n.category
}
