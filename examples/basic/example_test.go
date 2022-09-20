package di_test

import (
	"context"
	"fmt"

	"github.com/go-tk/di"
)

func Example() {
	var program di.Program

	provideBaz(&program)
	provideFoo(&program)
	showAll(&program)
	provideBar(&program)
	// NOTE: Program will rearrange Functions properly basing on dependency analysis.

	defer program.Clean()
	program.MustRun(context.Background())
	// Output:
	// foo = 100
	// bar = 200
	// baz = 300
	// foo, bar, baz = 100, 200, 300
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

func provideBar(program *di.Program) {
	var (
		foo int
		bar int
	)
	program.MustNewFunction(
		di.Argument("FOO", &foo),
		di.Result("BAR", &bar),
		di.Body(func(context.Context) error {
			bar = foo * 2
			fmt.Printf("bar = %d\n", bar)
			return nil
		}),
	)
}

func provideBaz(program *di.Program) {
	var (
		foo int
		bar int
		baz int
	)
	program.MustNewFunction(
		di.Argument("FOO", &foo),
		di.Argument("BAR", &bar),
		di.Result("BAZ", &baz),
		di.Body(func(context.Context) error {
			baz = foo + bar
			fmt.Printf("baz = %d\n", baz)
			return nil
		}),
	)
}

func showAll(program *di.Program) {
	var (
		foo int
		bar int
		baz int
	)
	program.MustNewFunction(
		di.Argument("FOO", &foo),
		di.Argument("BAR", &bar),
		di.Argument("BAZ", &baz),
		di.Body(func(context.Context) error {
			fmt.Printf("foo, bar, baz = %d, %d, %d\n", foo, bar, baz)
			return nil
		}),
	)
}
