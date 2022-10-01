package di

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// Program consists of DI Functions which are containers for dependency injection.
type Program struct {
	functions             []function
	arguments             []argument
	results               []result
	hooks                 []hook
	sortedFunctionIndexes []int
	calledFunctionCount   int
}

type function struct {
	Index           int
	Name            string
	ArgumentIndexes []int
	ResultIndexes   []int
	Body            func(context.Context) error
	HookIndexes     []int
	Cleanup         func()
}

// FunctionBuilder is the type of function that constructs a DI Function.
type FunctionBuilder func(function *function, program *Program) (err error)

// NewFunction add a DI Function into the Program.
func (p *Program) NewFunction(functionBuilders ...FunctionBuilder) error {
	pc, _, _, _ := runtime.Caller(1)
	functionName := runtime.FuncForPC(pc).Name()
	return p.doNewFunction(functionName, functionBuilders...)
}

// MustNewFunction likes NewFunction but panics when an error occurs.
func (p *Program) MustNewFunction(functionBuilders ...FunctionBuilder) {
	pc, _, _, _ := runtime.Caller(1)
	functionName := runtime.FuncForPC(pc).Name()
	if err := p.doNewFunction(functionName, functionBuilders...); err != nil {
		panic(fmt.Sprintf("new function: %v", err))
	}
}

func (p *Program) doNewFunction(functionName string, functionBuilders ...FunctionBuilder) (returnedErr error) {
	functionIndex := len(p.functions)
	p.functions = append(p.functions, function{Index: functionIndex})
	defer func() {
		if returnedErr != nil {
			p.functions = p.functions[:functionIndex]
		}
	}()
	function := &p.functions[functionIndex]
	function.Name = functionName
	for _, functionBuilder := range functionBuilders {
		if err := functionBuilder(function, p); err != nil {
			return err
		}
	}
	if function.Body == nil {
		return fmt.Errorf("%w; functionName=%q", ErrBodyRequired, functionName)
	}
	return nil
}

// ErrBodyRequired is returned by Program.NewFunction() when no body is specified.
var ErrBodyRequired = errors.New("di: body required")

type argument struct {
	FunctionIndex    int
	ValueRef         string
	ValueReceiver    reflect.Value
	IsOptional       bool
	ResultIndex      int
	ReceiveValueAddr bool
}

// Argument specifies an argument for a DI Function.
func Argument(valueRef string, rawValueReceiverPtr interface{}) FunctionBuilder {
	return argument1(valueRef, rawValueReceiverPtr, false)
}

// OptionalArgument specifies an optional argument for a DI Function.
func OptionalArgument(valueRef string, rawValueReceiverPtr interface{}) FunctionBuilder {
	return argument1(valueRef, rawValueReceiverPtr, true)
}

func argument1(valueRef string, rawValueReceiverPtr interface{}, isOptional bool) FunctionBuilder {
	return func(function *function, program *Program) error {
		if valueRef == "" {
			return fmt.Errorf("%w: empty value ref; functionName=%q", ErrInvalidArgument, function.Name)
		}
		if rawValueReceiverPtr == nil {
			return fmt.Errorf("%w: no value receiver; functionName=%q valueRef=%q", ErrInvalidArgument, function.Name, valueRef)
		}
		valueReceiverPtr := reflect.ValueOf(rawValueReceiverPtr)
		if valueReceiverPtr.Kind() != reflect.Ptr {
			return fmt.Errorf("%w: invalid value receiver pointer; valueReceiverPtrType=%q functionName=%q valueRef=%q",
				ErrInvalidArgument, valueReceiverPtr.Type(), function.Name, valueRef)
		}
		if valueReceiverPtr.IsNil() {
			return fmt.Errorf("%w: no value receiver; functionName=%q valueRef=%q", ErrInvalidArgument, function.Name, valueRef)
		}
		argumentIndex := len(program.arguments)
		program.arguments = append(program.arguments, argument{})
		argument := &program.arguments[argumentIndex]
		argument.FunctionIndex = function.Index
		argument.ValueRef = valueRef
		argument.ValueReceiver = valueReceiverPtr.Elem()
		argument.IsOptional = isOptional
		argument.ResultIndex = -1
		function.ArgumentIndexes = append(function.ArgumentIndexes, argumentIndex)
		return nil
	}
}

