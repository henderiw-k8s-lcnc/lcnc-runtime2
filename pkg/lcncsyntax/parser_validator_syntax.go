package lcncsyntax

import (
	"fmt"
	"sync"
)

func (r *lcncparser) ValidateSyntax() []Result {
	vs := &vs{
		result: []Result{},
	}

	fnc := WalkConfig{
		lcncCfgPreHookFn: vs.validateLcncCfgPreHookFn,
		lcncGvrObjectFn:  vs.validateLcncGvrObjectFn,
		lcncBlockFn:      vs.validateBlockFn,
		lcncVarFn:        vs.validateVarFn,
		lcncFunctionFn:   vs.validateFunctionFn,
		lcncServiceFn:    vs.validateServiceFn,
	}

	// walk the config to validate the syntax
	r.walkLcncConfig(fnc)
	return vs.result

}

type vs struct {
	mr     sync.RWMutex
	result []Result
}

func (r *vs) recordResult(result Result) {
	r.mr.Lock()
	defer r.mr.Unlock()
	r.result = append(r.result, result)
}

func (r *vs) validateLcncCfgPreHookFn(lcncCfg *LcncConfig) {
	if len(lcncCfg.For) != 1 {
		r.recordResult(Result{
			Origin: OriginFor,
			Index:  0,
			Error:  fmt.Errorf("lcnc config must have just 1 for statement, got: %v", lcncCfg.For).Error(),
		})
	}
}

func (r *vs) validateLcncGvrObjectFn(o Origin, idx int, vertexName string, v LcncGvrObject) {
	r.validateContext(&OriginContext{Origin: o, Index: idx, VertexName: vertexName}, v.Gvr)
}

// TODO add if statement ??
func (r *vs) validateBlockFn(o Origin, idx int, v LcncBlock) {
	if v.For == nil {
		r.recordResult(Result{
			Origin: o,
			Index:  idx,
			Error:  fmt.Errorf("for must be present in block").Error(),
		})
		return
	}
	if v.For.Range == nil {
		r.recordResult(Result{
			Origin: o,
			Index:  idx,
			Error:  fmt.Errorf("range must be present in for block").Error(),
		})
		return
	}
	r.validateContext(&OriginContext{Origin: o, Index: 0, VertexName: ""}, *v.For.Range)
}

func (r *vs) validateVarFn(o Origin, block bool, idx int, vertexName string, v LcncVar) {
	if v.Slice == nil && v.Map == nil {
		r.recordResult(Result{
			Origin: o,
			Index:  idx,
			Name:   vertexName,
			Error:  fmt.Errorf("slice or map must be present in a variable").Error(),
		})
		return
	}
	if v.Slice != nil {
		r.validateValue(o, block, idx, vertexName, v.Slice.LcncValue)
	}
	if v.Map != nil {
		r.validateValue(o, block, idx, vertexName, v.Map.LcncValue)
		// validate key
		input := false
		if o == OriginVariable {
			input = true
		}
		if v.Map.Key == nil {
			r.recordResult(Result{
				Origin: o,
				Index:  idx,
				Name:   vertexName,
				Error:  fmt.Errorf("key must be present in a map variable").Error(),
			})
			return
		}
		if v.Map.Key != nil {
			r.validateContext(&OriginContext{
				Origin:     o,
				Index:      idx,
				VertexName: vertexName,
				ForBlock:   block,
				Input:      input,
				Output:     true}, *v.Map.Key)
		}
	}
}

func (r *vs) validateValue(o Origin, block bool, idx int, vertexName string, v LcncValue) {
	if v.String == nil && v.LcncQuery.Query == nil {
		r.recordResult(Result{
			Origin: o,
			Index:  idx,
			Name:   vertexName,
			Error:  fmt.Errorf("query or string must be present in a slice variable").Error(),
		})
	}
	input := false
	if o == OriginVariable {
		input = true
	}
	if v.String != nil {
		r.validateContext(&OriginContext{
			Origin:     o,
			Index:      idx,
			VertexName: vertexName,
			ForBlock:   block,
			Input:      input,
			Output:     true}, *v.String)
	}
	if v.LcncQuery.Query != nil {
		r.validateContext(&OriginContext{
			Origin:     o,
			Index:      idx,
			VertexName: vertexName,
			ForBlock:   block,
			Input:      input,
			Query:      true,
			Output:     true}, *v.LcncQuery.Query)
	}
}

