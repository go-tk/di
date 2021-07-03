package di_test

import (
	"context"
	"fmt"

	"github.com/go-tk/di"
)

func Example() {
	var p di.Program
	defer p.Clean()
	p.MustAddFunction(Foo())
	p.MustAddFunction(Baz())
	p.MustAddFunction(Qux())
	p.MustAddFunction(Bar())
	p.MustRun(context.Background())
	// Output:
	// user name list: [tom jeff kim]
}

func Foo() di.Function {
	var userNameList *[]string
	return di.Function{
		Tag: "foo",
		Results: []di.Result{
			{ValueID: "user-name-list", ValuePtr: &userNameList},
		},
		Body: func(_ context.Context) error {
			userNameList = &[]string{"tom", "jeff"}
			return nil
		},
	}
}

func Bar() di.Function {
	var additionalUserName string
	return di.Function{
		Tag: "bar",
		Results: []di.Result{
			{ValueID: "additional-user-name", ValuePtr: &additionalUserName},
		},
		Body: func(_ context.Context) error {
			additionalUserName = "kim"
			return nil
		},
	}
}

func Baz() di.Function {
	var userNameList *[]string
	return di.Function{
		Tag: "baz",
		Arguments: []di.Argument{
			{ValueID: "user-name-list", ValuePtr: &userNameList},
		},
		Body: func(_ context.Context) error {
			fmt.Printf("user name list: %v\n", *userNameList)
			return nil
		},
	}
}

func Qux() di.Function {
	var additionalUserName string
	var userNameList *[]string
	var callback func(_ context.Context) error
	return di.Function{
		Tag: "qux",
		Arguments: []di.Argument{
			{ValueID: "additional-user-name", ValuePtr: &additionalUserName},
		},
		Hooks: []di.Hook{
			{ValueID: "user-name-list", ValuePtr: &userNameList, CallbackPtr: &callback},
		},
		Body: func(_ context.Context) error {
			callback = func(_ context.Context) error {
				*userNameList = append(*userNameList, additionalUserName)
				return nil
			}
			return nil
		},
	}
}
