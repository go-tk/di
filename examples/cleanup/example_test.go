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
	var p di.Program
	p.MustAddFunctions(
		Baz(),
		Foo(),
		Bar(),
		// NOTE: Program will rearrange above Functions properly basing on dependency analysis.
	)
	defer p.Clean()
	p.MustRun(context.Background())
	// Output:
	// create temp dir
	// create and open temp file
	// write temp file
	// close and delete temp file
	// delete temp dir
}

func Foo() di.Function {
	var tempDirName string
	var cleanup func()
	return di.Function{
		Tag: di.FullFunctionName(Foo),
		Results: []di.Result{
			{OutValueID: "temp-dir-name", OutValuePtr: &tempDirName},
		},
		CleanupPtr: &cleanup,
		Body: func(_ context.Context) error {
			fmt.Println("create temp dir")
			var err error
			tempDirName, err = ioutil.TempDir("", "")
			if err != nil {
				return err
			}
			cleanup = func() {
				fmt.Println("delete temp dir")
				os.Remove(tempDirName)
			}
			return nil
		},
	}
}

func Bar() di.Function {
	var tempDirName string
	var tempFile *os.File
	var cleanup func()
	return di.Function{
		Tag: di.FullFunctionName(Bar),
		Arguments: []di.Argument{
			{InValueID: "temp-dir-name", InValuePtr: &tempDirName},
		},
		Results: []di.Result{
			{OutValueID: "temp-file", OutValuePtr: &tempFile},
		},
		CleanupPtr: &cleanup,
		Body: func(_ context.Context) error {
			fmt.Println("create and open temp file")
			var err error
			tempFile, err = os.Create(filepath.Join(tempDirName, "temp"))
			if err != nil {
				return err
			}
			cleanup = func() {
				fmt.Println("close and delete temp file")
				tempFile.Close()
				os.Remove(tempFile.Name())
			}
			return nil
		},
	}
}

func Baz() di.Function {
	var tempFile *os.File
	return di.Function{
		Tag: di.FullFunctionName(Baz),
		Arguments: []di.Argument{
			{InValueID: "temp-file", InValuePtr: &tempFile},
		},
		Body: func(_ context.Context) error {
			fmt.Println("write temp file")
			_, err := tempFile.WriteString("hello world")
			return err
		},
	}
}
