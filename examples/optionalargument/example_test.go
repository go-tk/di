package di_test

import (
	"context"
	"fmt"

	"github.com/go-tk/di"
)

func Example() {
	var p di.Program
	p.MustAddFunctions(
		Bar(),
		Foo(),
		// NOTE: Program will rearrange above Functions properly basing on dependency analysis.
	)
	defer p.Clean()
	p.MustRun(context.Background())
	// Output:
	// y - x = 99
}

func Foo() di.Function {
	var y int
	return di.Function{
		Tag: di.FullFunctionName(Foo),
		Results: []di.Result{
			{OutValueID: "y", OutValuePtr: &y},
		},
		Body: func(_ context.Context) error {
			y = 199
			return nil
		},
	}
}

func Bar() di.Function {
	x := 100
	var y int
	return di.Function{
		Tag: di.FullFunctionName(Bar),
		Arguments: []di.Argument{
			{InValueID: "x", InValuePtr: &x, IsOptional: true},
			{InValueID: "y", InValuePtr: &y},
		},
		Body: func(_ context.Context) error {
			fmt.Printf("y - x = %d\n", y-x)
			return nil
		},
	}
}
