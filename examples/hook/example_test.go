package di_test

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-tk/di"
)

func Example() {
	var program di.Program

	showUserNameList(&program)
	modifyUserNameList(&program)
	provideUserNameList(&program)
	provideAdditionalUserName(&program)
	// NOTE: Program will rearrange Functions properly basing on dependency analysis.

	defer program.Clean()
	program.MustRun(context.Background())
	// Output:
	// user name list: tom,jeff,spike
}

func provideUserNameList(program *di.Program) {
	var userNameList []string
	program.MustNewFunction(
		di.Result("USER_NAME_LIST", &userNameList),
		di.Body(func(context.Context) error {
			userNameList = []string{"tom", "jeff"}
			return nil
		}),
	)
}

func provideAdditionalUserName(program *di.Program) {
	var additionalUserName string
	program.MustNewFunction(
		di.Result("ADDITIONAL_USER_NAME", &additionalUserName),
		di.Body(func(context.Context) error {
			additionalUserName = "spike"
			return nil
		}),
	)
}

func showUserNameList(program *di.Program) {
	var userNameList []string
	program.MustNewFunction(
		di.Argument("USER_NAME_LIST", &userNameList),
		di.Body(func(context.Context) error {
			fmt.Printf("user name list: %v\n", strings.Join(userNameList, ","))
			return nil
		}),
	)
}

func modifyUserNameList(program *di.Program) {
	var (
		additionalUserName string
		userNameList       *[]string
	)
	program.MustNewFunction(
		di.Argument("ADDITIONAL_USER_NAME", &additionalUserName),
		di.Body(func(context.Context) error { return nil }),
		di.Hook("USER_NAME_LIST", &userNameList, func(context.Context) error {
			*userNameList = append(*userNameList, additionalUserName)
			return nil
		}),
	)
}
