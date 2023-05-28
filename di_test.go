package di_test

import (
	"context"
	"testing"

	. "github.com/go-tk/di"
	"github.com/go-tk/testcase"
	"github.com/stretchr/testify/assert"
)

func TestArgument(t *testing.T) {
	type C struct {
		functionBuilder FunctionBuilder

		errStr string
		err    error
		repr   string
	}
	tc := testcase.New(func(t *testing.T, c *C) {
		testcase.Callback(t, "0")

		var p Program
		err := p.NewFunction(c.functionBuilder, Body(func(context.Context) error { return nil }))
		if c.errStr == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, c.errStr)
			if c.err != nil {
				assert.ErrorIs(t, err, c.err)
			}
		}
		if t.Failed() {
			t.FailNow()
		}
		assert.Equal(t, c.repr, p.DumpAsString())
	})

	tc.WithCallback("0", func(t *testing.T, c *C) {
		c.functionBuilder = Argument("", nil)
		c.err = ErrInvalidArgument
		c.errStr = c.err.Error() + `: empty value ref; functionName="github.com/go-tk/di_test.TestArgument.func1"`
		c.repr = `
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		c.functionBuilder = Argument("foo", nil)
		c.err = ErrInvalidArgument
		c.errStr = c.err.Error() + `: no value receiver; functionName="github.com/go-tk/di_test.TestArgument.func1" valueRef="foo"`
		c.repr = `
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		c.functionBuilder = Argument("foo", 0)
		c.err = ErrInvalidArgument
		c.errStr = c.err.Error() + `: invalid value receiver pointer; valueReceiverPtrType="int" functionName="github.com/go-tk/di_test.TestArgument.func1" valueRef="foo"`
		c.repr = `
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		c.functionBuilder = Argument("foo", (*string)(nil))
		c.err = ErrInvalidArgument
		c.errStr = c.err.Error() + `: no value receiver; functionName="github.com/go-tk/di_test.TestArgument.func1" valueRef="foo"`
		c.repr = `
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		c.functionBuilder = Argument("foo", new(string))
		c.repr = `
Function[0]:
	Index: 0
	Name: github.com/go-tk/di_test.TestArgument.func1
	ArgumentIndexes: [0]
	ResultIndexes: []
	HasBody: true
	HookIndexes: []
	HasCleanup: false
Argument[0]:
	FunctionIndex: 0
	ValueRef: foo
	HasValueReceiver: true
	IsOptional: false
	ResultIndex: -1
	ReceiveValueAddr: false
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		c.functionBuilder = OptionalArgument("foo", new(string))
		c.repr = `
Function[0]:
	Index: 0
	Name: github.com/go-tk/di_test.TestArgument.func1
	ArgumentIndexes: [0]
	ResultIndexes: []
	HasBody: true
	HookIndexes: []
	HasCleanup: false
Argument[0]:
	FunctionIndex: 0
	ValueRef: foo
	HasValueReceiver: true
	IsOptional: true
	ResultIndex: -1
	ReceiveValueAddr: false
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)
}

func TestHook(t *testing.T) {
	type C struct {
		functionBuilder FunctionBuilder

		errStr string
		err    error
		repr   string
	}
	tc := testcase.New(func(t *testing.T, c *C) {
		testcase.Callback(t, "0")

		var p Program
		err := p.NewFunction(c.functionBuilder, Body(func(context.Context) error { return nil }))
		if c.errStr == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, c.errStr)
			if c.err != nil {
				assert.ErrorIs(t, err, c.err)
			}
		}
		if t.Failed() {
			t.FailNow()
		}
		assert.Equal(t, c.repr, p.DumpAsString())
	})

	tc.WithCallback("0", func(t *testing.T, c *C) {
		c.functionBuilder = Hook("", nil, nil)
		c.err = ErrInvalidHook
		c.errStr = c.err.Error() + `: empty value ref; functionName="github.com/go-tk/di_test.TestHook.func1"`
		c.repr = `
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		c.functionBuilder = Hook("foo", nil, nil)
		c.err = ErrInvalidHook
		c.errStr = c.err.Error() + `: no value receiver; functionName="github.com/go-tk/di_test.TestHook.func1" valueRef="foo"`
		c.repr = `
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		c.functionBuilder = Hook("foo", 0, nil)
		c.err = ErrInvalidHook
		c.errStr = c.err.Error() + `: invalid value receiver pointer; valueReceiverPtrType="int" functionName="github.com/go-tk/di_test.TestHook.func1" valueRef="foo"`
		c.repr = `
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		c.functionBuilder = Hook("foo", (*string)(nil), nil)
		c.err = ErrInvalidHook
		c.errStr = c.err.Error() + `: no value receiver; functionName="github.com/go-tk/di_test.TestHook.func1" valueRef="foo"`
		c.repr = `
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		c.functionBuilder = Hook("foo", new(string), nil)
		c.err = ErrInvalidHook
		c.errStr = c.err.Error() + `: nil callback; functionName="github.com/go-tk/di_test.TestHook.func1" valueRef="foo"`
		c.repr = `
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		c.functionBuilder = Hook("foo", new(string), func(context.Context) error { return nil })
		c.repr = `
Function[0]:
	Index: 0
	Name: github.com/go-tk/di_test.TestHook.func1
	ArgumentIndexes: []
	ResultIndexes: []
	HasBody: true
	HookIndexes: [0]
	HasCleanup: false
Hook[0]:
	FunctionIndex: 0
	ValueRef: foo
	HasValueReceiver: true
	HasCallback: true
	ReceiveValueAddr: false
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)
}

