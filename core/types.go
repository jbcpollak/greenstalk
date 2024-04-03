package core

import (
	"context"
	"fmt"
)

// Category denotes whether a node is a composite, decorator or leaf.
type Category string

// A list of behavior tree node categories.
const (
	CategoryInvalid   = Category("invalid")
	CategoryComposite = Category("composite")
	CategoryDecorator = Category("decorator")
	CategoryLeaf      = Category("leaf")
)

type NodeResult interface {
	Status() Status
}

// Status denotes the return value of the execution of a node.
type Status int

// A list of possible statuses.
const (
	StatusInvalid Status = iota
	StatusSuccess
	StatusFailure
	StatusRunning
	StatusError
)

func (s Status) Status() Status { return s }

type EnqueueFn func(Event) error
type NodeAsyncRunning func(ctx context.Context, enqueue EnqueueFn) error

func (NodeAsyncRunning) Status() Status { return StatusRunning }

type NodeRuntimeError struct {
	Err error
}

func (n NodeRuntimeError) Status() Status { return StatusError }
func (n NodeRuntimeError) String() string { return n.Err.Error() }

type BaseParams string

func (b BaseParams) Name() string {
	return string(b)
}

type EmptyReturns struct {
}

type (
	// Params denotes a list of parameters to a node.
	// Obsolete, do not use on new nodes
	DefaultParams map[string]interface{}

	// Returns is just a type alias for Params.
	// Obsolete, do not use on new nodes
	Returns = DefaultParams
)

func (p DefaultParams) Name() (string, error) {
	return p.GetString("name")
}

func (p DefaultParams) Get(key string) (any, error) {
	val, ok := p[key]
	if !ok {
		return 0, ErrParamNotFound(key)
	}
	return val, nil
}

func (p DefaultParams) GetInt(key string) (int, error) {
	val, ok := p[key]
	if !ok {
		return 0, ErrParamNotFound(key)
	}
	n, ok := val.(int)
	if !ok {
		return 0, ErrInvalidType(key)
	}
	return n, nil
}

func (p DefaultParams) GetString(key string) (string, error) {
	val, ok := p[key]
	if !ok {
		return "", ErrParamNotFound(key)
	}
	s, ok := val.(string)
	if !ok {
		return "", ErrInvalidType(key)
	}
	return s, nil
}

func ErrParamNotFound(name string) error {
	return fmt.Errorf("parameter %s not found", name)
}

func ErrInvalidType(name string) error {
	return fmt.Errorf("invalid type for %s", name)
}
