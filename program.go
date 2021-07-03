package di

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

type Program struct {
	functionDescs              []functionDesc
	orderedFunctionDescIndexes []int
}

func (p *Program) AddFunction(function Function) error {
	functionDesc, err := describeFunction(&function)
	if err != nil {
		return err
	}
	functionDesc.Index = len(p.functionDescs)
	p.functionDescs = append(p.functionDescs, functionDesc)
	return nil
}

func (p *Program) Run(ctx context.Context) error {
	if err := p.resolve(); err != nil {
		return err
	}
	if n, err := p.callFunctions(ctx); err != nil {
		p.clean(n)
		return err
	}
	return nil
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
			return i, fmt.Errorf("function failed; tag=%q: %w", functionDesc.Tag, err)
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
					return i, fmt.Errorf("callback failed; tag=%q valueID=%q: %w",
						tag, resultDesc.ValueID, err)
				}
			}
		}
	}
	return i, nil
}

func (p *Program) Clean() {
	p.clean(len(p.orderedFunctionDescIndexes))
}

func (p *Program) clean(n int) {
	for i := n - 1; i >= 0; i-- {
		functionDescIndex := p.orderedFunctionDescIndexes[i]
		functionDesc := &p.functionDescs[functionDescIndex]
		for j := range functionDesc.Results {
			resultDesc := &functionDesc.Results[j]
			if resultDesc.CleanupPtr != nil {
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
