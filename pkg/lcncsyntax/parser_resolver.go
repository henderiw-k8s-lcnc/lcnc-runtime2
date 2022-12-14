package lcncsyntax

import (
	"fmt"
	"sync"

	"github.com/henderiw-k8s-lcnc/lcnc-runtime2/pkg/dag"
)

func (r *lcncparser) Resolve(d dag.DAG) []Result {
	rs := &rs{
		d:              d,
		result:         []Result{},
		localVariable:  map[string]interface{}{},
		output:         map[string]string{},
		rootVertexName: r.rootVertexName,
	}

	fnc := WalkConfig{
		lcncGvrObjectFn:         rs.resolveLcncGvrObjectFn,
		lcncVarFn:               rs.resolveLcncVarFn,
		lcncVarsPostHookFn:      rs.resolveLcncVarsPostHookFn,
		lcncFunctionFn:          rs.resolveLcncFunctionFn,
		lcncFunctionsPostHookFn: rs.resolveLcncFunctionsPostHookFn,
	}

	// walk the config to resolve the vars/functions/etc
	// we add the verteces in the graph and check for duplicate entries
	// we create an output mapping, which will be used in the 2nd step (edges/dependencies)
	// local variables are resolved within the function
	r.walkLcncConfig(fnc)
	// stop if errors were found
	if len(rs.result) != 0 {
		return rs.result
	}

	// the 2nd walk adds the dependencies and edges in the graph
	fnc = WalkConfig{
		lcncVarsPostHookFn:      rs.addDependenciesLcncVarsPostHookFn,
		lcncFunctionsPostHookFn: rs.addDependenciesLcncFunctionsPostHookFn,
	}
	r.walkLcncConfig(fnc)

	return rs.result
}

type rs struct {
	rootVertexName string
	d              dag.DAG
	mr             sync.RWMutex
	result         []Result
	// key is the local variable
	ml            sync.RWMutex
	localVariable map[string]interface{}
	// key is the outputKey
	mo     sync.RWMutex
	output map[string]string
}

func (r *rs) recordResult(result Result) {
	r.mr.Lock()
	defer r.mr.Unlock()
	r.result = append(r.result, result)
}

func (r *rs) resolveLcncGvrObjectFn(o Origin, idx int, vertexName string, v LcncGvrObject) {
	if o == OriginFor {
		if err := r.d.AddVertex(vertexName, v); err != nil {
			r.recordResult(Result{
				Origin: o,
				Index:  idx,
				Name:   vertexName,
				Error:  err.Error(),
			})
		}
	}
}

func (r *rs) resolveLcncVarFn(o Origin, block bool, idx int, vertexName string, v LcncVar) {
	if err := r.d.AddVertex(vertexName, v); err != nil {
		r.recordResult(Result{
			Origin: o,
			Index:  idx,
			Name:   vertexName,
			Error:  err.Error(),
		})
	}
}

func (r *rs) resolveLcncFunctionFn(o Origin, block bool, idx int, vertexName string, v LcncFunction) {
	if err := r.d.AddVertex(vertexName, v); err != nil {
		r.recordResult(Result{
			Origin: o,
			Index:  idx,
			Name:   vertexName,
			Error:  err.Error(),
		})
	}

	for k := range v.Output {
		// for output we generate a new output mapping
		if err := r.AddOutputMapping(k, vertexName); err != nil {
			r.recordResult(Result{
				Origin: o,
				Index:  idx,
				Name:   vertexName,
				Error:  err.Error(),
			})
		}
	}
}

func (r *rs) resolveLcncVarsPostHookFn(v []LcncVarBlock) {
	newvars := CopyVariables(v)
	unresolved := r.resolveUnresolvedVars(newvars)
	if len(unresolved) != 0 {
		for idxName := range unresolved {
			vertexName, idx := GetIdxName(idxName)
			r.recordResult(Result{
				Origin: OriginVariable,
				Index:  idx,
				Name:   vertexName,
				Error:  fmt.Errorf("unresolved variable").Error(),
			})
		}
	}
}

func (r *rs) resolveLcncFunctionsPostHookFn(v []LcncFunctionsBlock) {
	newfns := CopyFunctions(v)
	unresolved := r.resolveUnresolvedFunctions(newfns)
	if len(unresolved) != 0 {
		for idxName := range unresolved {
			vertexName, idx := GetIdxName(idxName)
			r.recordResult(Result{
				Origin: OriginFunction,
				Index:  idx,
				Name:   vertexName,
				Error:  fmt.Errorf("unresolved function").Error(),
			})
		}
	}
}

func (r *rs) resolveUnresolvedVars(unresolved map[string]LcncVarBlock) map[string]LcncVarBlock {
	totalUnresolved := len(unresolved)
	for idxName, v := range unresolved {
		if r.resolveVariable(v) {
			delete(unresolved, idxName)
		}
	}
	// when the new unresolved is 0 we are done and all variabled are resolved
	newUnresolved := len(unresolved)
	if newUnresolved == 0 {
		return unresolved
	}
	if newUnresolved < totalUnresolved {
		r.resolveUnresolvedVars(unresolved)
	}
	return unresolved
}

