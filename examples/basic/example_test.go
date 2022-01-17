package di_test

import (
	"context"
	"fmt"

	"github.com/go-tk/di"
)

func Example() {
	var p di.Program
	defer p.Clean()
	p.MustAddFunctions(Bar(), Qux(), Baz(), Foo())
	// NOTE: Program will rearrange above Functions properly basing on dependency analysis.
	p.MustRun(context.Background())
	// Output:
	// x = 100
	// y = 200
	// z = 300
	// x, y, z = 100, 200, 300
}

func Foo() di.Function {
	var x int
	return di.Function{
		Tag: di.FullFunctionName(Foo),
		Results: []di.Result{
			{OutValueID: "x", OutValuePtr: &x},
		},
		Body: func(_ context.Context) error {
			x = 100
			fmt.Printf("x = %d\n", x)
			return nil
		},
	}
}

func Bar() di.Function {
	var x int
	var y int
	return di.Function{
		Tag: di.FullFunctionName(Bar),
		Arguments: []di.Argument{
			{InValueID: "x", InValuePtr: &x},
		},
		Results: []di.Result{
			{OutValueID: "y", OutValuePtr: &y},
		},
		Body: func(_ context.Context) error {
			y = 2 * x
			fmt.Printf("y = %d\n", y)
			return nil
		},
	}
}

func Baz() di.Function {
	var x int
	var y int
	var z int
	return di.Function{
		Tag: di.FullFunctionName(Baz),
		Arguments: []di.Argument{
			{InValueID: "x", InValuePtr: &x},
			{InValueID: "y", InValuePtr: &y},
		},
		Results: []di.Result{
			{OutValueID: "z", OutValuePtr: &z},
		},
		Body: func(_ context.Context) error {
			z = x + y
			fmt.Printf("z = %d\n", z)
			return nil
		},
	}
}

func Qux() di.Function {
	var x int
	var y int
	var z int
	return di.Function{
		Tag: di.FullFunctionName(Qux),
		Arguments: []di.Argument{
			{InValueID: "x", InValuePtr: &x},
			{InValueID: "y", InValuePtr: &y},
			{InValueID: "z", InValuePtr: &z},
		},
		Body: func(_ context.Context) error {
			fmt.Printf("x, y, z = %d, %d, %d\n", x, y, z)
			return nil
		},
	}
}
