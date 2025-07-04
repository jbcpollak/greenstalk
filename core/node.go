package core

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const NAME_PREFIX_SEPARATOR = "."

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
	FullName() string
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
	SetNamePrefix(string)
}

type Params interface {
	Name() string
}

// BaseNode contains properties shared by all categories of node.
// Do not use this type directly.
type BaseNode[P Params] struct {
	id         uuid.UUID
	category   Category
	result     ResultDetails
	Params     P
	namePrefix string
}

func newBaseNode[P Params](category Category, params P) BaseNode[P] {
	if strings.Contains(params.Name(), NAME_PREFIX_SEPARATOR) {
		err := fmt.Errorf("Node '%s' name may not contain node name separator string '%s'", params.Name(), NAME_PREFIX_SEPARATOR)
		panic(err)
	}
	return BaseNode[P]{
		id:         uuid.New(),
		category:   category,
		result:     InvalidResult(),
		Params:     params,
		namePrefix: "",
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

func (n *BaseNode[P]) FullName() string {
	return n.namePrefix + n.Name()
}

func (n *BaseNode[P]) SetNamePrefix(name string) {
	n.namePrefix = name + NAME_PREFIX_SEPARATOR
}
