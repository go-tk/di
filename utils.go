package di

import (
	"runtime"
	"strings"
)

// PackagePath returns the package path of the caller.
func PackagePath(skip int) string {
	pc, _, _, _ := runtime.Caller(skip + 1)
	functionFullName := runtime.FuncForPC(pc).Name()
	i := strings.LastIndexByte(functionFullName, '.')
	packagePath := functionFullName[:i]
	return packagePath
}

// PackageName returns the package name of the caller.
func PackageName(skip int) string {
	packagePath := PackagePath(skip + 1)
	i := strings.LastIndexByte(packagePath, '/')
	packageName := packagePath[i+1:]
	return packageName
}
