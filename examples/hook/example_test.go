package di_test

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-tk/di"
)

func Example() {
	var program di.Program

	showPetNameList(&program)
	modifyPetNameList(&program)
	providePetNameList(&program)
	provideAdditionalPetName(&program)
	// NOTE: Program will rearrange Functions properly basing on dependency analysis.

	defer program.Clean()
	program.MustRun(context.Background())
	// Output:
	// user name list: tom,jeff,spike
}

func providePetNameList(program *di.Program) {
	var userNameList []string
	program.MustNewFunction(
		di.Result("PET_NAME_LIST", &userNameList),
		di.Body(func(context.Context) error {
			userNameList = []string{"tom", "jeff"}
			return nil
		}),
	)
}

func provideAdditionalPetName(program *di.Program) {
	var additionalPetName string
	program.MustNewFunction(
		di.Result("ADDITIONAL_PET_NAME", &additionalPetName),
		di.Body(func(context.Context) error {
			additionalPetName = "spike"
			return nil
		}),
	)
}

func showPetNameList(program *di.Program) {
	var userNameList []string
	program.MustNewFunction(
		di.Argument("PET_NAME_LIST", &userNameList),
		di.Body(func(context.Context) error {
			fmt.Printf("user name list: %v\n", strings.Join(userNameList, ","))
			return nil
		}),
	)
}

func modifyPetNameList(program *di.Program) {
	var (
		additionalPetName string
		userNameList      *[]string
	)
	program.MustNewFunction(
		di.Argument("ADDITIONAL_PET_NAME", &additionalPetName),
		di.Body(func(context.Context) error { return nil }),
		di.Hook("PET_NAME_LIST", &userNameList, func(context.Context) error {
			*userNameList = append(*userNameList, additionalPetName)
			return nil
		}),
	)
}
