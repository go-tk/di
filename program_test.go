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
			Given("argument with empty in-value id").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				var a1 int
				w.Input.F.Arguments = []Argument{{InValuePtr: &a1}}
				w.ExpectedOutput.Err = ErrInvalidArgument
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": empty in-value id; tag=\"foo\""
			}),
		tc.Copy().
			Given("argument without in-value pointer").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				w.Input.F.Arguments = []Argument{{InValueID: "a1"}}
				w.ExpectedOutput.Err = ErrInvalidArgument
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": no in-value pointer; tag=\"foo\" inValueID=\"a1\""
			}),
		tc.Copy().
			Given("argument with invalid in-value pointer type").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				var a1 string
				w.Input.F.Arguments = []Argument{{InValueID: "a1", InValuePtr: a1}}
				w.ExpectedOutput.Err = ErrInvalidArgument
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": invalid in-value pointer type; tag=\"foo\" inValueID=\"a1\" inValuePtrType=string"
			}),
		tc.Copy().
			Given("argument with nil in-value pointer").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				w.Input.F.Arguments = []Argument{{InValueID: "a1", InValuePtr: (*int)(nil)}}
				w.ExpectedOutput.Err = ErrInvalidArgument
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": nil in-value pointer; tag=\"foo\" inValueID=\"a1\""
			}),
		tc.Copy().
			Given("result with empty out-value id").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				var a1 int
				w.Input.F.Results = []Result{{OutValuePtr: &a1}}
				w.ExpectedOutput.Err = ErrInvalidResult
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": empty out-value id; tag=\"foo\""
			}),
		tc.Copy().
			Given("result without out-value pointer").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				w.Input.F.Results = []Result{{OutValueID: "r1"}}
				w.ExpectedOutput.Err = ErrInvalidResult
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": no out-value pointer; tag=\"foo\" outValueID=\"r1\""
			}),
		tc.Copy().
			Given("result with invalid out-value pointer type").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				var r1 string
				w.Input.F.Results = []Result{{OutValueID: "r1", OutValuePtr: r1}}
				w.ExpectedOutput.Err = ErrInvalidResult
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": invalid out-value pointer type; tag=\"foo\" outValueID=\"r1\" outValuePtrType=string"
			}),
		tc.Copy().
			Given("result with nil out-value pointer").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				w.Input.F.Results = []Result{{OutValueID: "r1", OutValuePtr: (*int)(nil)}}
				w.ExpectedOutput.Err = ErrInvalidResult
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": nil out-value pointer; tag=\"foo\" outValueID=\"r1\""
			}),
		tc.Copy().
			Given("hook with empty in-value id").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				var a1 int
				w.Input.F.Hooks = []Hook{{InValuePtr: &a1}}
				w.ExpectedOutput.Err = ErrInvalidHook
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": empty in-value id; tag=\"foo\""
			}),
		tc.Copy().
			Given("hook without in-value pointer").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				w.Input.F.Hooks = []Hook{{InValueID: "r1"}}
				w.ExpectedOutput.Err = ErrInvalidHook
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": no in-value pointer; tag=\"foo\" inValueID=\"r1\""
			}),
		tc.Copy().
			Given("hook with invalid in-value pointer type").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				var r1 string
				w.Input.F.Hooks = []Hook{{InValueID: "r1", InValuePtr: r1}}
				w.ExpectedOutput.Err = ErrInvalidHook
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": invalid in-value pointer type; tag=\"foo\" inValueID=\"r1\" inValuePtrType=string"
			}),
		tc.Copy().
			Given("hook with nil in-value pointer").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				w.Input.F.Hooks = []Hook{{InValueID: "r1", InValuePtr: (*int)(nil)}}
				w.ExpectedOutput.Err = ErrInvalidHook
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": nil in-value pointer; tag=\"foo\" inValueID=\"r1\""
			}),
		tc.Copy().
			Given("hook with nil callback pointer").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				var r1 *int
				w.Input.F.Hooks = []Hook{{InValueID: "r1", InValuePtr: &r1}}
				w.ExpectedOutput.Err = ErrInvalidHook
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + ": nil callback pointer; tag=\"foo\" inValueID=\"r1\""
			}),
		tc.Copy().
			Then("should succeed").
			AddTask(1999, func(w *Workspace) {
				w.Input.F.Tag = "foo"
				w.Input.F.Body = func(context.Context) error { return nil }
				var a1 *int
				w.Input.F.Arguments = []Argument{{InValueID: "a1", InValuePtr: &a1}}
				var r1 *int
				w.Input.F.Results = []Result{{OutValueID: "r1", OutValuePtr: &r1}}
				cb := func(context.Context) error { return nil }
				w.Input.F.Hooks = []Hook{{InValueID: "r1", InValuePtr: &r1, CallbackPtr: &cb}}
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
			When("results with same out-value ids").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				var var1 int
				var var2 int
				w.P.MustAddFunction(Function{
					Tag: "foo",
					Results: []Result{
						{OutValueID: "var", OutValuePtr: &var1},
					},
					Body: func(context.Context) error { return nil },
				})
				w.P.MustAddFunction(Function{
					Tag: "bar",
					Results: []Result{
						{OutValueID: "var", OutValuePtr: &var2},
					},
					Body: func(context.Context) error { return nil },
				})
				w.ExpectedOutput.Err = ErrValueAlreadyExists
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + "; tag1=\"bar\" tag2=\"foo\" outValueID=\"var\""
			}),
		tc.Copy().
			When("in-value of argument not provisioned").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				{
					var x int
					w.P.MustAddFunction(Function{
						Tag: "foo",
						Results: []Result{
							{OutValueID: "x", OutValuePtr: &x},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				{
					var x, y int
					w.P.MustAddFunction(Function{
						Tag: "bar",
						Arguments: []Argument{
							{InValueID: "x", InValuePtr: &x},
							{InValueID: "y", InValuePtr: &y},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				w.ExpectedOutput.Err = ErrValueNotFound
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + "; tag=\"bar\" inValueID=\"y\""
			}),
		tc.Copy().
			When("in-value of optional argument not provisioned").
			Then("should not fail").
			AddTask(1999, func(w *Workspace) {
				{
					var x int
					w.P.MustAddFunction(Function{
						Tag: "foo",
						Results: []Result{
							{OutValueID: "x", OutValuePtr: &x},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				{
					var x, y int
					w.P.MustAddFunction(Function{
						Tag: "bar",
						Arguments: []Argument{
							{InValueID: "x", InValuePtr: &x},
							{InValueID: "y", InValuePtr: &y, IsOptional: true},
						},
						Body: func(context.Context) error { return nil },
					})
				}
			}),
		tc.Copy().
			When("in-value type of argument and out-value type of result not matched").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				{
					var x int
					w.P.MustAddFunction(Function{
						Tag: "foo",
						Results: []Result{
							{OutValueID: "x", OutValuePtr: &x},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				{
					var x string
					w.P.MustAddFunction(Function{
						Tag: "bar",
						Arguments: []Argument{
							{InValueID: "x", InValuePtr: &x},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				w.ExpectedOutput.Err = ErrValueTypeMismatch
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + "; tag1=\"bar\" tag2=\"foo\" valueID=\"x\" inValueType=string outValueType=int"
			}),
		tc.Copy().
			When("in-value of hook not provisioned").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				{
					var x int
					w.P.MustAddFunction(Function{
						Tag: "foo",
						Results: []Result{
							{OutValueID: "x", OutValuePtr: &x},
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
							{InValueID: "x", InValuePtr: &x, CallbackPtr: &cb},
							{InValueID: "y", InValuePtr: &y, CallbackPtr: &cb},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				w.ExpectedOutput.Err = ErrValueNotFound
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + "; tag=\"bar\" inValueID=\"y\""
			}),
		tc.Copy().
			When("in-value type of hook and out-value type of result not matched").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				{
					var x int
					w.P.MustAddFunction(Function{
						Tag: "foo",
						Results: []Result{
							{OutValueID: "x", OutValuePtr: &x},
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
							{InValueID: "x", InValuePtr: &x, CallbackPtr: &cb},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				w.ExpectedOutput.Err = ErrValueTypeMismatch
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + "; tag1=\"bar\" tag2=\"foo\" valueID=\"x\" inValueType=string outValueType=int"
			}),
		tc.Copy().
			When("circular dependencies exist (1)").
			Then("should fail").
			AddTask(1999, func(w *Workspace) {
				var x int
				w.P.MustAddFunction(Function{
					Tag: "foo",
					Arguments: []Argument{
						{InValueID: "x", InValuePtr: &x},
					},
					Results: []Result{
						{OutValueID: "x", OutValuePtr: &x},
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
						{OutValueID: "x", OutValuePtr: &x},
					},
					Hooks: []Hook{
						{InValueID: "x", InValuePtr: &x, CallbackPtr: &cb},
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
							{InValueID: "x", InValuePtr: &x},
						},
						Results: []Result{
							{OutValueID: "y", OutValuePtr: &y},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				{
					var x, y int
					w.P.MustAddFunction(Function{
						Tag: "bar",
						Arguments: []Argument{
							{InValueID: "y", InValuePtr: &y},
						},
						Results: []Result{
							{OutValueID: "x", OutValuePtr: &x},
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
							{OutValueID: "x", OutValuePtr: &x},
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
							{InValueID: "x", InValuePtr: &x},
						},
						Hooks: []Hook{
							{InValueID: "x", InValuePtr: &x, CallbackPtr: &cb},
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
						{OutValueID: "x", OutValuePtr: &x, CleanupPtr: &c},
					},
					Body: func(context.Context) error { return nil },
				})
				w.ExpectedOutput.Err = ErrNilCleanup
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + "; tag=\"foo\" outValueID=\"x\""
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
							{OutValueID: "x", OutValuePtr: &x},
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
							{InValueID: "x", InValuePtr: &x, CallbackPtr: &cb},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				w.ExpectedOutput.Err = ErrNilCallback
				w.ExpectedOutput.ErrStr = w.ExpectedOutput.Err.Error() + "; tag=\"bar\" inValueID=\"x\""
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
							{OutValueID: "x", OutValuePtr: &x},
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
							{InValueID: "x", InValuePtr: &x, CallbackPtr: &cb},
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
				w.ExpectedOutput.ErrStr = "di: callback failed; tag=\"bar\" inValueID=\"x\": " + w.ExpectedOutput.Err.Error()
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
							{InValueID: "x", InValuePtr: &x},
							{InValueID: "y", InValuePtr: &y},
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
							{InValueID: "x", InValuePtr: &x},
						},
						Results: []Result{
							{OutValueID: "y", OutValuePtr: &y, CleanupPtr: &c},
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
							{OutValueID: "x", OutValuePtr: &x, CleanupPtr: &c},
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
							{InValueID: "y", InValuePtr: &y, CallbackPtr: &cb},
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
