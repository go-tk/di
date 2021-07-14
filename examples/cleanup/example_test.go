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
	// NOTE: The order that adds these Functions into the Program is insignificant, the
	//       Program will rearrange these Functions properly basing on dependency analysis.
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
			{OutValueID: "temp-dir-name", OutValuePtr: &tempDirName, CleanupPtr: &cleanup},
		},
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
		Tag: "bar",
		Arguments: []di.Argument{
			{InValueID: "temp-dir-name", InValuePtr: &tempDirName},
		},
		Results: []di.Result{
			{OutValueID: "temp-file", OutValuePtr: &tempFile, CleanupPtr: &cleanup},
		},
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
		Tag: "baz",
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
