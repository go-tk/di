package di

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

// Program represents a program consists of Functions.
type Program struct {
	functionDescs              []functionDesc
	orderedFunctionDescIndexes []int
}

// AddFunction adds a Function into the Program.
func (p *Program) AddFunction(function Function) error {
	functionDesc, err := describeFunction(&function)
	if err != nil {
		return err
	}
	functionDesc.Index = len(p.functionDescs)
	p.functionDescs = append(p.functionDescs, functionDesc)
	return nil
}

// MustAddFunction wraps AddFunction and panics when an error occurs.
func (p *Program) MustAddFunction(function Function) {
	if err := p.AddFunction(function); err != nil {
		panic(err)
	}
}

// Run arranges Functions basing on dependency analysis and calls them in order.
func (p *Program) Run(ctx context.Context) error {
	if err := p.resolve(); err != nil {
		return err
	}
	if n, err := p.callFunctions(ctx); err != nil {
		p.doClean(n)
		return err
	}
	return nil
}

// MustRun wraps Run and panics when an error occurs.
func (p *Program) MustRun(ctx context.Context) {
	if err := p.Run(ctx); err != nil {
		panic(err)
	}
}

func (p *Program) resolve() error {
	resolution12 := resolution12{
		FunctionDescs: p.functionDescs,
	}
	if err := resolution12.ExecutePhase1(); err != nil {
		return err
	}
	if err := resolution12.ExecutePhase2(); err != nil {
		return err
	}
	resolution3 := resolution3{
		FunctionDescs: p.functionDescs,
	}
	if err := resolution3.ExecutePhase3(); err != nil {
		return err
	}
	p.orderedFunctionDescIndexes = resolution3.OrderedFunctionDescIndexes
	return nil
}

func (p *Program) callFunctions(ctx context.Context) (int, error) {
	var i int
	for n := len(p.orderedFunctionDescIndexes); i < n; {
		functionDescIndex := p.orderedFunctionDescIndexes[i]
		functionDesc := &p.functionDescs[functionDescIndex]
		for j := range functionDesc.Arguments {
			argumentDesc := &functionDesc.Arguments[j]
			if argumentDesc.Result == nil {
				continue
			}
			argumentDesc.Value.Set(argumentDesc.Result.Value)
		}
		if err := functionDesc.Body(ctx); err != nil {
			return i, fmt.Errorf("di: function failed; tag=%q: %w", functionDesc.Tag, err)
		}
		i++
		for j := range functionDesc.Results {
			resultDesc := &functionDesc.Results[j]
			if resultDesc.CleanupPtr != nil && *resultDesc.CleanupPtr == nil {
				return i, fmt.Errorf("%w; tag=%q valueID=%q",
					ErrNilCleanup, functionDesc.Tag, resultDesc.ValueID)
			}
		}
		for j := range functionDesc.Hooks {
			hookDesc := &functionDesc.Hooks[j]
			if *hookDesc.CallbackPtr == nil {
				return i, fmt.Errorf("%w; tag=%q valueID=%q",
					ErrNilCallback, functionDesc.Tag, hookDesc.ValueID)
			}
		}
		for j := range functionDesc.Results {
			resultDesc := &functionDesc.Results[j]
			for _, hookDesc := range resultDesc.Hooks {
				hookDesc.Value.Set(resultDesc.Value)
				if err := (*hookDesc.CallbackPtr)(ctx); err != nil {
					tag := p.functionDescs[hookDesc.FunctionIndex].Tag
					return i, fmt.Errorf("di: callback failed; tag=%q valueID=%q: %w",
						tag, resultDesc.ValueID, err)
				}
			}
		}
	}
	return i, nil
}

// Clean does cleanups associated with Results.
func (p *Program) Clean() {
	p.doClean(len(p.orderedFunctionDescIndexes))
}

func (p *Program) doClean(n int) {
	for i := n - 1; i >= 0; i-- {
		functionDescIndex := p.orderedFunctionDescIndexes[i]
		functionDesc := &p.functionDescs[functionDescIndex]
		for j := range functionDesc.Results {
			resultDesc := &functionDesc.Results[j]
			if resultDesc.CleanupPtr != nil && *resultDesc.CleanupPtr != nil {
				(*resultDesc.CleanupPtr)()
			}
		}
	}
}

type functionDesc struct {
	Tag       string
	Arguments []argumentDesc
	Results   []resultDesc
	Hooks     []hookDesc
	Body      func(context.Context) error
	Index     int
}

type argumentDesc struct {
	ValueID    string
	Value      reflect.Value
	IsOptional bool
	Result     *resultDesc
}

type resultDesc struct {
	ValueID       string
	Value         reflect.Value
	CleanupPtr    *func()
	FunctionIndex int
	Hooks         []*hookDesc
}

type hookDesc struct {
	ValueID       string
	Value         reflect.Value
	CallbackPtr   *func(context.Context) error
	FunctionIndex int
}

var (
	ErrNilCleanup  = errors.New("di: nil cleanup")
	ErrNilCallback = errors.New("di: nil callback")
)