func (r *vs) validateFunctionFn(o Origin, block bool, idx int, vertexName string, v LcncFunction) {
	//v.ImageName -> must be present
	if v.LcncImage.ImageName == nil {
		r.recordResult(Result{
			Origin: OriginFunction,
			Index:  idx,
			Name:   vertexName,
			Error:  fmt.Errorf("imageName must be present in a function").Error(),
		})
	}
	// input must be present
	if len(v.Input) == 0 {
		r.recordResult(Result{
			Origin: OriginFunction,
			Index:  idx,
			Name:   vertexName,
			Error:  fmt.Errorf("input must be present in a function").Error(),
		})
	}
	// output must be present
	if len(v.Output) == 0 {
		r.recordResult(Result{
			Origin: OriginFunction,
			Index:  idx,
			Name:   vertexName,
			Error:  fmt.Errorf("output must be present in a function").Error(),
		})
	}
	// uniqueness of local variables are checked by the way we defined the api
	// map[string]LcncVar
	for _, v := range v.Vars {
		r.validateVarFn(OriginFunction, block, idx, vertexName, v)
	}
	for _, v := range v.Input {
		r.validateContext(&OriginContext{
			Origin:     OriginFunction,
			Index:      idx,
			VertexName: vertexName,
			ForBlock:   block,
			Input:      true}, v)
	}
	for _, v := range v.Output {
		r.validateContext(&OriginContext{
			Origin:     OriginFunction,
			Index:      idx,
			VertexName: vertexName,
			ForBlock:   block,
			Output:     true}, v)
	}
}

func (r *vs) validateServiceFn(o Origin, block bool, idx int, vertexName string, v LcncFunction) {}

func (r *vs) validateContext(o *OriginContext, s string) {
	value, err := GetValue(s)
	if err != nil {
		r.recordResult(Result{
			Origin: o.Origin,
			Index:  o.Index,
			Name:   o.VertexName,
			Error:  err.Error(),
		})
		return
	}
	//fmt.Printf("validate ctxName: %s, value: %s, kind: %s, variable: %v\n", o.VertexName, s, value.Kind, value.Variable)
	switch value.Kind {
	case GVRKind:
		// only allowed for variables and output
		if o.Origin == OriginFunction && !o.Output {
			r.recordResult(Result{
				Origin: o.Origin,
				Index:  o.Index,
				Name:   o.VertexName,
				Error:  fmt.Errorf("cannot use gvr encoding syntax in function statements, use variables instead").Error(),
			})
			return
		}
	case ChildVariableReferenceKind, RootVariableReferenceKind:
		if o.Origin == OriginFor {
			r.recordResult(Result{
				Origin: o.Origin,
				Index:  o.Index,
				Name:   o.VertexName,
				Error:  fmt.Errorf("can only use GVR resources in for statements").Error(),
			})
			return
		}
	case KeyVariableReferenceKind:
		if o.Origin == OriginFor {
			r.recordResult(Result{
				Origin: o.Origin,
				Index:  o.Index,
				Name:   o.VertexName,
				Error:  fmt.Errorf("can only use GVR resources in for statements").Error(),
			})
			return
		}
		if !o.ForBlock {
			r.recordResult(Result{
				Origin: o.Origin,
				Index:  o.Index,
				Name:   o.VertexName,
				Error:  fmt.Errorf("cannot use Key variabales without a for statement").Error(),
			})
			return
		}
	case VariableKind:
		if o.Origin == OriginFor {
			r.recordResult(Result{
				Origin: o.Origin,
				Index:  o.Index,
				Name:   o.VertexName,
				Error:  fmt.Errorf("can only use GVR resources in for statements").Error(),
			})
			return
		}
	default:
		r.recordResult(Result{
			Origin: o.Origin,
			Index:  o.Index,
			Name:   o.VertexName,
			Error:  fmt.Errorf("unknown variable, got: %s", s).Error(),
		})
		return
	}
}
