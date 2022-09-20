package di_test

import (
	"context"
	"fmt"

	"github.com/go-tk/di"
)

func Example() {
	var program di.Program

	substractFooWithBar(&program)
	provideFoo(&program)
	// NOTE: Program will rearrange Functions properly basing on dependency analysis.

	defer program.Clean()
	program.MustRun(context.Background())
	// Output:
	// foo = 100
	// foo - bar = 99
}

func provideFoo(program *di.Program) {
	var foo int
	program.MustNewFunction(
		di.Result("FOO", &foo),
		di.Body(func(context.Context) error {
			foo = 100
			fmt.Printf("foo = %d\n", foo)
			return nil
		}),
	)
}

func substractFooWithBar(program *di.Program) {
	var (
		foo int
		bar int = 1
	)
	program.MustNewFunction(
		di.Argument("FOO", &foo),
		di.OptionalArgument("BAR", &bar),
		di.Body(func(context.Context) error {
			fmt.Printf("foo - bar = %d\n", foo-bar)
			return nil
		}),
	)
}
