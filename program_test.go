package di_test

import (
	"context"
	"errors"
	"testing"

	. "github.com/go-tk/di"
	"github.com/go-tk/testcase"
	"github.com/stretchr/testify/assert"
)

func TestProgram_AddFunctions(t *testing.T) {
	type Workspace struct {
		P  Program
		In struct {
			F Function
		}
		ExpOut, ActOut struct {
			ErrStr string
			Err    error
		}
	}
	tc := testcase.New().
		Step(1, func(t *testing.T, w *Workspace) {}).
		Step(2, func(t *testing.T, w *Workspace) {
			err := w.P.AddFunctions(w.In.F)
			if err != nil {
				w.ActOut.ErrStr = err.Error()
				for err2 := errors.Unwrap(err); err2 != nil; err, err2 = err2, errors.Unwrap(err2) {
				}
				w.ActOut.Err = err
			}
		}).
		Step(3, func(t *testing.T, w *Workspace) {
			assert.Equal(t, w.ExpOut, w.ActOut)
		})
	testcase.RunListParallel(t,
		tc.Copy().
			Given("function with empty tag").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.In.F.Body = func(context.Context) error { return nil }
				w.ExpOut.ErrStr = ErrInvalidFunction.Error() + ": empty tag"
				w.ExpOut.Err = ErrInvalidFunction
			}),
		tc.Copy().
			Given("function with nil body").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.In.F.Tag = "foo"
				w.ExpOut.ErrStr = ErrInvalidFunction.Error() + ": nil body; tag=\"foo\""
				w.ExpOut.Err = ErrInvalidFunction
			}),
		tc.Copy().
			Given("argument with empty in-value id").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.In.F.Tag = "foo"
				w.In.F.Body = func(context.Context) error { return nil }
				var a1 int
				w.In.F.Arguments = []Argument{{InValuePtr: &a1}}
				w.ExpOut.ErrStr = ErrInvalidArgument.Error() + ": empty in-value id; tag=\"foo\""
				w.ExpOut.Err = ErrInvalidFunction
			}),
		tc.Copy().
			Given("argument without in-value pointer").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.In.F.Tag = "foo"
				w.In.F.Body = func(context.Context) error { return nil }
				w.In.F.Arguments = []Argument{{InValueID: "a1"}}
				w.ExpOut.ErrStr = ErrInvalidArgument.Error() + ": no in-value pointer; tag=\"foo\" inValueID=\"a1\""
				w.ExpOut.Err = ErrInvalidFunction
			}),
		tc.Copy().
			Given("argument with invalid in-value pointer type").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.In.F.Tag = "foo"
				w.In.F.Body = func(context.Context) error { return nil }
				var a1 string
				w.In.F.Arguments = []Argument{{InValueID: "a1", InValuePtr: a1}}
				w.ExpOut.ErrStr = ErrInvalidArgument.Error() + ": invalid in-value pointer type; tag=\"foo\" inValueID=\"a1\" inValuePtrType=\"string\""
				w.ExpOut.Err = ErrInvalidFunction
			}),
		tc.Copy().
			Given("argument with nil in-value pointer").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.In.F.Tag = "foo"
				w.In.F.Body = func(context.Context) error { return nil }
				w.In.F.Arguments = []Argument{{InValueID: "a1", InValuePtr: (*int)(nil)}}
				w.ExpOut.ErrStr = ErrInvalidArgument.Error() + ": nil in-value pointer; tag=\"foo\" inValueID=\"a1\""
				w.ExpOut.Err = ErrInvalidFunction
			}),
		tc.Copy().
			Given("result with empty out-value id").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.In.F.Tag = "foo"
				w.In.F.Body = func(context.Context) error { return nil }
				var a1 int
				w.In.F.Results = []Result{{OutValuePtr: &a1}}
				w.ExpOut.ErrStr = ErrInvalidResult.Error() + ": empty out-value id; tag=\"foo\""
				w.ExpOut.Err = ErrInvalidFunction
			}),
		tc.Copy().
			Given("result without out-value pointer").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.In.F.Tag = "foo"
				w.In.F.Body = func(context.Context) error { return nil }
				w.In.F.Results = []Result{{OutValueID: "r1"}}
				w.ExpOut.ErrStr = ErrInvalidResult.Error() + ": no out-value pointer; tag=\"foo\" outValueID=\"r1\""
				w.ExpOut.Err = ErrInvalidFunction
			}),
		tc.Copy().
			Given("result with invalid out-value pointer type").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.In.F.Tag = "foo"
				w.In.F.Body = func(context.Context) error { return nil }
				var r1 string
				w.In.F.Results = []Result{{OutValueID: "r1", OutValuePtr: r1}}
				w.ExpOut.ErrStr = ErrInvalidResult.Error() + ": invalid out-value pointer type; tag=\"foo\" outValueID=\"r1\" outValuePtrType=\"string\""
				w.ExpOut.Err = ErrInvalidFunction
			}),
		tc.Copy().
			Given("result with nil out-value pointer").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.In.F.Tag = "foo"
				w.In.F.Body = func(context.Context) error { return nil }
				w.In.F.Results = []Result{{OutValueID: "r1", OutValuePtr: (*int)(nil)}}
				w.ExpOut.ErrStr = ErrInvalidResult.Error() + ": nil out-value pointer; tag=\"foo\" outValueID=\"r1\""
				w.ExpOut.Err = ErrInvalidFunction
			}),
		tc.Copy().
			Given("hook with empty in-value id").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.In.F.Tag = "foo"
				w.In.F.Body = func(context.Context) error { return nil }
				var a1 int
				w.In.F.Hooks = []Hook{{InValuePtr: &a1}}
				w.ExpOut.ErrStr = ErrInvalidHook.Error() + ": empty in-value id; tag=\"foo\""
				w.ExpOut.Err = ErrInvalidFunction
			}),
		tc.Copy().
			Given("hook without in-value pointer").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.In.F.Tag = "foo"
				w.In.F.Body = func(context.Context) error { return nil }
				w.In.F.Hooks = []Hook{{InValueID: "r1"}}
				w.ExpOut.ErrStr = ErrInvalidHook.Error() + ": no in-value pointer; tag=\"foo\" inValueID=\"r1\""
				w.ExpOut.Err = ErrInvalidFunction
			}),
		tc.Copy().
			Given("hook with invalid in-value pointer type").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.In.F.Tag = "foo"
				w.In.F.Body = func(context.Context) error { return nil }
				var r1 string
				w.In.F.Hooks = []Hook{{InValueID: "r1", InValuePtr: r1}}
				w.ExpOut.ErrStr = ErrInvalidHook.Error() + ": invalid in-value pointer type; tag=\"foo\" inValueID=\"r1\" inValuePtrType=\"string\""
				w.ExpOut.Err = ErrInvalidFunction
			}),
		tc.Copy().
			Given("hook with nil in-value pointer").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.In.F.Tag = "foo"
				w.In.F.Body = func(context.Context) error { return nil }
				w.In.F.Hooks = []Hook{{InValueID: "r1", InValuePtr: (*int)(nil)}}
				w.ExpOut.ErrStr = ErrInvalidHook.Error() + ": nil in-value pointer; tag=\"foo\" inValueID=\"r1\""
				w.ExpOut.Err = ErrInvalidFunction
			}),
		tc.Copy().
			Given("hook with nil callback pointer").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.In.F.Tag = "foo"
				w.In.F.Body = func(context.Context) error { return nil }
				var r1 *int
				w.In.F.Hooks = []Hook{{InValueID: "r1", InValuePtr: &r1}}
				w.ExpOut.ErrStr = ErrInvalidHook.Error() + ": nil callback pointer; tag=\"foo\" inValueID=\"r1\""
				w.ExpOut.Err = ErrInvalidFunction
			}),
		tc.Copy().
			Then("should succeed").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.In.F.Tag = "foo"
				w.In.F.Body = func(context.Context) error { return nil }
				var a1 *int
				w.In.F.Arguments = []Argument{{InValueID: "a1", InValuePtr: &a1}}
				var r1 *int
				w.In.F.Results = []Result{{OutValueID: "r1", OutValuePtr: &r1}}
				cb := func(context.Context) error { return nil }
				w.In.F.Hooks = []Hook{{InValueID: "r1", InValuePtr: &r1, CallbackPtr: &cb}}
			}),
	)
}

