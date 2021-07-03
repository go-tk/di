package di_test

import (
	"context"
	"fmt"

	"github.com/go-tk/di"
)

func Example() {
	var p di.Program
	defer p.Clean()
	p.MustAddFunction(Bar())
	p.MustAddFunction(Qux())
	p.MustAddFunction(Baz())
	p.MustAddFunction(Foo())
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
		Tag: "foo",
		Results: []di.Result{
			{ValueID: "x", ValuePtr: &x},
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
		Tag: "bar",
		Arguments: []di.Argument{
			{ValueID: "x", ValuePtr: &x},
		},
		Results: []di.Result{
			{ValueID: "y", ValuePtr: &y},
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
		Tag: "baz",
		Arguments: []di.Argument{
			{ValueID: "x", ValuePtr: &x},
			{ValueID: "y", ValuePtr: &y},
		},
		Results: []di.Result{
			{ValueID: "z", ValuePtr: &z},
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
		Tag: "baz",
		Arguments: []di.Argument{
			{ValueID: "x", ValuePtr: &x},
			{ValueID: "y", ValuePtr: &y},
			{ValueID: "z", ValuePtr: &z},
		},
		Body: func(_ context.Context) error {
			fmt.Printf("x, y, z = %d, %d, %d\n", x, y, z)
			return nil
		},
	}
}
