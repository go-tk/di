package di

import (
	"reflect"
	"runtime"
)

// FullFunctionName returns the full name of the given function.
func FullFunctionName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