// ErrInvalidArgument is returned by Program.NewFunction() when an invalid argument is specified.
var ErrInvalidArgument = errors.New("di: invalid augment")

type result struct {
	FunctionIndex int
	ValueName     string
	Value         reflect.Value
	HookIndexes   []int
}

// Result specifies a result for a DI Function.
func Result(valueName string, rawValuePtr interface{}) FunctionBuilder {
	return func(function *function, program *Program) error {
		if valueName == "" {
			return fmt.Errorf("%w: empty value name; functionName=%q", ErrInvalidResult, function.Name)
		}
		if rawValuePtr == nil {
			return fmt.Errorf("%w: no value; functionName=%q valueName=%q", ErrInvalidResult, function.Name, valueName)
		}
		valuePtr := reflect.ValueOf(rawValuePtr)
		if valuePtr.Kind() != reflect.Ptr {
			return fmt.Errorf("%w: invalid value pointer; valuePtrType=%q functionName=%q valueName=%q",
				ErrInvalidResult, valuePtr.Type(), function.Name, valueName)
		}
		if valuePtr.IsNil() {
			return fmt.Errorf("%w: no value; functionName=%q valueName=%q", ErrInvalidResult, function.Name, valueName)
		}
		resultIndex := len(program.results)
		program.results = append(program.results, result{})
		result := &program.results[resultIndex]
		result.FunctionIndex = function.Index
		result.ValueName = valueName
		result.Value = valuePtr.Elem()
		function.ResultIndexes = append(function.ResultIndexes, resultIndex)
		return nil
	}
}

// ErrInvalidResult is returned by Program.NewFunction() when an invalid result is specified.
var ErrInvalidResult = errors.New("di: invalid result")

// Body specifies the body for a DI Function.
func Body(body func(context.Context) error) FunctionBuilder {
	return func(function *function, program *Program) error {
		if body == nil {
			return fmt.Errorf("%w; functionName=%q", ErrNilBody, function.Name)
		}
		function.Body = body
		return nil
	}
}

// ErrNilBody is returned by Program.NewFunction() when nil body is specified.
var ErrNilBody = errors.New("di: nil body")

// Cleanup specifies the cleanup for a DI Function.
func Cleanup(cleanup func()) FunctionBuilder {
	return func(function *function, program *Program) error {
		if cleanup == nil {
			return fmt.Errorf("%w; functionName=%q", ErrNilCleanup, function.Name)
		}
		function.Cleanup = cleanup
		return nil
	}
}

// ErrNilCleanup is returned by Program.NewFunction() when nil cleanup is specified.
var ErrNilCleanup = errors.New("di: nil cleanup")

type hook struct {
	FunctionIndex    int
	ValueRef         string
	ValueReceiver    reflect.Value
	Callback         func(context.Context) error
	ReceiveValueAddr bool
}

// Hook specifies a hook for a DI Function.
func Hook(valueRef string, rawValueReceiverPtr interface{}, callback func(context.Context) error) FunctionBuilder {
	return func(function *function, program *Program) error {
		if valueRef == "" {
			return fmt.Errorf("%w: empty value ref; functionName=%q", ErrInvalidHook, function.Name)
		}
		if rawValueReceiverPtr == nil {
			return fmt.Errorf("%w: no value receiver; functionName=%q valueRef=%q", ErrInvalidHook, function.Name, valueRef)
		}
		valueReceiverPtr := reflect.ValueOf(rawValueReceiverPtr)
		if valueReceiverPtr.Kind() != reflect.Ptr {
			return fmt.Errorf("%w: invalid value receiver pointer; valueReceiverPtrType=%q functionName=%q valueRef=%q",
				ErrInvalidHook, valueReceiverPtr.Type(), function.Name, valueRef)
		}
		if valueReceiverPtr.IsNil() {
			return fmt.Errorf("%w: no value receiver; functionName=%q valueRef=%q", ErrInvalidHook, function.Name, valueRef)
		}
		if callback == nil {
			return fmt.Errorf("%w: nil callback; functionName=%q valueRef=%q", ErrInvalidHook, function.Name, valueRef)
		}
		hookIndex := len(program.hooks)
		program.hooks = append(program.hooks, hook{})
		hook := &program.hooks[hookIndex]
		hook.FunctionIndex = function.Index
		hook.ValueRef = valueRef
		hook.ValueReceiver = valueReceiverPtr.Elem()
		hook.Callback = callback
		function.HookIndexes = append(function.HookIndexes, hookIndex)
		return nil
	}
}

