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
		p.doCleanups(n)
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
		if err := callFunction(ctx, functionDesc); err != nil {
			return i, err
		}
		i++
		if err := checkCleanups(functionDesc); err != nil {
			return i, err
		}
		if err := checkCallbacks(functionDesc); err != nil {
			return i, err
		}
		if err := p.doCallbacks(ctx, functionDesc); err != nil {
			return i, err
		}
	}
	return i, nil
}

func callFunction(ctx context.Context, functionDesc *functionDesc) error {
	for i := range functionDesc.Arguments {
		argumentDesc := &functionDesc.Arguments[i]
		if argumentDesc.Result == nil {
			continue
		}
		argumentDesc.InValue.Set(argumentDesc.Result.OutValue)
	}
	if err := functionDesc.Body(ctx); err != nil {
		return fmt.Errorf("di: function failed; tag=%q: %w", functionDesc.Tag, err)
	}
	return nil
}

func checkCleanups(functionDesc *functionDesc) error {
	for i := range functionDesc.Results {
		resultDesc := &functionDesc.Results[i]
		if resultDesc.CleanupPtr != nil && *resultDesc.CleanupPtr == nil {
			return fmt.Errorf("%w; tag=%q outValueID=%q",
				ErrNilCleanup, functionDesc.Tag, resultDesc.OutValueID)
		}
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

func (p *Program) doCallbacks(ctx context.Context, functionDesc *functionDesc) error {
	for j := range functionDesc.Results {
		resultDesc := &functionDesc.Results[j]
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

// Clean does cleanups associated with Results.
func (p *Program) Clean() {
	p.doCleanups(len(p.orderedFunctionDescIndexes))
}

func (p *Program) doCleanups(n int) {
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
	InValueID  string
	InValue    reflect.Value
	IsOptional bool
	Result     *resultDesc
}

type resultDesc struct {
	OutValueID    string
	OutValue      reflect.Value
	CleanupPtr    *func()
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
	// ErrNilCleanup is return by Program.Run() when the cleanup pointed by CleanupPtr
	// is nil after the Function body is executed.
	ErrNilCleanup = errors.New("di: nil cleanup")

	// ErrNilCallback is return by Program.Run() when the callback pointed by CallbackPtr
	// is nil after the Function body is executed.
	ErrNilCallback = errors.New("di: nil callback")
)
