package di_test

import (
	"context"
	"fmt"

	. "github.com/go-tk/di"
)

func ExampleProgram() {
	var p Program
	p.MustAddFunction(baz())
	p.MustAddFunction(bar())
	p.MustAddFunction(foo())
	p.MustRun(context.Background())
	// Output:
	// x = 100
	// y = 200
	// z = 300
}

func foo() Function {
	var x int
	return Function{
		Tag: "foo",
		Results: []Result{
			{ValueID: "x", ValuePtr: &x},
		},
		Body: func(_ context.Context) error {
			x = 100
			fmt.Printf("x = %d\n", x)
			return nil
		},
	}
}

func bar() Function {
	var x int
	var y int
	return Function{
		Tag: "bar",
		Arguments: []Argument{
			{ValueID: "x", ValuePtr: &x},
		},
		Results: []Result{
			{ValueID: "y", ValuePtr: &y},
		},
		Body: func(_ context.Context) error {
			y = 2 * x
			fmt.Printf("y = %d\n", y)
			return nil
		},
	}
}

func baz() Function {
	var x int
	var y int
	var z int
	return Function{
		Tag: "baz",
		Arguments: []Argument{
			{ValueID: "x", ValuePtr: &x},
			{ValueID: "y", ValuePtr: &y},
		},
		Results: []Result{
			{ValueID: "z", ValuePtr: &z},
		},
		Body: func(_ context.Context) error {
			z = x + y
			fmt.Printf("z = %d\n", z)
			return nil
		},
	}
}
