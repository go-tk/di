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
	defer p.Clean()
	p.MustAddFunction(Baz())
	p.MustAddFunction(Foo())
	p.MustAddFunction(Bar())
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
		Tag: "foo",
		Results: []di.Result{
			{ValueID: "temp-dir-name", ValuePtr: &tempDirName, CleanupPtr: &cleanup},
		},
		Body: func(_ context.Context) error {
			fmt.Println("create temp dir")
			var err error
			tempDirName, err = ioutil.TempDir("", "")
			cleanup = func() {
				fmt.Println("delete temp dir")
				os.Remove(tempDirName)
			}
			return err
		},
	}
}

func Bar() di.Function {
	var tempDirName string
	var tempFile *os.File
	var cleanup func()
	return di.Function{
		Tag: "bar",
		Arguments: []di.Argument{
			{ValueID: "temp-dir-name", ValuePtr: &tempDirName},
		},
		Results: []di.Result{
			{ValueID: "temp-file", ValuePtr: &tempFile, CleanupPtr: &cleanup},
		},
		Body: func(_ context.Context) error {
			fmt.Println("create and open temp file")
			var err error
			tempFile, err = os.Create(filepath.Join(tempDirName, "temp"))
			cleanup = func() {
				fmt.Println("close and delete temp file")
				tempFile.Close()
				os.Remove(tempFile.Name())
			}
			return err
		},
	}
}

func Baz() di.Function {
	var tempFile *os.File
	return di.Function{
		Tag: "baz",
		Arguments: []di.Argument{
			{ValueID: "temp-file", ValuePtr: &tempFile},
		},
		Body: func(_ context.Context) error {
			fmt.Println("write temp file")
			_, err := tempFile.WriteString("hello world")
			return err
		},
	}
}
