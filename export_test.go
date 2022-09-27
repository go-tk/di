package di

import (
	"bytes"
	"fmt"
	"reflect"
)

type Function = function

func (p *Program) Dump(buffer *bytes.Buffer) {
	for i := range p.functions {
		function := &p.functions[i]
		p.dumpFunction(i, function, buffer)
	}
	for i := range p.arguments {
		argument := &p.arguments[i]
		p.dumpArgument(i, argument, buffer)
	}
	for i := range p.results {
		result := &p.results[i]
		p.dumpResult(i, result, buffer)
	}
	for i := range p.hooks {
		hook := &p.hooks[i]
		p.dumpHook(i, hook, buffer)
	}
	fmt.Fprintf(buffer, "SortedFunctionIndexes: %v\n", p.sortedFunctionIndexes)
	fmt.Fprintf(buffer, "CalledFunctionCount: %v\n", p.calledFunctionCount)
}

func (p *Program) dumpFunction(i int, function *function, buffer *bytes.Buffer) {
	fmt.Fprintf(buffer, "Function[%d]:\n", i)
	fmt.Fprintf(buffer, "\tIndex: %v\n", function.Index)
	fmt.Fprintf(buffer, "\tName: %v\n", function.Name)
	fmt.Fprintf(buffer, "\tArgumentIndexes: %v\n", function.ArgumentIndexes)
	fmt.Fprintf(buffer, "\tResultIndexes: %v\n", function.ResultIndexes)
	fmt.Fprintf(buffer, "\tHasBody: %v\n", function.Body != nil)
	fmt.Fprintf(buffer, "\tHookIndexes: %v\n", function.HookIndexes)
	fmt.Fprintf(buffer, "\tHasCleanup: %v\n", function.Cleanup != nil)
}

func (p *Program) dumpArgument(i int, argument *argument, buffer *bytes.Buffer) {
	fmt.Fprintf(buffer, "Argument[%d]:\n", i)
	fmt.Fprintf(buffer, "\tFunctionIndex: %v\n", argument.FunctionIndex)
	fmt.Fprintf(buffer, "\tValueRef: %v\n", argument.ValueRef)
	fmt.Fprintf(buffer, "\tHasValueReceiver: %v\n", argument.ValueReceiver != (reflect.Value{}))
	fmt.Fprintf(buffer, "\tIsOptional: %v\n", argument.IsOptional)
	fmt.Fprintf(buffer, "\tResultIndex: %v\n", argument.ResultIndex)
	fmt.Fprintf(buffer, "\tReceiveValueAddr: %v\n", argument.ReceiveValueAddr)
}

func (p *Program) dumpResult(i int, result *result, buffer *bytes.Buffer) {
	fmt.Fprintf(buffer, "Result[%d]:\n", i)
	fmt.Fprintf(buffer, "\tFunctionIndex: %v\n", result.FunctionIndex)
	fmt.Fprintf(buffer, "\tValueName: %v\n", result.ValueName)
	fmt.Fprintf(buffer, "\tHasValue: %v\n", result.Value != (reflect.Value{}))
	fmt.Fprintf(buffer, "\tHookIndexes: %v\n", result.HookIndexes)
}

func (p *Program) dumpHook(i int, hook *hook, buffer *bytes.Buffer) {
	fmt.Fprintf(buffer, "Hook[%d]:\n", i)
	fmt.Fprintf(buffer, "\tFunctionIndex: %v\n", hook.FunctionIndex)
	fmt.Fprintf(buffer, "\tValueRef: %v\n", hook.ValueRef)
	fmt.Fprintf(buffer, "\tHasValueReceiver: %v\n", hook.ValueReceiver != (reflect.Value{}))
	fmt.Fprintf(buffer, "\tHasCallback: %v\n", hook.Callback != nil)
	fmt.Fprintf(buffer, "\tReceiveValueAddr: %v\n", hook.ReceiveValueAddr)
}

func (p *Program) DumpAsString() string {
	var buffer bytes.Buffer
	p.Dump(&buffer)
	return buffer.String()
}
