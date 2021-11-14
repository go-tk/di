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
	calledFunctionCount        int
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
	if err := p.callFunctions(ctx); err != nil {
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

func (p *Program) callFunctions(ctx context.Context) error {
	for _, functionDescIndex := range p.orderedFunctionDescIndexes {
		functionDesc := &p.functionDescs[functionDescIndex]
		if err := p.callFunction(ctx, functionDesc); err != nil {
			return err
		}
		if err := checkCallbacks(functionDesc); err != nil {
			return err
		}
		if err := checkCleanup(functionDesc); err != nil {
			return err
		}
		if err := p.doCallbacks(ctx, functionDesc); err != nil {
			return err
		}
	}
	return nil
}

func (p *Program) callFunction(ctx context.Context, functionDesc *functionDesc) error {
	for i := range functionDesc.Arguments {
		argumentDesc := &functionDesc.Arguments[i]
		if argumentDesc.Result == nil {
			continue
		}
		argumentDesc.InValue.Set(argumentDesc.Result.OutValue)
	}
	err := functionDesc.Body(ctx)
	p.calledFunctionCount++
	if err != nil {
		return fmt.Errorf("di: function failed; tag=%q: %w", functionDesc.Tag, err)
	}
	return nil
}

func checkCallbacks(functionDesc *functionDesc) error {
	for i := range functionDesc.Hooks {
		hookDesc := &functionDesc.Hooks[i]
		if *hookDesc.CallbackPtr == nil {
			return fmt.Errorf("%w; tag=%q inValueID=%q",
				ErrNilCallback, functionDesc.Tag, hookDesc.InValueID)
		}
	}
	return nil
}

func checkCleanup(functionDesc *functionDesc) error {
	if cleanupPtr := functionDesc.CleanupPtr; cleanupPtr != nil && *cleanupPtr == nil {
		return fmt.Errorf("%w; tag=%q", ErrNilCleanup, functionDesc.Tag)
	}
	return nil
}

func (p *Program) doCallbacks(ctx context.Context, functionDesc *functionDesc) error {
	for i := range functionDesc.Results {
		resultDesc := &functionDesc.Results[i]
		for _, hookDesc := range resultDesc.Hooks {
			hookDesc.InValue.Set(resultDesc.OutValue)
			if err := (*hookDesc.CallbackPtr)(ctx); err != nil {
				tag := p.functionDescs[hookDesc.FunctionIndex].Tag
				return fmt.Errorf("di: callback failed; tag=%q inValueID=%q: %w",
					tag, hookDesc.InValueID, err)
			}
		}
	}
	return nil
}

// Clean does cleanups associated with Functions.
func (p *Program) Clean() {
	for i := p.calledFunctionCount - 1; i >= 0; i-- {
		functionDescIndex := p.orderedFunctionDescIndexes[i]
		functionDesc := &p.functionDescs[functionDescIndex]
		doCleanup(functionDesc)
	}
}

func doCleanup(functionDesc *functionDesc) {
	if cleanupPtr := functionDesc.CleanupPtr; cleanupPtr != nil && *cleanupPtr != nil {
		(*cleanupPtr)()
	}
}

type functionDesc struct {
	Tag        string
	Arguments  []argumentDesc
	Results    []resultDesc
	Hooks      []hookDesc
	CleanupPtr *func()
	Body       func(context.Context) error
	Index      int
}

type argumentDesc struct {
	InValueID  string
	InValue    reflect.Value
	IsOptional bool
	Result     *resultDesc
}

type resultDesc struct {
	OutValueID    string
	OutValue      reflect.Value
	FunctionIndex int
	Hooks         []*hookDesc
}

type hookDesc struct {
	InValueID     string
	InValue       reflect.Value
	CallbackPtr   *func(context.Context) error
	FunctionIndex int
}

var (
	// ErrNilCallback is return by Program.Run() when the callback pointed by CallbackPtr
	// is nil after the Function body is executed.
	ErrNilCallback = errors.New("di: nil callback")

	// ErrNilCleanup is return by Program.Run() when the cleanup pointed by CleanupPtr
	// is nil after the Function body is executed.
	ErrNilCleanup = errors.New("di: nil cleanup")
)