func TestResult(t *testing.T) {
	type C struct {
		functionBuilder FunctionBuilder

		errStr string
		err    error
		repr   string
	}
	tc := testcase.New(func(t *testing.T, c *C) {
		testcase.Callback(t, "0")

		var p Program
		err := p.NewFunction(c.functionBuilder, Body(func(context.Context) error { return nil }))
		if c.errStr == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, c.errStr)
			if c.err != nil {
				assert.ErrorIs(t, err, c.err)
			}
		}
		if t.Failed() {
			t.FailNow()
		}
		assert.Equal(t, c.repr, p.DumpAsString())
	})

	tc.WithCallback("0", func(t *testing.T, c *C) {
		c.functionBuilder = Result("", nil)
		c.err = ErrInvalidResult
		c.errStr = c.err.Error() + `: empty value name; functionName="github.com/go-tk/di_test.TestResult.func1"`
		c.repr = `
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		c.functionBuilder = Result("foo", nil)
		c.err = ErrInvalidResult
		c.errStr = c.err.Error() + `: no value; functionName="github.com/go-tk/di_test.TestResult.func1" valueName="foo"`
		c.repr = `
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		c.functionBuilder = Result("foo", 0)
		c.err = ErrInvalidResult
		c.errStr = c.err.Error() + `: invalid value pointer; valuePtrType="int" functionName="github.com/go-tk/di_test.TestResult.func1" valueName="foo"`
		c.repr = `
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		c.functionBuilder = Result("foo", (*string)(nil))
		c.err = ErrInvalidResult
		c.errStr = c.err.Error() + `: no value; functionName="github.com/go-tk/di_test.TestResult.func1" valueName="foo"`
		c.repr = `
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		c.functionBuilder = Result("foo", new(string))
		c.repr = `
Function[0]:
	Index: 0
	Name: github.com/go-tk/di_test.TestResult.func1
	ArgumentIndexes: []
	ResultIndexes: [0]
	HasBody: true
	HookIndexes: []
	HasCleanup: false
Result[0]:
	FunctionIndex: 0
	ValueName: foo
	HasValue: true
	HookIndexes: []
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)
}

func TestBody(t *testing.T) {
	type C struct {
		functionBuilder FunctionBuilder

		errStr string
		err    error
		repr   string
	}
	tc := testcase.New(func(t *testing.T, c *C) {
		testcase.Callback(t, "0")

		var p Program
		err := p.NewFunction(c.functionBuilder, Body(func(context.Context) error { return nil }))
		if c.errStr == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, c.errStr)
			if c.err != nil {
				assert.ErrorIs(t, err, c.err)
			}
		}
		if t.Failed() {
			t.FailNow()
		}
		assert.Equal(t, c.repr, p.DumpAsString())
	})

	tc.WithCallback("0", func(t *testing.T, c *C) {
		c.functionBuilder = Body(nil)
		c.err = ErrNilBody
		c.errStr = c.err.Error() + `; functionName="github.com/go-tk/di_test.TestBody.func1"`
		c.repr = `
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		c.functionBuilder = Body(func(context.Context) error { return nil })
		c.repr = `
Function[0]:
	Index: 0
	Name: github.com/go-tk/di_test.TestBody.func1
	ArgumentIndexes: []
	ResultIndexes: []
	HasBody: true
	HookIndexes: []
	HasCleanup: false
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)
}

func TestCleanup(t *testing.T) {
	type C struct {
		functionBuilder FunctionBuilder

		errStr string
		err    error
		repr   string
	}
	tc := testcase.New(func(t *testing.T, c *C) {
		testcase.Callback(t, "0")

		var p Program
		err := p.NewFunction(c.functionBuilder, Body(func(context.Context) error { return nil }))
		if c.errStr == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, c.errStr)
			if c.err != nil {
				assert.ErrorIs(t, err, c.err)
			}
		}
		if t.Failed() {
			t.FailNow()
		}
		assert.Equal(t, c.repr, p.DumpAsString())
	})

	tc.WithCallback("0", func(t *testing.T, c *C) {
		c.functionBuilder = Cleanup(nil)
		c.err = ErrNilCleanup
		c.errStr = c.err.Error() + `; functionName="github.com/go-tk/di_test.TestCleanup.func1"`
		c.repr = `
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		c.functionBuilder = Cleanup(func() {})
		c.repr = `
Function[0]:
	Index: 0
	Name: github.com/go-tk/di_test.TestCleanup.func1
	ArgumentIndexes: []
	ResultIndexes: []
	HasBody: true
	HookIndexes: []
	HasCleanup: true
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)
}

func TestNewFunction(t *testing.T) {
	type C struct {
		functionBuilders []FunctionBuilder

		errStr string
		err    error
		repr   string
	}
	tc := testcase.New(func(t *testing.T, c *C) {
		testcase.Callback(t, "0")

		var p Program
		err := p.NewFunction(c.functionBuilders...)
		if c.errStr == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, c.errStr)
			if c.err != nil {
				assert.ErrorIs(t, err, c.err)
			}
		}
		if t.Failed() {
			t.FailNow()
		}
		assert.Equal(t, c.repr, p.DumpAsString())
	})

	tc.WithCallback("0", func(t *testing.T, c *C) {
		c.err = ErrBodyRequired
		c.errStr = c.err.Error() + `; functionName="github.com/go-tk/di_test.TestNewFunction.func1"`
		c.repr = `
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		var (
			arg1 int
			arg2 string
			res1 float64
			res2 byte
			var1 int32
			var2 int64
		)
		c.functionBuilders = []FunctionBuilder{
			Argument("arg1", &arg1),
			Argument("arg2", &arg2),
			Result("res1", &res1),
			Result("res2", &res2),
			Body(func(context.Context) error { return nil }),
			Hook("var1", &var1, func(context.Context) error { return nil }),
			Hook("var2", &var2, func(context.Context) error { return nil }),
			Cleanup(func() {}),
		}
		c.repr = `
Function[0]:
	Index: 0
	Name: github.com/go-tk/di_test.TestNewFunction.func1
	ArgumentIndexes: [0 1]
	ResultIndexes: [0 1]
	HasBody: true
	HookIndexes: [0 1]
	HasCleanup: true
Argument[0]:
	FunctionIndex: 0
	ValueRef: arg1
	HasValueReceiver: true
	IsOptional: false
	ResultIndex: -1
	ReceiveValueAddr: false
Argument[1]:
	FunctionIndex: 0
	ValueRef: arg2
	HasValueReceiver: true
	IsOptional: false
	ResultIndex: -1
	ReceiveValueAddr: false
Result[0]:
	FunctionIndex: 0
	ValueName: res1
	HasValue: true
	HookIndexes: []
Result[1]:
	FunctionIndex: 0
	ValueName: res2
	HasValue: true
	HookIndexes: []
Hook[0]:
	FunctionIndex: 0
	ValueRef: var1
	HasValueReceiver: true
	HasCallback: true
	ReceiveValueAddr: false
Hook[1]:
	FunctionIndex: 0
	ValueRef: var2
	HasValueReceiver: true
	HasCallback: true
	ReceiveValueAddr: false
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)
}

func TestProgram_MustNewFunction(t *testing.T) {
	{
		var p Program
		p.MustNewFunction(Body(func(context.Context) error { return nil }))
	}
	assert.PanicsWithValue(t, `new function: di: body required; functionName="github.com/go-tk/di_test.TestProgram_MustNewFunction.func2"`, func() {
		var p Program
		p.MustNewFunction()
	})
}

func TestProgram_Run(t *testing.T) {
	type C struct {
		p   *Program
		ctx *context.Context

		errStr string
		err    error
		repr   string
	}
	tc := testcase.New(func(t *testing.T, c *C) {
		var p Program
		ctx := context.Background()
		c.p = &p
		c.ctx = &ctx

		testcase.Callback(t, "0")

		err := p.Run(ctx)
		if c.errStr == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, c.errStr)
			if c.err != nil {
				assert.ErrorIs(t, err, c.err)
			}
		}
		if t.Failed() {
			t.FailNow()
		}
		assert.Equal(t, c.repr, p.DumpAsString())
	})

	tc.WithCallback("0", func(t *testing.T, c *C) {
		func() {
			var res1 int
			err := c.p.NewFunction(Result("res1", &res1), Body(func(context.Context) error { return nil }))
			if err != nil {
				t.Fatal(err)
			}
		}()
		func() {
			var res1 int
			err := c.p.NewFunction(Result("res1", &res1), Body(func(context.Context) error { return nil }))
			if err != nil {
				t.Fatal(err)
			}
		}()
		c.err = ErrDuplicateValueName
		c.errStr = c.err.Error() + `; valueName="res1" functionName1="github.com/go-tk/di_test.TestProgram_Run.func2.2" functionName2="github.com/go-tk/di_test.TestProgram_Run.func2.1"`
		c.repr = `
Function[0]:
	Index: 0
	Name: github.com/go-tk/di_test.TestProgram_Run.func2.1
	ArgumentIndexes: []
	ResultIndexes: [0]
	HasBody: true
	HookIndexes: []
	HasCleanup: false
Function[1]:
	Index: 1
	Name: github.com/go-tk/di_test.TestProgram_Run.func2.2
	ArgumentIndexes: []
	ResultIndexes: [1]
	HasBody: true
	HookIndexes: []
	HasCleanup: false
Result[0]:
	FunctionIndex: 0
	ValueName: res1
	HasValue: true
	HookIndexes: []
Result[1]:
	FunctionIndex: 1
	ValueName: res1
	HasValue: true
	HookIndexes: []
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		func() {
			var arg int
			err := c.p.NewFunction(Argument("arg", &arg), Body(func(context.Context) error { return nil }))
			if err != nil {
				t.Fatal(err)
			}
		}()
		c.err = ErrValueNotFound
		c.errStr = c.err.Error() + `; valueRef="arg" functionName="github.com/go-tk/di_test.TestProgram_Run.func3.1"`
		c.repr = `
Function[0]:
	Index: 0
	Name: github.com/go-tk/di_test.TestProgram_Run.func3.1
	ArgumentIndexes: [0]
	ResultIndexes: []
	HasBody: true
	HookIndexes: []
	HasCleanup: false
Argument[0]:
	FunctionIndex: 0
	ValueRef: arg
	HasValueReceiver: true
	IsOptional: false
	ResultIndex: -1
	ReceiveValueAddr: false
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		func() {
			var arg int
			err := c.p.NewFunction(OptionalArgument("arg", &arg), Body(func(context.Context) error { return nil }))
			if err != nil {
				t.Fatal(err)
			}
		}()
		c.repr = `
Function[0]:
	Index: 0
	Name: github.com/go-tk/di_test.TestProgram_Run.func4.1
	ArgumentIndexes: [0]
	ResultIndexes: []
	HasBody: true
	HookIndexes: []
	HasCleanup: false
Argument[0]:
	FunctionIndex: 0
	ValueRef: arg
	HasValueReceiver: true
	IsOptional: true
	ResultIndex: -1
	ReceiveValueAddr: false
SortedFunctionIndexes: [0]
CalledFunctionCount: 1
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		func() {
			var val int
			err := c.p.NewFunction(Result("val", &val), Body(func(context.Context) error { return nil }))
			if err != nil {
				t.Fatal(err)
			}
		}()
		func() {
			var val string
			err := c.p.NewFunction(Argument("val", &val), Body(func(context.Context) error { return nil }))
			if err != nil {
				t.Fatal(err)
			}
		}()
		c.err = ErrIncompatibleValueReceiver
		c.errStr = c.err.Error() + `; valueReceiverType="string" valueType="int" valueRef="val" functionName="github.com/go-tk/di_test.TestProgram_Run.func5.2"`
		c.repr = `
Function[0]:
	Index: 0
	Name: github.com/go-tk/di_test.TestProgram_Run.func5.1
	ArgumentIndexes: []
	ResultIndexes: [0]
	HasBody: true
	HookIndexes: []
	HasCleanup: false
Function[1]:
	Index: 1
	Name: github.com/go-tk/di_test.TestProgram_Run.func5.2
	ArgumentIndexes: [0]
	ResultIndexes: []
	HasBody: true
	HookIndexes: []
	HasCleanup: false
Argument[0]:
	FunctionIndex: 1
	ValueRef: val
	HasValueReceiver: true
	IsOptional: false
	ResultIndex: -1
	ReceiveValueAddr: false
Result[0]:
	FunctionIndex: 0
	ValueName: val
	HasValue: true
	HookIndexes: []
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		func() {
			var hook int
			err := c.p.NewFunction(Body(func(context.Context) error { return nil }),
				Hook("hook", &hook, func(context.Context) error { return nil }))
			if err != nil {
				t.Fatal(err)
			}
		}()
		c.err = ErrValueNotFound
		c.errStr = c.err.Error() + `; valueRef="hook" functionName="github.com/go-tk/di_test.TestProgram_Run.func6.1"`
		c.repr = `
Function[0]:
	Index: 0
	Name: github.com/go-tk/di_test.TestProgram_Run.func6.1
	ArgumentIndexes: []
	ResultIndexes: []
	HasBody: true
	HookIndexes: [0]
	HasCleanup: false
Hook[0]:
	FunctionIndex: 0
	ValueRef: hook
	HasValueReceiver: true
	HasCallback: true
	ReceiveValueAddr: false
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		func() {
			var val int
			err := c.p.NewFunction(Result("val", &val), Body(func(context.Context) error { return nil }))
			if err != nil {
				t.Fatal(err)
			}
		}()
		func() {
			var val string
			err := c.p.NewFunction(Body(func(context.Context) error { return nil }),
				Hook("val", &val, func(context.Context) error { return nil }))
			if err != nil {
				t.Fatal(err)
			}
		}()
		c.err = ErrIncompatibleValueReceiver
		c.errStr = c.err.Error() + `; valueReceiverType="string" valueType="int" valueRef="val" functionName="github.com/go-tk/di_test.TestProgram_Run.func7.2"`
		c.repr = `
Function[0]:
	Index: 0
	Name: github.com/go-tk/di_test.TestProgram_Run.func7.1
	ArgumentIndexes: []
	ResultIndexes: [0]
	HasBody: true
	HookIndexes: []
	HasCleanup: false
Function[1]:
	Index: 1
	Name: github.com/go-tk/di_test.TestProgram_Run.func7.2
	ArgumentIndexes: []
	ResultIndexes: []
	HasBody: true
	HookIndexes: [0]
	HasCleanup: false
Result[0]:
	FunctionIndex: 0
	ValueName: val
	HasValue: true
	HookIndexes: []
Hook[0]:
	FunctionIndex: 1
	ValueRef: val
	HasValueReceiver: true
	HasCallback: true
	ReceiveValueAddr: false
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		func() {
			var val int
			var pval *int
			err := c.p.NewFunction(Argument("val", &pval), Result("val", &val), Body(func(context.Context) error { return nil }))
			if err != nil {
				t.Fatal(err)
			}
		}()
		c.err = ErrCircularDependencies
		c.errStr = c.err.Error() + `; path="github.com/go-tk/di_test.TestProgram_Run.func8.1@argument:val => github.com/go-tk/di_test.TestProgram_Run.func8.1"`
		c.repr = `
Function[0]:
	Index: -1
	Name: github.com/go-tk/di_test.TestProgram_Run.func8.1
	ArgumentIndexes: [0]
	ResultIndexes: [0]
	HasBody: true
	HookIndexes: []
	HasCleanup: false
Argument[0]:
	FunctionIndex: 0
	ValueRef: val
	HasValueReceiver: true
	IsOptional: false
	ResultIndex: 0
	ReceiveValueAddr: true
Result[0]:
	FunctionIndex: 0
	ValueName: val
	HasValue: true
	HookIndexes: []
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		func() {
			var val int
			var pval *int
			err := c.p.NewFunction(Result("val", &val), Body(func(context.Context) error { return nil }),
				Hook("val", &pval, func(context.Context) error { return nil }))
			if err != nil {
				t.Fatal(err)
			}
		}()
		c.err = ErrCircularDependencies
		c.errStr = c.err.Error() + `; path="github.com/go-tk/di_test.TestProgram_Run.func9.1@hook:val => github.com/go-tk/di_test.TestProgram_Run.func9.1"`
		c.repr = `
Function[0]:
	Index: -1
	Name: github.com/go-tk/di_test.TestProgram_Run.func9.1
	ArgumentIndexes: []
	ResultIndexes: [0]
	HasBody: true
	HookIndexes: [0]
	HasCleanup: false
Result[0]:
	FunctionIndex: 0
	ValueName: val
	HasValue: true
	HookIndexes: [0]
Hook[0]:
	FunctionIndex: 0
	ValueRef: val
	HasValueReceiver: true
	HasCallback: true
	ReceiveValueAddr: true
SortedFunctionIndexes: []
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		func() {
			var val int
			err := c.p.NewFunction(Argument("val", &val), Body(func(context.Context) error { return nil }))
			if err != nil {
				t.Fatal(err)
			}
		}()
		func() {
			var val int
			err := c.p.NewFunction(Argument("val", &val), Body(func(context.Context) error { return nil }))
			if err != nil {
				t.Fatal(err)
			}
		}()
		func() {
			var val int
			err := c.p.NewFunction(Result("val", &val), Body(func(context.Context) error { return nil }))
			if err != nil {
				t.Fatal(err)
			}
		}()
		c.repr = `
Function[0]:
	Index: 0
	Name: github.com/go-tk/di_test.TestProgram_Run.func10.1
	ArgumentIndexes: [0]
	ResultIndexes: []
	HasBody: true
	HookIndexes: []
	HasCleanup: false
Function[1]:
	Index: 1
	Name: github.com/go-tk/di_test.TestProgram_Run.func10.2
	ArgumentIndexes: [1]
	ResultIndexes: []
	HasBody: true
	HookIndexes: []
	HasCleanup: false
Function[2]:
	Index: 2
	Name: github.com/go-tk/di_test.TestProgram_Run.func10.3
	ArgumentIndexes: []
	ResultIndexes: [0]
	HasBody: true
	HookIndexes: []
	HasCleanup: false
Argument[0]:
	FunctionIndex: 0
	ValueRef: val
	HasValueReceiver: true
	IsOptional: false
	ResultIndex: 0
	ReceiveValueAddr: false
Argument[1]:
	FunctionIndex: 1
	ValueRef: val
	HasValueReceiver: true
	IsOptional: false
	ResultIndex: 0
	ReceiveValueAddr: false
Result[0]:
	FunctionIndex: 2
	ValueName: val
	HasValue: true
	HookIndexes: []
SortedFunctionIndexes: [2 0 1]
CalledFunctionCount: 3
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		func() {
			var (
				x int
				y *int
			)
			err := c.p.NewFunction(Argument("x", &x), Argument("y", &y), Body(func(context.Context) error {
				x += 1
				*y -= 1
				return nil
			}))
			if err != nil {
				t.Fatal(err)
			}
		}()
		func() {
			var (
				x int
				y int
			)
			err := c.p.NewFunction(Argument("x", &x), Argument("y", &y), Body(func(context.Context) error {
				assert.Equal(t, 100, x)
				assert.Equal(t, 98, y)
				y -= 1
				return nil
			}))
			if err != nil {
				t.Fatal(err)
			}
		}()
		func() {
			var (
				x int
				y int
			)
			err := c.p.NewFunction(Result("x", &x), Result("y", &y), Body(func(context.Context) error {
				x = 100
				y = 99
				return nil
			}))
			if err != nil {
				t.Fatal(err)
			}
		}()
		c.repr = `
Function[0]:
	Index: 0
	Name: github.com/go-tk/di_test.TestProgram_Run.func11.1
	ArgumentIndexes: [0 1]
	ResultIndexes: []
	HasBody: true
	HookIndexes: []
	HasCleanup: false
Function[1]:
	Index: 1
	Name: github.com/go-tk/di_test.TestProgram_Run.func11.2
	ArgumentIndexes: [2 3]
	ResultIndexes: []
	HasBody: true
	HookIndexes: []
	HasCleanup: false
Function[2]:
	Index: 2
	Name: github.com/go-tk/di_test.TestProgram_Run.func11.3
	ArgumentIndexes: []
	ResultIndexes: [0 1]
	HasBody: true
	HookIndexes: []
	HasCleanup: false
Argument[0]:
	FunctionIndex: 0
	ValueRef: x
	HasValueReceiver: true
	IsOptional: false
	ResultIndex: 0
	ReceiveValueAddr: false
Argument[1]:
	FunctionIndex: 0
	ValueRef: y
	HasValueReceiver: true
	IsOptional: false
	ResultIndex: 1
	ReceiveValueAddr: true
Argument[2]:
	FunctionIndex: 1
	ValueRef: x
	HasValueReceiver: true
	IsOptional: false
	ResultIndex: 0
	ReceiveValueAddr: false
Argument[3]:
	FunctionIndex: 1
	ValueRef: y
	HasValueReceiver: true
	IsOptional: false
	ResultIndex: 1
	ReceiveValueAddr: false
Result[0]:
	FunctionIndex: 2
	ValueName: x
	HasValue: true
	HookIndexes: []
Result[1]:
	FunctionIndex: 2
	ValueName: y
	HasValue: true
	HookIndexes: []
SortedFunctionIndexes: [2 0 1]
CalledFunctionCount: 3
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		func() {
			var (
				x int
				y int
			)
			err := c.p.NewFunction(Argument("x", &x), Argument("y", &y), Body(func(context.Context) error {
				assert.Equal(t, 100, x)
				assert.Equal(t, 98, y)
				y -= 1
				return nil
			}))
			if err != nil {
				t.Fatal(err)
			}
		}()
		func() {
			var (
				x int
				y int
			)
			err := c.p.NewFunction(Result("x", &x), Result("y", &y), Body(func(context.Context) error {
				x = 100
				y = 99
				return nil
			}))
			if err != nil {
				t.Fatal(err)
			}
		}()
		func() {
			var (
				x int
				y *int
			)
			err := c.p.NewFunction(Body(func(context.Context) error { return nil }),
				Hook("x", &x, func(context.Context) error { x += 1; return nil }),
				Hook("y", &y, func(context.Context) error { *y -= 1; return nil }))

			if err != nil {
				t.Fatal(err)
			}
		}()
		c.repr = `
Function[0]:
	Index: 0
	Name: github.com/go-tk/di_test.TestProgram_Run.func12.1
	ArgumentIndexes: [0 1]
	ResultIndexes: []
	HasBody: true
	HookIndexes: []
	HasCleanup: false
Function[1]:
	Index: 1
	Name: github.com/go-tk/di_test.TestProgram_Run.func12.2
	ArgumentIndexes: []
	ResultIndexes: [0 1]
	HasBody: true
	HookIndexes: []
	HasCleanup: false
Function[2]:
	Index: 2
	Name: github.com/go-tk/di_test.TestProgram_Run.func12.3
	ArgumentIndexes: []
	ResultIndexes: []
	HasBody: true
	HookIndexes: [0 1]
	HasCleanup: false
Argument[0]:
	FunctionIndex: 0
	ValueRef: x
	HasValueReceiver: true
	IsOptional: false
	ResultIndex: 0
	ReceiveValueAddr: false
Argument[1]:
	FunctionIndex: 0
	ValueRef: y
	HasValueReceiver: true
	IsOptional: false
	ResultIndex: 1
	ReceiveValueAddr: false
Result[0]:
	FunctionIndex: 1
	ValueName: x
	HasValue: true
	HookIndexes: [0]
Result[1]:
	FunctionIndex: 1
	ValueName: y
	HasValue: true
	HookIndexes: [1]
Hook[0]:
	FunctionIndex: 2
	ValueRef: x
	HasValueReceiver: true
	HasCallback: true
	ReceiveValueAddr: false
Hook[1]:
	FunctionIndex: 2
	ValueRef: y
	HasValueReceiver: true
	HasCallback: true
	ReceiveValueAddr: true
SortedFunctionIndexes: [2 1 0]
CalledFunctionCount: 3
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		var cancel context.CancelFunc
		*c.ctx, cancel = context.WithTimeout(context.Background(), 0)
		_ = cancel
		func() {
			err := c.p.NewFunction(Body(func(ctx context.Context) error {
				<-ctx.Done()
				return ctx.Err()
			}))
			if err != nil {
				t.Fatal(err)
			}
		}()
		c.err = context.DeadlineExceeded
		c.errStr = `call function; functionName="github.com/go-tk/di_test.TestProgram_Run.func13.1": ` + c.err.Error()
		c.repr = `
Function[0]:
	Index: 0
	Name: github.com/go-tk/di_test.TestProgram_Run.func13.1
	ArgumentIndexes: []
	ResultIndexes: []
	HasBody: true
	HookIndexes: []
	HasCleanup: false
SortedFunctionIndexes: [0]
CalledFunctionCount: 0
`[1:]
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		var cancel context.CancelFunc
		*c.ctx, cancel = context.WithCancel(context.Background())
		cancel()
		func() {
			var x int
			err := c.p.NewFunction(Result("x", &x), Body(func(context.Context) error {
				x = 100
				return nil
			}))
			if err != nil {
				t.Fatal(err)
			}
		}()
		func() {
			var x int
			err := c.p.NewFunction(Body(func(context.Context) error { return nil }),
				Hook("x", &x, func(context.Context) error { return context.Canceled }))
			if err != nil {
				t.Fatal(err)
			}
		}()
		c.err = context.Canceled
		c.errStr = `do callback; functionName="github.com/go-tk/di_test.TestProgram_Run.func14.2" valueRef="x": ` + c.err.Error()
		c.repr = `
Function[0]:
	Index: 0
	Name: github.com/go-tk/di_test.TestProgram_Run.func14.1
	ArgumentIndexes: []
	ResultIndexes: [0]
	HasBody: true
	HookIndexes: []
	HasCleanup: false
Function[1]:
	Index: 1
	Name: github.com/go-tk/di_test.TestProgram_Run.func14.2
	ArgumentIndexes: []
	ResultIndexes: []
	HasBody: true
	HookIndexes: [0]
	HasCleanup: false
Result[0]:
	FunctionIndex: 0
	ValueName: x
	HasValue: true
	HookIndexes: [0]
Hook[0]:
	FunctionIndex: 1
	ValueRef: x
	HasValueReceiver: true
	HasCallback: true
	ReceiveValueAddr: false
SortedFunctionIndexes: [1 0]
CalledFunctionCount: 2
`[1:]
	}).RunParallel(t)
}