func (r *rs) resolveVariable(v LcncVarBlock) bool {
	for vertexName, vv := range v.LcncVariables {
		forblock := false
		if v.For != nil && v.For.Range != nil {
			forblock = true
			if !r.isResolved(&OriginContext{Origin: OriginVariable, ForBlock: true}, *v.For.Range) {
				fmt.Printf("unresolved vertexName: %s\n", vertexName)
				return false
			}
		}
		if vv.Map != nil {
			if !r.resolveMap(&OriginContext{Origin: OriginVariable, ForBlock: forblock, Query: true}, vv.Map) {
				return false
			}
		}
		if vv.Slice != nil {
			if !r.resolveValue(&OriginContext{Origin: OriginVariable, ForBlock: forblock, Query: true}, vv.Slice.LcncValue) {
				return false
			}
		}
	}
	return true
}

func (r *rs) resolveUnresolvedFunctions(unresolved map[string]LcncFunctionsBlock) map[string]LcncFunctionsBlock {
	totalUnresolved := len(unresolved)
	for idxName, v := range unresolved {
		if r.resolveFunction(v) {
			delete(unresolved, idxName)
		}
	}
	// when the new unresolved is 0 we are done and all variabled are resolved
	newUnresolved := len(unresolved)
	if newUnresolved == 0 {
		return unresolved
	}
	if newUnresolved < totalUnresolved {
		r.resolveUnresolvedFunctions(unresolved)
	}
	return unresolved
}

func (r *rs) resolveFunction(v LcncFunctionsBlock) bool {
	for vertexName, vv := range v.LcncFunctions {
		// initialize the local variables for local resolution
		r.initLocalVariables()
		forblock := false
		if v.For != nil && v.For.Range != nil {
			forblock = true
			if !r.isResolved(&OriginContext{Origin: OriginFunction, ForBlock: true}, *v.For.Range) {
				fmt.Printf("unresolved vertexName: %s\n", vertexName)
				return false
			}
		}
		for localVarName, v := range vv.Vars {
			// TODO how to handle this error better
			if err := r.AddLocalVariable(localVarName, v); err != nil {
				return false
			}
			if v.Map != nil {
				if !r.resolveMap(&OriginContext{Origin: OriginFunction, ForBlock: forblock, Query: true}, v.Map) {
					return false
				}
			}
			if v.Slice != nil {
				if !r.resolveValue(&OriginContext{Origin: OriginFunction, ForBlock: forblock, Query: true}, v.Slice.LcncValue) {
					return false
				}
			}
		}
		for _, v := range vv.Input {
			if !r.isResolved(&OriginContext{Origin: OriginFunction, ForBlock: forblock, Input: true}, v) {
				return false
			}
		}
	}
	return true
}

func (r *rs) resolveMap(o *OriginContext, v *LcncMap) bool {
	if v.Key != nil {
		if !r.isResolved(o, *v.Key) {
			return false
		}
	}
	if !r.resolveValue(o, v.LcncValue) {
		return false
	}
	return true
}

func (r *rs) resolveValue(o *OriginContext, v LcncValue) bool {
	if v.LcncQuery.Query != nil {
		if !r.isResolved(o, *v.LcncQuery.Query) {
			return false
		}
	}
	if v.String != nil {
		if !r.isResolved(o, *v.String) {
			return false
		}
	}
	return true
}

func (r *rs) isResolved(o *OriginContext, s string) bool {
	// we dont handle the validation here, since is handled before
	value, err := GetValue(s)
	if err != nil {
		// this should never happen since validation should have happened before
		return false
	}
	resolved := false
	switch value.Kind {
	case GVRKind:
		// resolution is global, so the only resolution we can validate is if the resource exists
		// on the api server
		resolved = true
	case ChildVariableReferenceKind, RootVariableReferenceKind:
		// input of a function can resolve to a local variable
		// if so we should be ok and dont have to add an edge since the variable has already been
		// resolved to handle the dependency
		if o.Origin == OriginFunction && o.Input && r.GetLocalVariable(value.Variable[0]) {
			resolved = true
			break
		}
		// a fucntion can be dependent on another fn based on the output
		if o.Origin == OriginFunction && r.HasOutputMapping(value.Variable[0]) {
			resolved = true
			break
		}
		// check if the global variable/output exists
		if r.d.GetVertex(value.Variable[0]) {
			resolved = true
			break
		}
	case KeyVariableReferenceKind:
		if o.ForBlock {
			resolved = true
			break
		}
	}
	//fmt.Printf("isResolved originContext: %v string: %s resolved: %t\n", *o, s, resolved)
	return resolved
}
