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
		Tag:       function.Tag,
		Arguments: argumentDescs,
		Results:   resultDescs,
		Hooks:     hookDescs,
		Body:      function.Body,
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
		argumentDesc.ValueID = argument.ValueID
		argumentDesc.Value = reflect.ValueOf(argument.ValuePtr).Elem()
		argumentDesc.IsOptional = argument.IsOptional
	}
	return argumentDescs, nil
}

func validateArgument(argument *Argument, tag string) error {
	if argument.ValueID == "" {
		return fmt.Errorf("%w: empty value id; tag=%q", ErrInvalidArgument, tag)
	}
	if argument.ValuePtr == nil {
		return fmt.Errorf("%w: no value pointer; tag=%q valueID=%q",
			ErrInvalidArgument, tag, argument.ValueID)
	}
	valuePtr := reflect.ValueOf(argument.ValuePtr)
	if valuePtr.Kind() != reflect.Ptr {
		return fmt.Errorf("%w: invalid value pointer type; tag=%q valueID=%q valuePtrType=%T",
			ErrInvalidArgument, tag, argument.ValueID, argument.ValuePtr)
	}
	if valuePtr.IsNil() {
		return fmt.Errorf("%w: nil value pointer; tag=%q valueID=%q",
			ErrInvalidArgument, tag, argument.ValueID)
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
		resultDesc.ValueID = result.ValueID
		resultDesc.Value = reflect.ValueOf(result.ValuePtr).Elem()
		resultDesc.CleanupPtr = result.CleanupPtr
	}
	return resultDescs, nil
}

func validateResult(result *Result, tag string) error {
	if result.ValueID == "" {
		return fmt.Errorf("%w: empty value id; tag=%q", ErrInvalidResult, tag)
	}
	if result.ValuePtr == nil {
		return fmt.Errorf("%w: no value pointer; tag=%q valueID=%q",
			ErrInvalidResult, tag, result.ValueID)
	}
	valuePtr := reflect.ValueOf(result.ValuePtr)
	if valuePtr.Kind() != reflect.Ptr {
		return fmt.Errorf("%w: invalid value pointer type; tag=%q valueID=%q valuePtrType=%T",
			ErrInvalidResult, tag, result.ValueID, result.ValuePtr)
	}
	if valuePtr.IsNil() {
		return fmt.Errorf("%w: nil value pointer; tag=%q valueID=%q",
			ErrInvalidResult, tag, result.ValueID)
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
		hookDesc.ValueID = hook.ValueID
		hookDesc.Value = reflect.ValueOf(hook.ValuePtr).Elem()
		hookDesc.CallbackPtr = hook.CallbackPtr
	}
	return hookDescs, nil
}

func validateHook(hook *Hook, tag string) error {
	if hook.ValueID == "" {
		return fmt.Errorf("%w: empty value id; tag=%q", ErrInvalidHook, tag)
	}
	if hook.ValuePtr == nil {
		return fmt.Errorf("%w: no value pointer; tag=%q valueID=%q",
			ErrInvalidHook, tag, hook.ValueID)
	}
	valuePtr := reflect.ValueOf(hook.ValuePtr)
	if valuePtr.Kind() != reflect.Ptr {
		return fmt.Errorf("%w: invalid value pointer type; tag=%q valueID=%q valuePtrType=%T",
			ErrInvalidHook, tag, hook.ValueID, hook.ValuePtr)
	}
	if valuePtr.IsNil() {
		return fmt.Errorf("%w: nil value pointer; tag=%q valueID=%q",
			ErrInvalidHook, tag, hook.ValueID)
	}
	if hook.CallbackPtr == nil {
		return fmt.Errorf("%w: nil callback pointer; tag=%q valueID=%q",
			ErrInvalidHook, tag, hook.ValueID)
	}
	return nil
}

var (
	ErrInvalidFunction = errors.New("di: invalid function")
	ErrInvalidArgument = errors.New("di: invalid augment")
	ErrInvalidResult   = errors.New("di: invalid result")
	ErrInvalidHook     = errors.New("di: invalid hook")
)
