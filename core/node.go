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

// Preliminary interface to work around intermediate types like
// composite, decorator, etc not inplementing Enter/Tick/Leave
type Walkable[Blackboard any] interface {
	// Automatically implemented by embedding a pointer to a
	// Composite, Decorator or Leaf node in the custom node.
	Status() Status
	SetStatus(Status)
	Id() uuid.UUID
	Name() string
	Category() Category
	String() string

	Walk(WalkFunc[Blackboard], int)
}

type WalkFunc[Blackboard any] func(Walkable[Blackboard], int)

// The Node interface must be satisfied by any custom node.
type Node[Blackboard any] interface {
	Walkable[Blackboard]

	// Must be implemented by the custom node.
	Activate(context.Context, Blackboard, Event) NodeResult
	Tick(context.Context, Blackboard, Event) NodeResult
	Leave(Blackboard)
}

type Params interface {
	Name() string
}

// BaseNode contains properties shared by all categories of node.
// Do not use this type directly.
type BaseNode[P Params] struct {
	id       uuid.UUID
	category Category
	status   Status
	Params   P
}

func newBaseNode[P Params](category Category, params P) BaseNode[P] {
	return BaseNode[P]{
		id:       uuid.New(),
		category: category,
		Params:   params,
	}
}

// Status returns the status of this node.
func (n *BaseNode[P]) Id() uuid.UUID {
	return n.id
}

// Status returns the status of this node.
func (n *BaseNode[P]) Name() string {
	return n.Params.Name()
}

// Status returns the status of this node.
func (n *BaseNode[P]) Status() Status {
	return n.status
}

// SetStatus sets the status of this node.
func (n *BaseNode[P]) SetStatus(status Status) {
	n.status = status
}

// GetCategory returns the category of this node.
func (n *BaseNode[P]) Category() Category {
	return n.category
}
