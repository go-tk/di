package di

import "context"

// Function represents the container that describes the specification for dependency injection.
// Tag serves as additional information to locate the Function for debugging purposes.
// Body will be called alone with Program.Run().
// Cleanup is optional and will be called along with Program.Clean().
// If CleanupPtr is not nil, Body should provision the cleanup.
type Function struct {
	Tag        string
	Arguments  []Argument
	Results    []Result
	Hooks      []Hook
	CleanupPtr *func()
	Body       func(ctx context.Context) (err error)
}

// Argument describes a value the container requires.
// Function.Body could use the in-value.
type Argument struct {
	InValueID  string
	InValuePtr interface{}
	IsOptional bool
}

// Result describes a value the container provides.
// Function.Body should provision the out-value.
type Result struct {
	OutValueID  string
	OutValuePtr interface{}
}

// Hook describes a callback which should be called once the value is created by a Function
// (as an out-value) but has not yet been passed to other Functions (as in-values).
// Function.Body should provision the callback.
// The callback could use the in-value.
type Hook struct {
	InValueID   string
	InValuePtr  interface{}
	CallbackPtr *func(ctx context.Context) (err error)
}
