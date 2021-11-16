package di

import (
	"reflect"
	"runtime"
	"strings"
)

// FullFunctionName returns the full name of the given function.
func FullFunctionName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// PackagePath returns the package path of the caller.
func PackagePath(skip int) string {
	pc, _, _, _ := runtime.Caller(skip + 1)
	fullFunctionName := runtime.FuncForPC(pc).Name()
	i := strings.LastIndexByte(fullFunctionName, '.')
	packagePath := fullFunctionName[:i]
	return packagePath
}

// PackageName returns the package name of the caller.
func PackageName(skip int) string {
	packagePath := PackagePath(skip + 1)
	i := strings.LastIndexByte(packagePath, '/')
	packageName := packagePath[i+1:]
	return packageName
}
