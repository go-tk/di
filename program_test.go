package di_test

import (
	"context"
	"errors"
	"testing"

	. "github.com/go-tk/di"
	"github.com/go-tk/testcase"
	"github.com/stretchr/testify/assert"
)

func TestProgram_AddFunction(t *testing.T) {
	type Input struct {
		F Function
	}
	type Output struct {
		ErrStr string
		Err    error
	}
	type Workspace struct {
		testcase.WorkspaceBase

		P              Program
		Input          Input
		ExpectedOutput Output
	}
	tc := testcase.New().
		AddTask(1000, func(w *Workspace) {}).
		AddTask(2000, func(w *Workspace) {
			err := w.P.AddFunction(w.Input.F)
			var output Output
			if err != nil {
				output.ErrStr = err.Error()
				for err2 := errors.Unwrap(err); err2 != nil; err, err2 = err2, errors.Unwrap(err2) {
				}
				output.Err = err
			}
			assert.Equal(w.T(), w.ExpectedOutput, output)
		})
	testcase.RunListParallel(t,
		tc.Copy().
			Given("function with empty tag").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Body = func(context.Context) error { return nil }
				w.ExpectedOutput.Err = ErrInvalidFunction
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": empty tag"
			}),
		tc.Copy().
			Given("function with nil body").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.ExpectedOutput.Err = ErrInvalidFunction
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": nil body; tag=\"foo\""
			}),
		tc.Copy().
			Given("argument with empty value id").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				var a1 int
				w.Input.F.Arguments = []Argument{{ValuePtr: &a1}}
				w.ExpectedOutput.Err = ErrInvalidArgument
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": empty value id; tag=\"foo\""
			}),
		tc.Copy().
			Given("argument without value pointer").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				w.Input.F.Arguments = []Argument{{ValueID: "a1"}}
				w.ExpectedOutput.Err = ErrInvalidArgument
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": no value pointer; tag=\"foo\" valueID=\"a1\""
			}),
		tc.Copy().
			Given("argument with invalid value pointer type").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				var a1 string
				w.Input.F.Arguments = []Argument{{ValueID: "a1", ValuePtr: a1}}
				w.ExpectedOutput.Err = ErrInvalidArgument
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": invalid value pointer type; tag=\"foo\" valueID=\"a1\" valuePtrType=string"
			}),
		tc.Copy().
			Given("argument with nil value pointer").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				w.Input.F.Arguments = []Argument{{ValueID: "a1", ValuePtr: (*int)(nil)}}
				w.ExpectedOutput.Err = ErrInvalidArgument
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": nil value pointer; tag=\"foo\" valueID=\"a1\""
			}),
		tc.Copy().
			Given("result with empty value id").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				var a1 int
				w.Input.F.Results = []Result{{ValuePtr: &a1}}
				w.ExpectedOutput.Err = ErrInvalidResult
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": empty value id; tag=\"foo\""
			}),
		tc.Copy().
			Given("result without value pointer").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				w.Input.F.Results = []Result{{ValueID: "r1"}}
				w.ExpectedOutput.Err = ErrInvalidResult
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": no value pointer; tag=\"foo\" valueID=\"r1\""
			}),
		tc.Copy().
			Given("result with invalid value pointer type").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				var r1 string
				w.Input.F.Results = []Result{{ValueID: "r1", ValuePtr: r1}}
				w.ExpectedOutput.Err = ErrInvalidResult
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": invalid value pointer type; tag=\"foo\" valueID=\"r1\" valuePtrType=string"
			}),
		tc.Copy().
			Given("result with nil value pointer").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				w.Input.F.Results = []Result{{ValueID: "r1", ValuePtr: (*int)(nil)}}
				w.ExpectedOutput.Err = ErrInvalidResult
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": nil value pointer; tag=\"foo\" valueID=\"r1\""
			}),
		tc.Copy().
			Given("hook with empty value id").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				var a1 int
				w.Input.F.Hooks = []Hook{{ValuePtr: &a1}}
				w.ExpectedOutput.Err = ErrInvalidHook
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": empty value id; tag=\"foo\""
			}),
		tc.Copy().
			Given("hook without value pointer").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				w.Input.F.Hooks = []Hook{{ValueID: "r1"}}
				w.ExpectedOutput.Err = ErrInvalidHook
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": no value pointer; tag=\"foo\" valueID=\"r1\""
			}),
		tc.Copy().
			Given("hook with invalid value pointer type").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				var r1 string
				w.Input.F.Hooks = []Hook{{ValueID: "r1", ValuePtr: r1}}
				w.ExpectedOutput.Err = ErrInvalidHook
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": invalid value pointer type; tag=\"foo\" valueID=\"r1\" valuePtrType=string"
			}),
		tc.Copy().
			Given("hook with nil value pointer").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				w.Input.F.Hooks = []Hook{{ValueID: "r1", ValuePtr: (*int)(nil)}}
				w.ExpectedOutput.Err = ErrInvalidHook
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": nil value pointer; tag=\"foo\" valueID=\"r1\""
			}),
		tc.Copy().
			Given("hook with nil callback pointer").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				var r1 *int
				w.Input.F.Hooks = []Hook{{ValueID: "r1", ValuePtr: &r1}}
				w.ExpectedOutput.Err = ErrInvalidHook
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": nil callback pointer; tag=\"foo\" valueID=\"r1\""
			}),
		tc.Copy().
			Then("should succeed").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				var a1 *int
				w.Input.F.Arguments = []Argument{{ValueID: "a1", ValuePtr: &a1}}
				var r1 *int
				w.Input.F.Results = []Result{{ValueID: "r1", ValuePtr: &r1}}
				cb := func(context.Context) error { return nil }
				w.Input.F.Hooks = []Hook{{ValueID: "r1", ValuePtr: &r1, CallbackPtr: &cb}}
			}),
	)
}

