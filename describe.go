package di

import (
	"errors"
	"fmt"
	"reflect"
)

func describeFunction(function *Function) (functionDesc, error) {
	if err := validateFunction(function); err != nil {
		return functionDesc{}, err
	}
	argumentDescs, err := describeArguments(function.Arguments, function.Tag)
	if err != nil {
		return functionDesc{}, err
	}
	resultDescs, err := describeResults(function.Results, function.Tag)
	if err != nil {
		return functionDesc{}, err
	}
	hookDescs, err := describeHooks(function.Hooks, function.Tag)
	if err != nil {
		return functionDesc{}, err
	}
	return functionDesc{
		Tag:        function.Tag,
		Arguments:  argumentDescs,
		Results:    resultDescs,
		Hooks:      hookDescs,
		CleanupPtr: function.CleanupPtr,
		Body:       function.Body,
	}, nil
}

func validateFunction(function *Function) error {
	if function.Tag == "" {
		return fmt.Errorf("%w: empty tag", ErrInvalidFunction)
	}
	if function.Body == nil {
		return fmt.Errorf("%w: nil body; tag=%q", ErrInvalidFunction, function.Tag)
	}
	return nil
}

func describeArguments(arguments []Argument, tag string) ([]argumentDesc, error) {
	argumentDescs := make([]argumentDesc, len(arguments))
	for i := range arguments {
		argument := &arguments[i]
		if err := validateArgument(argument, tag); err != nil {
			return nil, err
		}
		argumentDesc := &argumentDescs[i]
		argumentDesc.InValueID = argument.InValueID
		argumentDesc.InValue = reflect.ValueOf(argument.InValuePtr).Elem()
		argumentDesc.IsOptional = argument.IsOptional
	}
	return argumentDescs, nil
}

func validateArgument(argument *Argument, tag string) error {
	if argument.InValueID == "" {
		return fmt.Errorf("%w: empty in-value id; tag=%q", errInvalidArgument, tag)
	}
	if argument.InValuePtr == nil {
		return fmt.Errorf("%w: no in-value pointer; tag=%q inValueID=%q",
			errInvalidArgument, tag, argument.InValueID)
	}
	inValuePtr := reflect.ValueOf(argument.InValuePtr)
	if inValuePtr.Kind() != reflect.Ptr {
		return fmt.Errorf("%w: invalid in-value pointer type; tag=%q inValueID=%q inValuePtrType=%q",
			errInvalidArgument, tag, argument.InValueID, reflect.TypeOf(argument.InValuePtr))
	}
	if inValuePtr.IsNil() {
		return fmt.Errorf("%w: nil in-value pointer; tag=%q inValueID=%q",
			errInvalidArgument, tag, argument.InValueID)
	}
	return nil
}

func describeResults(results []Result, tag string) ([]resultDesc, error) {
	resultDescs := make([]resultDesc, len(results))
	for i := range results {
		result := &results[i]
		if err := validateResult(result, tag); err != nil {
			return nil, err
		}
		resultDesc := &resultDescs[i]
		resultDesc.OutValueID = result.OutValueID
		resultDesc.OutValue = reflect.ValueOf(result.OutValuePtr).Elem()
	}
	return resultDescs, nil
}

func validateResult(result *Result, tag string) error {
	if result.OutValueID == "" {
		return fmt.Errorf("%w: empty out-value id; tag=%q", errInvalidResult, tag)
	}
	if result.OutValuePtr == nil {
		return fmt.Errorf("%w: no out-value pointer; tag=%q outValueID=%q",
			errInvalidResult, tag, result.OutValueID)
	}
	outValuePtr := reflect.ValueOf(result.OutValuePtr)
	if outValuePtr.Kind() != reflect.Ptr {
		return fmt.Errorf("%w: invalid out-value pointer type; tag=%q outValueID=%q outValuePtrType=%q",
			errInvalidResult, tag, result.OutValueID, reflect.TypeOf(result.OutValuePtr))
	}
	if outValuePtr.IsNil() {
		return fmt.Errorf("%w: nil out-value pointer; tag=%q outValueID=%q",
			errInvalidResult, tag, result.OutValueID)
	}
	return nil
}

func describeHooks(hooks []Hook, tag string) ([]hookDesc, error) {
	hookDescs := make([]hookDesc, len(hooks))
	for i := range hooks {
		hook := &hooks[i]
		if err := validateHook(hook, tag); err != nil {
			return nil, err
		}
		hookDesc := &hookDescs[i]
		hookDesc.InValueID = hook.InValueID
		hookDesc.InValue = reflect.ValueOf(hook.InValuePtr).Elem()
		hookDesc.CallbackPtr = hook.CallbackPtr
	}
	return hookDescs, nil
}

func validateHook(hook *Hook, tag string) error {
	if hook.InValueID == "" {
		return fmt.Errorf("%w: empty in-value id; tag=%q", errInvalidHook, tag)
	}
	if hook.InValuePtr == nil {
		return fmt.Errorf("%w: no in-value pointer; tag=%q inValueID=%q",
			errInvalidHook, tag, hook.InValueID)
	}
	inValuePtr := reflect.ValueOf(hook.InValuePtr)
	if inValuePtr.Kind() != reflect.Ptr {
		return fmt.Errorf("%w: invalid in-value pointer type; tag=%q inValueID=%q inValuePtrType=%q",
			errInvalidHook, tag, hook.InValueID, reflect.TypeOf(hook.InValuePtr))
	}
	if inValuePtr.IsNil() {
		return fmt.Errorf("%w: nil in-value pointer; tag=%q inValueID=%q",
			errInvalidHook, tag, hook.InValueID)
	}
	if hook.CallbackPtr == nil {
		return fmt.Errorf("%w: nil callback pointer; tag=%q inValueID=%q",
			errInvalidHook, tag, hook.InValueID)
	}
	return nil
}

var (
	// ErrInvalidFunction is returned by Program.AddFunctions() when a Function is invalid.
	ErrInvalidFunction = errors.New("di: invalid function")

	errInvalidArgument = fmt.Errorf("%w: invalid augment", ErrInvalidFunction)
	errInvalidResult   = fmt.Errorf("%w: invalid result", ErrInvalidFunction)
	errInvalidHook     = fmt.Errorf("%w: invalid hook", ErrInvalidFunction)
)
