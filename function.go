package di

import "context"

type Function struct {
	Tag       string
	Arguments []Argument
	Results   []Result
	Hooks     []Hook
	Body      func(ctx context.Context) (err error)
}

type Argument struct {
	ValueID    string
	ValuePtr   interface{}
	IsOptional bool
}

type Result struct {
	ValueID    string
	ValuePtr   interface{}
	CleanupPtr *func()
}

type Hook struct {
	ValueID     string
	ValuePtr    interface{}
	CallbackPtr *func(ctx context.Context) (err error)
}