func TestProgram_MustRun(t *testing.T) {
	{
		var p Program
		p.MustRun(context.Background())
	}
	assert.PanicsWithValue(t, `run program: di: value not found; valueRef="x" functionName="github.com/go-tk/di_test.TestProgram_MustRun.func1"`, func() {
		var p Program
		var x int
		p.MustNewFunction(Argument("x", &x), Body(func(context.Context) error { return nil }))
		p.MustRun(context.Background())
	})
}

func TestProgram_Clean(t *testing.T) {
	type C struct {
		p   *Program
		ctx *context.Context

		errStr string
		err    error
		seq    string
	}
	tc := testcase.New(func(t *testing.T, c *C) {
		var p Program
		ctx := context.Background()
		c.p = &p
		c.ctx = &ctx

		testcase.Callback(t, "0")

		err := p.Run(ctx)
		if c.errStr == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, c.errStr)
			if c.err != nil {
				assert.ErrorIs(t, err, c.err)
			}
		}
		if t.Failed() {
			t.FailNow()
		}
		p.Clean()

		testcase.Callback(t, "1")
	})

	tc.WithCallback("0", func(t *testing.T, c *C) {
		func() {
			var (
				x int
				y int
			)
			err := c.p.NewFunction(
				Argument("y", &y),
				Result("x", &x),
				Body(func(context.Context) error { c.seq += "A"; return nil }),
				Cleanup(func() { c.seq += "B" }))
			if err != nil {
				t.Fatal(err)
			}
		}()
		func() {
			var y int
			err := c.p.NewFunction(
				Result("y", &y),
				Body(func(context.Context) error { c.seq += "C"; return nil }),
				Cleanup(func() { c.seq += "D" }))
			if err != nil {
				t.Fatal(err)
			}
		}()
		func() {
			var y int
			err := c.p.NewFunction(
				Body(func(context.Context) error { c.seq += "E"; return nil }),
				Hook("y", &y, func(context.Context) error { c.seq += "F"; return nil }),
				Cleanup(func() { c.seq += "G" }))
			if err != nil {
				t.Fatal(err)
			}
		}()
	}).WithCallback("1", func(t *testing.T, c *C) {
		assert.Equal(t, c.seq, "ECFABDG")
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		func() {
			var (
				x int
				y int
			)
			err := c.p.NewFunction(
				Argument("y", &y),
				Result("x", &x),
				Body(func(context.Context) error { c.seq += "A"; return nil }),
				Cleanup(func() { c.seq += "B" }))
			if err != nil {
				t.Fatal(err)
			}
		}()
		func() {
			var y int
			err := c.p.NewFunction(
				Result("y", &y),
				Body(func(context.Context) error { c.seq += "C"; return nil }),
				Cleanup(func() { c.seq += "D" }))
			if err != nil {
				t.Fatal(err)
			}
		}()
		func() {
			var y int
			err := c.p.NewFunction(
				Body(func(context.Context) error { c.seq += "E"; return nil }),
				Hook("y", &y, func(context.Context) error { c.seq += "F"; return context.Canceled }),
				Cleanup(func() { c.seq += "G" }))
			if err != nil {
				t.Fatal(err)
			}
		}()
		c.errStr = `do callback; functionName="github.com/go-tk/di_test.TestProgram_Clean.func4.3" valueRef="y": ` + context.Canceled.Error()
	}).WithCallback("1", func(t *testing.T, c *C) {
		assert.Equal(t, c.seq, "ECFDG")
	}).RunParallel(t)

	tc.WithCallback("0", func(t *testing.T, c *C) {
		func() {
			var (
				x int
				y int
			)
			err := c.p.NewFunction(
				Argument("y", &y),
				Result("x", &x),
				Body(func(context.Context) error { c.seq += "A"; return context.Canceled }),
				Cleanup(func() { c.seq += "B" }))
			if err != nil {
				t.Fatal(err)
			}
		}()
		func() {
			var y int
			err := c.p.NewFunction(
				Result("y", &y),
				Body(func(context.Context) error { c.seq += "C"; return nil }),
				Cleanup(func() { c.seq += "D" }))
			if err != nil {
				t.Fatal(err)
			}
		}()
		func() {
			var y int
			err := c.p.NewFunction(
				Body(func(context.Context) error { c.seq += "E"; return nil }),
				Hook("y", &y, func(context.Context) error { c.seq += "F"; return nil }),
				Cleanup(func() { c.seq += "G" }))
			if err != nil {
				t.Fatal(err)
			}
		}()
		c.errStr = `call function; functionName="github.com/go-tk/di_test.TestProgram_Clean.func6.1": ` + context.Canceled.Error()
	}).WithCallback("1", func(t *testing.T, c *C) {
		assert.Equal(t, c.seq, "ECFADG")
	}).RunParallel(t)
}