// ErrInvalidHook is returned by Program.NewFunction() when an invalid hook is specified.
var ErrInvalidHook = errors.New("di: invalid hook")

// Run calls all DI Functions added into the Program, the order in which DI Functions are to be called
// is based on dependency analysis.
func (p *Program) Run(ctx context.Context) error {
	if err := p.resolve(); err != nil {
		return err
	}
	if err := p.sortFunctions(); err != nil {
		return err
	}
	return p.callFunctions(ctx)
}

func (p *Program) resolve() error {
	valueName2ResultIndex := make(map[string]int, len(p.results))
	for resultIndex := range p.results {
		result := &p.results[resultIndex]
		if resultIndex2, ok := valueName2ResultIndex[result.ValueName]; ok {
			result2 := &p.results[resultIndex2]
			return fmt.Errorf("%w; valueName=%q functionName1=%q functionName2=%q",
				ErrDuplicateValueName, result.ValueName, p.functions[result.FunctionIndex].Name,
				p.functions[result2.FunctionIndex].Name)
		}
		valueName2ResultIndex[result.ValueName] = resultIndex
	}
	for argumentIndex := range p.arguments {
		argument := &p.arguments[argumentIndex]
		resultIndex, ok := valueName2ResultIndex[argument.ValueRef]
		if !ok {
			if argument.IsOptional {
				continue
			}
			return fmt.Errorf("%w; valueRef=%q functionName=%q",
				ErrValueNotFound, argument.ValueRef, p.functions[argument.FunctionIndex].Name)
		}
		result := &p.results[resultIndex]
		valueType := result.Value.Type()
		valueReceiverType := argument.ValueReceiver.Type()
		if valueReceiverType == reflect.PtrTo(valueType) {
			argument.ReceiveValueAddr = true
		} else {
			if !valueType.AssignableTo(valueReceiverType) {
				return fmt.Errorf("%w; valueReceiverType=%q valueType=%q valueRef=%q functionName=%q",
					ErrIncompatibleValueReceiver, valueReceiverType, valueType, argument.ValueRef,
					p.functions[argument.FunctionIndex].Name)
			}
		}
		argument.ResultIndex = resultIndex
	}
	for hookIndex := range p.hooks {
		hook := &p.hooks[hookIndex]
		resultIndex, ok := valueName2ResultIndex[hook.ValueRef]
		if !ok {
			return fmt.Errorf("%w; valueRef=%q functionName=%q",
				ErrValueNotFound, hook.ValueRef, p.functions[hook.FunctionIndex].Name)
		}
		result := &p.results[resultIndex]
		valueType := result.Value.Type()
		valueReceiverType := hook.ValueReceiver.Type()
		if valueReceiverType == reflect.PtrTo(valueType) {
			hook.ReceiveValueAddr = true
		} else {
			if !valueType.AssignableTo(valueReceiverType) {
				return fmt.Errorf("%w; valueReceiverType=%q valueType=%q valueRef=%q functionName=%q",
					ErrIncompatibleValueReceiver, valueReceiverType, valueType, hook.ValueRef,
					p.functions[hook.FunctionIndex].Name)
			}
		}
		result.HookIndexes = append(result.HookIndexes, hookIndex)
	}
	return nil
}

