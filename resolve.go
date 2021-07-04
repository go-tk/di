package di

import (
	"bytes"
	"errors"
	"fmt"
)

type resolution12 struct {
	FunctionDescs []functionDesc

	ResultDescs map[string]*resultDesc
}

func (r *resolution12) ExecutePhase1() error {
	for i := range r.FunctionDescs {
		functionDesc := &r.FunctionDescs[i]
		for j := range functionDesc.Results {
			resultDesc := &functionDesc.Results[j]
			if err := r.processResultDesc(resultDesc, functionDesc); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *resolution12) processResultDesc(resultDesc1 *resultDesc, functionDesc *functionDesc) error {
	resultDescs := r.ResultDescs
	if resultDesc2, ok := resultDescs[resultDesc1.ValueID]; ok {
		tag := r.FunctionDescs[resultDesc2.FunctionIndex].Tag
		return fmt.Errorf("%w; tag1=%q tag2=%q valueID=%q",
			ErrValueAlreadyExists, functionDesc.Tag, tag, resultDesc1.ValueID)
	}
	if resultDescs == nil {
		resultDescs = make(map[string]*resultDesc)
		r.ResultDescs = resultDescs
	}
	resultDesc1.FunctionIndex = functionDesc.Index
	resultDescs[resultDesc1.ValueID] = resultDesc1
	return nil
}

func (r *resolution12) ExecutePhase2() error {
	for i := range r.FunctionDescs {
		functionDesc := &r.FunctionDescs[i]
		for j := range functionDesc.Arguments {
			argumentDesc := &functionDesc.Arguments[j]
			if err := r.processArgumentDesc(argumentDesc, functionDesc); err != nil {
				return err
			}
		}
		for j := range functionDesc.Hooks {
			hookDesc := &functionDesc.Hooks[j]
			if err := r.processHookDesc(hookDesc, functionDesc); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *resolution12) processArgumentDesc(argumentDesc *argumentDesc, functionDesc *functionDesc) error {
	resultDesc, ok := r.ResultDescs[argumentDesc.ValueID]
	if !ok {
		if !argumentDesc.IsOptional {
			return fmt.Errorf("%w; tag=%q valueID=%q",
				ErrValueNotFound, functionDesc.Tag, argumentDesc.ValueID)
		}
		return nil
	}
	valueType1 := argumentDesc.Value.Type()
	valueType2 := resultDesc.Value.Type()
	if valueType1 != valueType2 {
		tag := r.FunctionDescs[resultDesc.FunctionIndex].Tag
		return fmt.Errorf("%w; tag1=%q tag2=%q valueID=%q valueType1=%v valueType2=%v",
			ErrValueTypeMismatch, functionDesc.Tag, tag, argumentDesc.ValueID,
			valueType1, valueType2)
	}
	argumentDesc.Result = resultDesc
	return nil
}

func (r *resolution12) processHookDesc(hookDesc *hookDesc, functionDesc *functionDesc) error {
	resultDesc, ok := r.ResultDescs[hookDesc.ValueID]
	if !ok {
		return fmt.Errorf("%w; tag=%q valueID=%q",
			ErrValueNotFound, functionDesc.Tag, hookDesc.ValueID)
	}
	valueType1 := hookDesc.Value.Type()
	valueType2 := resultDesc.Value.Type()
	if valueType1 != valueType2 {
		tag := r.FunctionDescs[resultDesc.FunctionIndex].Tag
		return fmt.Errorf("%w; tag1=%q tag2=%q valueID=%q valueType1=%v valueType2=%v",
			ErrValueTypeMismatch, functionDesc.Tag, tag, hookDesc.ValueID,
			valueType1, valueType2)
	}
	hookDesc.FunctionIndex = functionDesc.Index
	resultDesc.Hooks = append(resultDesc.Hooks, hookDesc)
	return nil
}

type resolution3 struct {
	FunctionDescs []functionDesc

	OrderedFunctionDescIndexes []int

	stackTrace stackTrace
}

func (r *resolution3) ExecutePhase3() error {
	r.OrderedFunctionDescIndexes = make([]int, 0, len(r.FunctionDescs))
	for i := range r.FunctionDescs {
		if err := r.processFunctionDesc(i); err != nil {
			return err
		}
	}
	return nil
}

func (r *resolution3) processFunctionDesc(functionDescIndex int) error {
	functionDesc := &r.FunctionDescs[functionDescIndex]
	if functionDesc.Index < 0 { // functionDesc has already been processed
		return nil
	}
	if err := r.stackTrace.PushEntry(functionDesc); err != nil {
		return err
	}
	for i := range functionDesc.Arguments {
		argumentDesc := &functionDesc.Arguments[i]
		if argumentDesc.Result == nil {
			continue
		}
		r.stackTrace.ReferValue("argument", argumentDesc.ValueID)
		if err := r.processFunctionDesc(argumentDesc.Result.FunctionIndex); err != nil {
			return err
		}
	}
	for i := range functionDesc.Results {
		resultDesc := &functionDesc.Results[i]
		for _, hookDesc := range resultDesc.Hooks {
			r.stackTrace.ReferValue("hook", hookDesc.ValueID)
			if err := r.processFunctionDesc(hookDesc.FunctionIndex); err != nil {
				return err
			}
		}
	}
	r.stackTrace.PopEntry()
	r.OrderedFunctionDescIndexes = append(r.OrderedFunctionDescIndexes, functionDesc.Index)
	functionDesc.Index = -1 // mark functionDesc as processed
	return nil
}

type stackTrace struct {
	entries             []stackTraceEntry
	functionDescIndexes map[int]struct{}
}

func (st *stackTrace) PushEntry(functionDesc *functionDesc) error {
	st.entries = append(st.entries, stackTraceEntry{
		FunctionDesc: functionDesc,
	})
	functionDescIndexes := st.functionDescIndexes
	if _, ok := functionDescIndexes[functionDesc.Index]; ok {
		return fmt.Errorf("%w; %s", ErrCircularDependencies, st.dump())
	}
	if functionDescIndexes == nil {
		functionDescIndexes = make(map[int]struct{})
		st.functionDescIndexes = functionDescIndexes
	}
	functionDescIndexes[functionDesc.Index] = struct{}{}
	return nil
}

func (st *stackTrace) ReferValue(valueReferer string, valueID string) {
	entry := &st.entries[len(st.entries)-1]
	entry.ValueReferer = valueReferer
	entry.ValueID = valueID
}

func (st stackTrace) dump() string {
	var buffer bytes.Buffer
	i := len(st.entries) - 1
	for _, entry := range st.entries[:i] {
		buffer.WriteString(fmt.Sprintf("{tag: %q, %s: %q}", entry.FunctionDesc.Tag, entry.ValueReferer, entry.ValueID))
		buffer.WriteString(" => ")
	}
	entry := st.entries[i]
	buffer.WriteString(fmt.Sprintf("{tag: %q}", entry.FunctionDesc.Tag))
	return buffer.String()
}

func (st *stackTrace) PopEntry() {
	i := len(st.entries) - 1
	functionDescIndex := st.entries[i].FunctionDesc.Index
	st.entries = st.entries[:i]
	delete(st.functionDescIndexes, functionDescIndex)
}

type stackTraceEntry struct {
	FunctionDesc *functionDesc
	ValueReferer string
	ValueID      string
}

var (
	ErrValueAlreadyExists   = errors.New("di: value already exists")
	ErrValueNotFound        = errors.New("di: value not found")
	ErrValueTypeMismatch    = errors.New("di: value type mismatch")
	ErrCircularDependencies = errors.New("di: circular dependencies")
)
