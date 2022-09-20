package di_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-tk/di"
)

func Example() {
	var program di.Program

	writeTempFile(&program)
	provideTempDirName(&program)
	provideTempFile(&program)
	// NOTE: Program will rearrange Functions properly basing on dependency analysis.

	defer program.Clean()
	program.MustRun(context.Background())
	// Output:
	// 1. create temp dir
	// 2. create and open temp file
	// 3. write temp file
	// 4. close and delete temp file
	// 5. delete temp dir
}

func provideTempDirName(program *di.Program) {
	var tempDirName string
	program.MustNewFunction(
		di.Result("TEMP_DIR_NAME", &tempDirName),
		di.Body(func(context.Context) error {
			fmt.Println("1. create temp dir")
			var err error
			tempDirName, err = ioutil.TempDir("", "")
			return err
		}),
		di.Cleanup(func() {
			fmt.Println("5. delete temp dir")
			os.Remove(tempDirName)
		}),
	)
}

func provideTempFile(program *di.Program) {
	var (
		tempDirName string
		tempFile    *os.File
	)
	program.MustNewFunction(
		di.Argument("TEMP_DIR_NAME", &tempDirName),
		di.Result("TEMP_FILE", &tempFile),
		di.Body(func(context.Context) error {
			fmt.Println("2. create and open temp file")
			var err error
			tempFile, err = os.Create(filepath.Join(tempDirName, "temp"))
			return err
		}),
		di.Cleanup(func() {
			fmt.Println("4. close and delete temp file")
			tempFile.Close()
			os.Remove(tempFile.Name())
		}),
	)
}

func writeTempFile(program *di.Program) {
	var tempFile *os.File
	program.MustNewFunction(
		di.Argument("TEMP_FILE", &tempFile),
		di.Body(func(context.Context) error {
			fmt.Println("3. write temp file")
			_, err := tempFile.WriteString("hello world")
			return err
		}),
	)
}
