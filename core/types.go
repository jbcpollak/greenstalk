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

type ResultDetails interface {
	Status() Status
}

type simpleResultDetails struct {
	status Status
}

func (d simpleResultDetails) Status() Status { return d.status }

func SuccessResult() ResultDetails {
	return simpleResultDetails{status: StatusSuccess}
}
func FailureResult() ResultDetails {
	return simpleResultDetails{status: StatusFailure}
}
func InvalidResult() ResultDetails {
	return simpleResultDetails{status: StatusInvalid}
}
func RunningResult() ResultDetails {
	return simpleResultDetails{status: StatusRunning}
}

type EnqueueFn func(Event) error
type RunningFn func(ctx context.Context, enqueue EnqueueFn) error

func InitRunningResult(fn RunningFn) InitRunningResultDetails {
	return InitRunningResultDetails{fn}
}

type InitRunningResultDetails struct {
	RunningFn RunningFn
}

func (InitRunningResultDetails) Status() Status { return StatusRunning }

func ErrorResult(err error) ErrorResultDetails {
	return ErrorResultDetails{err}
}

type ErrorResultDetails struct {
	Err error
}

func (n ErrorResultDetails) Status() Status { return StatusError }
func (n ErrorResultDetails) Error() error   { return n.Err }

type BaseParams string

func (b BaseParams) Name() string {
	return string(b)
}

type EmptyReturns struct{}

type EmptyBlackboard struct{}

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
