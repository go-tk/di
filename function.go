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
	ValueID    string
	ValuePtr   interface{}
	IsOptional bool
}

// Argument describes a value the container provides.
// It is optional to do a cleanup for the value when Program.Clean() is executed.
type Result struct {
	ValueID    string
	ValuePtr   interface{}
	CleanupPtr *func()
}

// Hook describes a callback which should be called once the value is created but
// has not yet been provisioned to other Functions.
type Hook struct {
	ValueID     string
	ValuePtr    interface{}
	CallbackPtr *func(ctx context.Context) (err error)
}