func TestProgram_MustAddFunctions(t *testing.T) {
	assert.Panics(t, func() {
		var p Program
		p.MustAddFunctions(Function{})
	})
}

func TestProgram_Run(t *testing.T) {
	type Workspace struct {
		P  Program
		In struct {
			Ctx context.Context
		}
		ExpOut, ActOut struct {
			ErrStr string
			Err    error
		}
	}
	tc := testcase.New().
		Step(1, func(t *testing.T, w *Workspace) {
			w.In.Ctx = context.Background()
		}).
		Step(2, func(t *testing.T, w *Workspace) {
			err := w.P.Run(w.In.Ctx)
			if err != nil {
				w.ActOut.ErrStr = err.Error()
				for err2 := errors.Unwrap(err); err2 != nil; err, err2 = err2, errors.Unwrap(err2) {
				}
				w.ActOut.Err = err
			}
		}).
		Step(3, func(t *testing.T, w *Workspace) {
			assert.Equal(t, w.ExpOut, w.ActOut)
		})
	testcase.RunListParallel(t,
		tc.Copy().
			Given("results with identical out-value ids").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				var var1 int
				var var2 int
				w.P.MustAddFunctions(
					Function{
						Tag: "foo",
						Results: []Result{
							{OutValueID: "var", OutValuePtr: &var1},
						},
						Body: func(context.Context) error { return nil },
					},
					Function{
						Tag: "bar",
						Results: []Result{
							{OutValueID: "var", OutValuePtr: &var2},
						},
						Body: func(context.Context) error { return nil },
					},
				)
				w.ExpOut.ErrStr = ErrValueAlreadyExists.Error() + "; tag1=\"bar\" tag2=\"foo\" outValueID=\"var\""
				w.ExpOut.Err = ErrValueAlreadyExists
			}),
		tc.Copy().
			Given("in-value of argument not found by id").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				{
					var x int
					w.P.MustAddFunctions(Function{
						Tag: "foo",
						Results: []Result{
							{OutValueID: "x", OutValuePtr: &x},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				{
					var x, y int
					w.P.MustAddFunctions(Function{
						Tag: "bar",
						Arguments: []Argument{
							{InValueID: "x", InValuePtr: &x},
							{InValueID: "y", InValuePtr: &y},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				w.ExpOut.ErrStr = ErrValueNotFound.Error() + "; tag=\"bar\" inValueID=\"y\""
				w.ExpOut.Err = ErrValueNotFound
			}),
		tc.Copy().
			Given("in-value of optional argument not found by id").
			Then("should not fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				{
					var x int
					w.P.MustAddFunctions(Function{
						Tag: "foo",
						Results: []Result{
							{OutValueID: "x", OutValuePtr: &x},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				{
					var x, y int
					w.P.MustAddFunctions(Function{
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
			Given("in-value type of argument and out-value type of result not matched").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				{
					var x int
					w.P.MustAddFunctions(Function{
						Tag: "foo",
						Results: []Result{
							{OutValueID: "x", OutValuePtr: &x},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				{
					var x string
					w.P.MustAddFunctions(Function{
						Tag: "bar",
						Arguments: []Argument{
							{InValueID: "x", InValuePtr: &x},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				w.ExpOut.ErrStr = ErrValueTypeMismatch.Error() + "; tag1=\"bar\" tag2=\"foo\" valueID=\"x\" inValueType=\"string\" outValueType=\"int\""
				w.ExpOut.Err = ErrValueTypeMismatch
			}),
		tc.Copy().
			Given("in-value of hook not found by id").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				{
					var x int
					w.P.MustAddFunctions(Function{
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
					w.P.MustAddFunctions(Function{
						Tag: "bar",
						Hooks: []Hook{
							{InValueID: "x", InValuePtr: &x, CallbackPtr: &cb},
							{InValueID: "y", InValuePtr: &y, CallbackPtr: &cb},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				w.ExpOut.ErrStr = ErrValueNotFound.Error() + "; tag=\"bar\" inValueID=\"y\""
				w.ExpOut.Err = ErrValueNotFound
			}),
		tc.Copy().
			Given("in-value type of hook and out-value type of result not matched").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				{
					var x int
					w.P.MustAddFunctions(Function{
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
					w.P.MustAddFunctions(Function{
						Tag: "bar",
						Hooks: []Hook{
							{InValueID: "x", InValuePtr: &x, CallbackPtr: &cb},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				w.ExpOut.Err = ErrValueTypeMismatch
				w.ExpOut.ErrStr = w.ExpOut.Err.Error() + "; tag1=\"bar\" tag2=\"foo\" valueID=\"x\" inValueType=\"string\" outValueType=\"int\""
			}),
		tc.Copy().
			Given("circular dependencies (1)").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				var x int
				w.P.MustAddFunctions(Function{
					Tag: "foo",
					Arguments: []Argument{
						{InValueID: "x", InValuePtr: &x},
					},
					Results: []Result{
						{OutValueID: "x", OutValuePtr: &x},
					},
					Body: func(context.Context) error { return nil },
				})
				w.ExpOut.ErrStr = ErrCircularDependencies.Error() +
					"; {tag: \"foo\", argument: \"x\"} => {tag: \"foo\"}"
				w.ExpOut.Err = ErrCircularDependencies
			}),
		tc.Copy().
			Given("circular dependencies (2)").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				var x int
				cb := func(context.Context) error { return nil }
				w.P.MustAddFunctions(Function{
					Tag: "foo",
					Results: []Result{
						{OutValueID: "x", OutValuePtr: &x},
					},
					Hooks: []Hook{
						{InValueID: "x", InValuePtr: &x, CallbackPtr: &cb},
					},
					Body: func(context.Context) error { return nil },
				})
				w.ExpOut.ErrStr = ErrCircularDependencies.Error() +
					"; {tag: \"foo\", hook: \"x\"} => {tag: \"foo\"}"
				w.ExpOut.Err = ErrCircularDependencies
			}),
		tc.Copy().
			Given("circular dependencies (3)").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				{
					var x, y int
					w.P.MustAddFunctions(Function{
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
					w.P.MustAddFunctions(Function{
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
				w.ExpOut.ErrStr = ErrCircularDependencies.Error() +
					"; {tag: \"foo\", argument: \"x\"} => {tag: \"bar\", argument: \"y\"} => {tag: \"foo\"}"
				w.ExpOut.Err = ErrCircularDependencies
			}),
		tc.Copy().
			Given("circular dependencies (4)").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				{
					var x int
					w.P.MustAddFunctions(Function{
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
					w.P.MustAddFunctions(Function{
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
				w.ExpOut.ErrStr = ErrCircularDependencies.Error() +
					"; {tag: \"foo\", hook: \"x\"} => {tag: \"bar\", argument: \"x\"} => {tag: \"foo\"}"
				w.ExpOut.Err = ErrCircularDependencies
			}),
		tc.Copy().
			Given("function body returning error").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.P.MustAddFunctions(Function{
					Tag:  "foo",
					Body: func(context.Context) error { return context.DeadlineExceeded },
				})
				w.ExpOut.ErrStr = "call function; tag=\"foo\": " + context.DeadlineExceeded.Error()
				w.ExpOut.Err = context.DeadlineExceeded
			}),
		tc.Copy().
			Given("function body not provisioning cleanup").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				var x int
				var c func()
				w.P.MustAddFunctions(Function{
					Tag: "foo",
					Results: []Result{
						{OutValueID: "x", OutValuePtr: &x},
					},
					CleanupPtr: &c,
					Body:       func(context.Context) error { return nil },
				})
				w.ExpOut.ErrStr = ErrNilCleanup.Error() + "; tag=\"foo\""
				w.ExpOut.Err = ErrNilCleanup
			}),
		tc.Copy().
			Given("function body not provisioning callback").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				{
					var x int
					w.P.MustAddFunctions(Function{
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
					w.P.MustAddFunctions(Function{
						Tag: "bar",
						Hooks: []Hook{
							{InValueID: "x", InValuePtr: &x, CallbackPtr: &cb},
						},
						Body: func(context.Context) error { return nil },
					})
				}
				w.ExpOut.ErrStr = ErrNilCallback.Error() + "; tag=\"bar\" inValueID=\"x\""
				w.ExpOut.Err = ErrNilCallback
			}),
		tc.Copy().
			Given("function body provisioning callback returning error").
			Then("should fail").
			Step(0.5, func(t *testing.T, w *Workspace) {
				{
					var x int
					w.P.MustAddFunctions(Function{
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
					w.P.MustAddFunctions(Function{
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
				w.ExpOut.ErrStr = "do callback; tag=\"bar\" inValueID=\"x\": " + context.DeadlineExceeded.Error()
				w.ExpOut.Err = context.DeadlineExceeded
			}),
		tc.Copy().
			Then("should succeed").
			Step(0.5, func(t *testing.T, w *Workspace) {
				var s string
				t.Cleanup(func() {
					if !t.Failed() {
						assert.Equal(t, "12345", s)
					}
				})
				{
					var x, y int
					w.P.MustAddFunctions(Function{
						Tag: "foo",
						Arguments: []Argument{
							{InValueID: "x", InValuePtr: &x},
							{InValueID: "y", InValuePtr: &y},
						},
						Body: func(context.Context) error {
							s += "5"
							assert.Equal(t, 101, x)
							assert.Equal(t, 404, y)
							return nil
						},
					})
				}
				{
					var x, y int
					w.P.MustAddFunctions(Function{
						Tag: "bar",
						Arguments: []Argument{
							{InValueID: "x", InValuePtr: &x},
						},
						Results: []Result{
							{OutValueID: "y", OutValuePtr: &y},
						},
						Body: func(context.Context) error {
							s += "3"
							assert.Equal(t, 101, x)
							y = 404
							return nil
						},
					})
				}
				{
					var x int
					w.P.MustAddFunctions(Function{
						Tag: "baz",
						Results: []Result{
							{OutValueID: "x", OutValuePtr: &x},
						},
						Body: func(context.Context) error {
							s += "1"
							x = 101
							return nil
						},
					})
				}
				{
					var y int
					var cb func(context.Context) error
					w.P.MustAddFunctions(Function{
						Tag: "qux",
						Hooks: []Hook{
							{InValueID: "y", InValuePtr: &y, CallbackPtr: &cb},
						},
						Body: func(context.Context) error {
							s += "2"
							cb = func(context.Context) error {
								s += "4"
								assert.Equal(t, 404, y)
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
		p.MustAddFunctions(Function{
			Tag: "foo",
			Body: func(context.Context) error {
				return errors.New("")
			},
		})
		p.MustRun(context.Background())
	})
}

func TestProgram_Clean(t *testing.T) {
	type Workspace struct {
		P Program
	}
	tc := testcase.New().
		Step(1, func(t *testing.T, w *Workspace) {
			w.P.Clean()
		})
	testcase.RunListParallel(t,
		tc.Copy().
			Given("successful Program.Run() and cleanups provisioned").
			Then("should do cleanups").
			Step(0.5, func(t *testing.T, w *Workspace) {
				var s string
				t.Cleanup(func() {
					if !t.Failed() {
						assert.Equal(t, "12", s)
					}
				})
				{
					var x int
					var c func()
					w.P.MustAddFunctions(Function{
						Tag: "foo",
						Results: []Result{
							{OutValueID: "x", OutValuePtr: &x},
						},
						CleanupPtr: &c,
						Body: func(context.Context) error {
							c = func() { s += "2" }
							return nil
						},
					})
				}
				{
					var x int
					var c func()
					w.P.MustAddFunctions(Function{
						Tag: "bar",
						Arguments: []Argument{
							{InValueID: "x", InValuePtr: &x},
						},
						CleanupPtr: &c,
						Body: func(context.Context) error {
							c = func() { s += "1" }
							return nil
						},
					})
				}
				err := w.P.Run(context.Background())
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			}),
		tc.Copy().
			Given("failed Program.Run() and cleanups provisioned").
			Then("should do cleanups").
			Step(0.5, func(t *testing.T, w *Workspace) {
				var s string
				t.Cleanup(func() {
					if !t.Failed() {
						assert.Equal(t, "12", s)
					}
				})
				{
					var x int
					var c func()
					w.P.MustAddFunctions(Function{
						Tag: "foo",
						Results: []Result{
							{OutValueID: "x", OutValuePtr: &x},
						},
						CleanupPtr: &c,
						Body: func(context.Context) error {
							c = func() { s += "2" }
							return nil
						},
					})
				}
				{
					var x int
					var c func()
					w.P.MustAddFunctions(Function{
						Tag: "bar",
						Arguments: []Argument{
							{InValueID: "x", InValuePtr: &x},
						},
						CleanupPtr: &c,
						Body: func(context.Context) error {
							c = func() { s += "1" }
							return context.DeadlineExceeded
						},
					})
				}
				err := w.P.Run(context.Background())
				if !assert.ErrorIs(t, err, context.DeadlineExceeded) {
					t.FailNow()
				}
			}),
		tc.Copy().
			Given("cleanups not yet provisioned").
			Then("should not do cleanups").
			Step(0.5, func(t *testing.T, w *Workspace) {
				var s string
				t.Cleanup(func() {
					if !t.Failed() {
						assert.Equal(t, "", s)
					}
				})
				{
					var x int
					var c func()
					w.P.MustAddFunctions(Function{
						Tag: "foo",
						Results: []Result{
							{OutValueID: "x", OutValuePtr: &x},
						},
						CleanupPtr: &c,
						Body: func(context.Context) error {
							c = func() { s += "2" }
							return nil
						},
					})
				}
				{
					var x int
					var c func()
					w.P.MustAddFunctions(Function{
						Tag: "bar",
						Arguments: []Argument{
							{InValueID: "x", InValuePtr: &x},
						},
						CleanupPtr: &c,
						Body: func(context.Context) error {
							c = func() { s += "1" }
							return nil
						},
					})
				}
			}),
	)
}