func TestProgram_MustAddFunction(t *testing.T) {
	assert.Panics(t, func() {
		var p Program
		p.MustAddFunction(Function{})
	})
}

func TestProgram_Run(t *testing.T) {
	type Input struct {
		Ctx context.Context
	}
	type Output struct {
		ErrStr string
		Err    error
	}
	type Workspace struct {
		testcase.WorkspaceBase

		P              Program
		Input          Input
		ExpectedOutput Output
	}
	tc := testcase.New().
		AddTask(1000, func(w *Workspace) {
			w.Input.Ctx = context.Background()
		}).
		AddTask(2000, func(w *Workspace) {
			err := w.P.Run(w.Input.Ctx)
			var output Output
			if err == nil {
				w.AddCleanup(w.P.Clean)
			} else {
				output.ErrStr = err.Error()
				for err2 := errors.Unwrap(err); err2 != nil; err, err2 = err2, errors.Unwrap(err2) {
				}
				output.Err = err
			}
			assert.Equal(w.T(), w.ExpectedOutput, output)
		})
	testcase.RunListParallel(t,
		tc.Copy().
			When("results with same value ids").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				var var1 int
				var var2 int
				w.P.MustAddFunction(Function{
					Tag: "foo",
					Results: []Result{
						{ValueID: "var", ValuePtr: &var1},
					},
					Body: func(context.Context) error { return nil },
				})
				w.P.MustAddFunction(Function{
					Tag: "bar",
					Results: []Result{
						{ValueID: "var", ValuePtr: &var2},
					},
					Body: func(context.Context) error { return nil },
				})
				w.ExpectedOutput.Err = ErrValueAlreadyExists
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + "; tag1=\"bar\" tag2=\"foo\" valueID=\"var\""
			}),
		tc.Copy().
			When("value of argument not provisioned").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				{
					var x int
					w.P.MustAddFunction(Function{
						Tag: "foo",
						Results: []Result{
							{ValueID: "x", ValuePtr: &x},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				{
					var x, y int
					w.P.MustAddFunction(Function{
						Tag: "bar",
						Arguments: []Argument{
							{ValueID: "x", ValuePtr: &x},
							{ValueID: "y", ValuePtr: &y},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				w.ExpectedOutput.Err = ErrValueNotFound
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + "; tag=\"bar\" valueID=\"y\""
			}),
		tc.Copy().
			When("value of optional argument not provisioned").
			Then("should not fail").
			AddTask(1999, func(w *Workspace) {
				{
					var x int
					w.P.MustAddFunction(Function{
						Tag: "foo",
						Results: []Result{
							{ValueID: "x", ValuePtr: &x},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				{
					var x, y int
					w.P.MustAddFunction(Function{
						Tag: "bar",
						Arguments: []Argument{
							{ValueID: "x", ValuePtr: &x},
							{ValueID: "y", ValuePtr: &y, IsOptional: true},
						},
						Body: func(context.Context) error { return nil },
					})
				}
			}),
		tc.Copy().
			When("value types of argument and result not matched").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				{
					var x int
					w.P.MustAddFunction(Function{
						Tag: "foo",
						Results: []Result{
							{ValueID: "x", ValuePtr: &x},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				{
					var x string
					w.P.MustAddFunction(Function{
						Tag: "bar",
						Arguments: []Argument{
							{ValueID: "x", ValuePtr: &x},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				w.ExpectedOutput.Err = ErrValueTypeMismatch
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + "; tag1=\"bar\" tag2=\"foo\" valueID=\"x\" valueType1=string valueType2=int"
			}),
		tc.Copy().
			When("value of hook not provisioned").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				{
					var x int
					w.P.MustAddFunction(Function{
						Tag: "foo",
						Results: []Result{
							{ValueID: "x", ValuePtr: &x},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				{
					var x, y int
					cb := func(context.Context) error { return nil }
					w.P.MustAddFunction(Function{
						Tag: "bar",
						Hooks: []Hook{
							{ValueID: "x", ValuePtr: &x, CallbackPtr: &cb},
							{ValueID: "y", ValuePtr: &y, CallbackPtr: &cb},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				w.ExpectedOutput.Err = ErrValueNotFound
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + "; tag=\"bar\" valueID=\"y\""
			}),
		tc.Copy().
			When("value types of hook and result not matched").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				{
					var x int
					w.P.MustAddFunction(Function{
						Tag: "foo",
						Results: []Result{
							{ValueID: "x", ValuePtr: &x},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				{
					var x string
					cb := func(context.Context) error { return nil }
					w.P.MustAddFunction(Function{
						Tag: "bar",
						Hooks: []Hook{
							{ValueID: "x", ValuePtr: &x, CallbackPtr: &cb},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				w.ExpectedOutput.Err = ErrValueTypeMismatch
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + "; tag1=\"bar\" tag2=\"foo\" valueID=\"x\" valueType1=string valueType2=int"
			}),
		tc.Copy().
			When("circular dependencies exist (1)").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				var x int
				w.P.MustAddFunction(Function{
					Tag: "foo",
					Arguments: []Argument{
						{ValueID: "x", ValuePtr: &x},
					},
					Results: []Result{
						{ValueID: "x", ValuePtr: &x},
					},
					Body: func(context.Context) error { return nil },
				})
				w.ExpectedOutput.Err = ErrCircularDependencies
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() +
					"; {tag: \"foo\", argument: \"x\"} => {tag: \"foo\"}"
			}),
		tc.Copy().
			When("circular dependencies exist (2)").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				var x int
				cb := func(context.Context) error { return nil }
				w.P.MustAddFunction(Function{
					Tag: "foo",
					Results: []Result{
						{ValueID: "x", ValuePtr: &x},
					},
					Hooks: []Hook{
						{ValueID: "x", ValuePtr: &x, CallbackPtr: &cb},
					},
					Body: func(context.Context) error { return nil },
				})
				w.ExpectedOutput.Err = ErrCircularDependencies
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() +
					"; {tag: \"foo\", hook: \"x\"} => {tag: \"foo\"}"
			}),
		tc.Copy().
			When("circular dependencies exist (3)").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				{
					var x, y int
					w.P.MustAddFunction(Function{
						Tag: "foo",
						Arguments: []Argument{
							{ValueID: "x", ValuePtr: &x},
						},
						Results: []Result{
							{ValueID: "y", ValuePtr: &y},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				{
					var x, y int
					w.P.MustAddFunction(Function{
						Tag: "bar",
						Arguments: []Argument{
							{ValueID: "y", ValuePtr: &y},
						},
						Results: []Result{
							{ValueID: "x", ValuePtr: &x},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				w.ExpectedOutput.Err = ErrCircularDependencies
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() +
					"; {tag: \"foo\", argument: \"x\"} => {tag: \"bar\", argument: \"y\"} => {tag: \"foo\"}"
			}),
		tc.Copy().
			When("circular dependencies exist (4)").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				{
					var x int
					w.P.MustAddFunction(Function{
						Tag: "foo",
						Results: []Result{
							{ValueID: "x", ValuePtr: &x},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				{
					var x int
					cb := func(context.Context) error { return nil }
					w.P.MustAddFunction(Function{
						Tag: "bar",
						Arguments: []Argument{
							{ValueID: "x", ValuePtr: &x},
						},
						Hooks: []Hook{
							{ValueID: "x", ValuePtr: &x, CallbackPtr: &cb},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				w.ExpectedOutput.Err = ErrCircularDependencies
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() +
					"; {tag: \"foo\", hook: \"x\"} => {tag: \"bar\", argument: \"x\"} => {tag: \"foo\"}"
			}),
		tc.Copy().
			When("error returned by function body").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.P.MustAddFunction(Function{
					Tag:  "foo",
					Body: func(context.Context) error { return context.DeadlineExceeded },
				})
				w.ExpectedOutput.Err = context.DeadlineExceeded
				w.ExpectedOutput.ErrStr = "di: function failed; tag=\"foo\": " + w.ExpectedOutput.Err.Error()
			}),
		tc.Copy().
			When("nil cleanup").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				var x int
				var c func()
				w.P.MustAddFunction(Function{
					Tag: "foo",
					Results: []Result{
						{ValueID: "x", ValuePtr: &x, CleanupPtr: &c},
					},
					Body: func(context.Context) error { return nil },
				})
				w.ExpectedOutput.Err = ErrNilCleanup
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + "; tag=\"foo\" valueID=\"x\""
			}),
		tc.Copy().
			When("nil callback").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				{
					var x int
					w.P.MustAddFunction(Function{
						Tag: "foo",
						Results: []Result{
							{ValueID: "x", ValuePtr: &x},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				{
					var x int
					var cb func(context.Context) error
					w.P.MustAddFunction(Function{
						Tag: "bar",
						Hooks: []Hook{
							{ValueID: "x", ValuePtr: &x, CallbackPtr: &cb},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				w.ExpectedOutput.Err = ErrNilCallback
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + "; tag=\"bar\" valueID=\"x\""
			}),
		tc.Copy().
			When("callback failed").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				{
					var x int
					w.P.MustAddFunction(Function{
						Tag: "foo",
						Results: []Result{
							{ValueID: "x", ValuePtr: &x},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				{
					var x int
					var cb func(context.Context) error
					w.P.MustAddFunction(Function{
						Tag: "bar",
						Hooks: []Hook{
							{ValueID: "x", ValuePtr: &x, CallbackPtr: &cb},
						},
						Body: func(context.Context) error {
							cb = func(context.Context) error {
								return context.DeadlineExceeded
							}
							return nil
						},
					})
				}
				w.ExpectedOutput.Err = context.DeadlineExceeded
				w.ExpectedOutput.ErrStr = "di: callback failed; tag=\"bar\" valueID=\"x\": " + w.ExpectedOutput.Err.Error()
			}),
		tc.Copy().
			Then("should succeed").
			AddTask(1999, func(w *Workspace) {
				var s string
				w.AddCleanup(func() { assert.Equal(w.T(), "1234567", s) })
				{
					var x, y int
					w.P.MustAddFunction(Function{
						Tag: "foo",
						Arguments: []Argument{
							{ValueID: "x", ValuePtr: &x},
							{ValueID: "y", ValuePtr: &y},
						},
						Body: func(context.Context) error {
							s += "5"
							assert.Equal(w.T(), 101, x)
							assert.Equal(w.T(), 404, y)
							return nil
						},
					})
				}
				{
					var x, y int
					var c func()
					w.P.MustAddFunction(Function{
						Tag: "bar",
						Arguments: []Argument{
							{ValueID: "x", ValuePtr: &x},
						},
						Results: []Result{
							{ValueID: "y", ValuePtr: &y, CleanupPtr: &c},
						},
						Body: func(context.Context) error {
							s += "3"
							c = func() { s += "6" }
							assert.Equal(w.T(), 101, x)
							y = 404
							return nil
						},
					})
				}
				{
					var x int
					var c func()
					w.P.MustAddFunction(Function{
						Tag: "baz",
						Results: []Result{
							{ValueID: "x", ValuePtr: &x, CleanupPtr: &c},
						},
						Body: func(context.Context) error {
							s += "1"
							c = func() { s += "7" }
							x = 101
							return nil
						},
					})
				}
				{
					var y int
					var cb func(context.Context) error
					w.P.MustAddFunction(Function{
						Tag: "qux",
						Hooks: []Hook{
							{ValueID: "y", ValuePtr: &y, CallbackPtr: &cb},
						},
						Body: func(context.Context) error {
							s += "2"
							cb = func(context.Context) error {
								s += "4"
								assert.Equal(w.T(), 404, y)
								return nil
							}
							return nil
						},
					})
				}
			}),
	)
}

func TestProgram_MustRun(t *testing.T) {
	assert.Panics(t, func() {
		var p Program
		p.MustAddFunction(Function{
			Tag: "foo",
			Body: func(context.Context) error {
				return errors.New("")
			},
		})
		p.MustRun(context.Background())
	})
}