func (p *Program) sortFunctions() error {
	var walk func(*function, interface{}) bool
	var path []interface{}
	visitedFunctionIndexes := make(map[int]struct{}, len(p.functions))
	walk = func(function *function, from interface{}) bool {
		functionIndex := function.Index
		if _, ok := visitedFunctionIndexes[functionIndex]; ok {
			return true
		}
		pathLength := len(path)
		if from == nil {
			path = append(path, function)
		} else {
			path = append(path, from, function)
		}
		if functionIndex < 0 {
			return false
		}
		function.Index = -1
		for _, argumentIndex := range function.ArgumentIndexes {
			argument := &p.arguments[argumentIndex]
			if argument.ResultIndex < 0 {
				continue
			}
			result := &p.results[argument.ResultIndex]
			function2 := &p.functions[result.FunctionIndex]
			if !walk(function2, argument) {
				return false
			}
		}
		for _, resultIndex := range function.ResultIndexes {
			result := &p.results[resultIndex]
			for _, hookIndex := range result.HookIndexes {
				hook := &p.hooks[hookIndex]
				function2 := &p.functions[hook.FunctionIndex]
				if !walk(function2, hook) {
					return false
				}
			}
		}
		function.Index = functionIndex
		path = path[:pathLength]
		visitedFunctionIndexes[functionIndex] = struct{}{}
		p.sortedFunctionIndexes = append(p.sortedFunctionIndexes, functionIndex)
		return true
	}
	dumpPath := func() string {
		var builder strings.Builder
		n := len(path)
		for i, j := 0, n-1; i < j; i += 2 {
			function, from := path[i].(*function), path[i+1]
			switch from := from.(type) {
			case *argument:
				builder.WriteString(fmt.Sprintf("%s@argument:%s => ", function.Name, from.ValueRef))
			case *hook:
				builder.WriteString(fmt.Sprintf("%s@hook:%s => ", function.Name, from.ValueRef))
			default:
				panic("unreachable code")
			}
		}
		function := path[n-1].(*function)
		builder.WriteString(function.Name)
		return builder.String()
	}
	for functionIndex := range p.functions {
		function := &p.functions[functionIndex]
		if !walk(function, nil) {
			return fmt.Errorf("%w; path=%q", ErrCircularDependencies, dumpPath())
		}
	}
	return nil
}

func (p *Program) callFunctions(ctx context.Context) error {
	for _, functionIndex := range p.sortedFunctionIndexes {
		function := &p.functions[functionIndex]
		for _, argumentIndex := range function.ArgumentIndexes {
			argument := &p.arguments[argumentIndex]
			if argument.ResultIndex < 0 {
				continue
			}
			result := &p.results[argument.ResultIndex]
			if argument.ReceiveValueAddr {
				argument.ValueReceiver.Set(result.Value.Addr())
			} else {
				argument.ValueReceiver.Set(result.Value)
			}
		}
		if err := function.Body(ctx); err != nil {
			return fmt.Errorf("call function; functionName=%q: %w", function.Name, err)
		}
		p.calledFunctionCount++
		for _, resultIndex := range function.ResultIndexes {
			result := &p.results[resultIndex]
			for _, hookIndex := range result.HookIndexes {
				hook := &p.hooks[hookIndex]
				if hook.ReceiveValueAddr {
					hook.ValueReceiver.Set(result.Value.Addr())
				} else {
					hook.ValueReceiver.Set(result.Value)
				}
				if err := hook.Callback(ctx); err != nil {
					function2 := &p.functions[hook.FunctionIndex]
					return fmt.Errorf("do callback; functionName=%q valueRef=%q: %w", function2.Name, hook.ValueRef, err)
				}
			}
		}
	}
	return nil
}

var (
	// ErrDuplicateValueName is returned by Program.Run() when a value name used by Result() is duplicate.
	ErrDuplicateValueName = errors.New("di: duplicate value name")

	// ErrValueNotFound is returned by Program.Run() when a value used by Argument()/Hook() does not exist.
	ErrValueNotFound = errors.New("di: value not found")

	// ErrIncompatibleValueReceiver is returned by Program.Run() when a value receiver used by Argument()/Hook() is incompatible.
	ErrIncompatibleValueReceiver = errors.New("di: incompatible value receiver")

	// ErrCircularDependencies is returned by Program.Run() when circular dependencies are detected.
	ErrCircularDependencies = errors.New("di: circular dependencies")
)

// MustRun likes Run but panics when an error occurs.
func (p *Program) MustRun(ctx context.Context) {
	if err := p.Run(ctx); err != nil {
		panic(fmt.Sprintf("run program: %v", err))
	}
}

// Clean calls cleanups of DI Functions, the order in which cleanups are to be called
// is reversed to the order in which DI Functions are called.
func (p *Program) Clean() {
	for i := p.calledFunctionCount - 1; i >= 0; i-- {
		functionIndex := p.sortedFunctionIndexes[i]
		function := &p.functions[functionIndex]
		if cleanup := function.Cleanup; cleanup != nil {
			cleanup()
		}
	}
}
