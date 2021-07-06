package di

import "context"

// Function represents the container that describes the specification for dependency injection.
type Function struct {
	Tag       string
	Arguments []Argument
	Results   []Result
	Hooks     []Hook
	Body      func(ctx context.Context) (err error)
}

// Argument describes a value the container requires.
type Argument struct {
	InValueID  string
	InValuePtr interface{}
	IsOptional bool
}

// Result describes a value the container provides.
// It is optional to do a cleanup for the value when Program.Clean() is executed.
type Result struct {
	OutValueID  string
	OutValuePtr interface{}
	CleanupPtr  *func()
}

// Hook describes a callback which should be called once the value is created but
// has not yet been provisioned to other Functions.
type Hook struct {
	InValueID   string
	InValuePtr  interface{}
	CallbackPtr *func(ctx context.Context) (err error)
}
