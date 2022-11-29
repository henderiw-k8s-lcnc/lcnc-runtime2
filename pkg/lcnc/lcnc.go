package lcnc

import (
	"fmt"
	"sync"

	"github.com/henderiw-k8s-lcnc/lcnc-runtime2/pkg/dag"
	"github.com/henderiw-k8s-lcnc/lcnc-runtime2/pkg/lcncsyntax"
)

// TODO
// add output to the graph
// handle unresolved edges
// deal with $VALUE/$KEY

type Lcnc interface {
	//GetExternalResources() ([]schema.GroupVersionResource, error)
	//Parse() (dag.DAG, error)
}

/*
func New(cfg *lcncsyntax.LcncConfig) (Lcnc, error) {
	rootVertexName, err := cfg.GetRootVertexName()
	if err != nil {
		return nil, err
	}
	return &lcnc{
		lcncCfg:        cfg,
		d:              dag.NewDAG(),
		rootVertexName: rootVertexName,
		output:         map[string]string{},
	}, nil
}
*/

type lcnc struct {
	lcncCfg        *lcncsyntax.LcncConfig
	d              dag.DAG
	rootVertexName string
	// key is the local variable
	ml            sync.RWMutex
	localVariable map[string]interface{}
	// key is the outputKey
	mo     sync.RWMutex
	output map[string]string
}

/*
func (r *lcnc) Parse() (dag.DAG, error) {
	//if err := r.validate(); err != nil {
	//	return nil, err
	//}
	if err := r.resolve(); err != nil {
		return nil, err
	}
	if err := r.addDependencies(); err != nil {
		return nil, err
	}
	r.d.Walk("root")

	newdag := r.d.TransitiveReduction()
	newdag.Walk("root")

	return newdag, nil
}
*/

/*
func (r *lcnc) validate() error {
	for vertexName, v := range r.lcncCfg.For {
		if err := r.validateContext(&OriginContext{Origin: OriginFor}, vertexName, v.Gvr); err != nil {
			return err
		}
	}
	for vertexName, v := range r.lcncCfg.Vars {
		if err := r.validateVar(vertexName, v); err != nil {
			return err
		}
	}
	for vertexName, v := range r.lcncCfg.Resources {
		if err := r.validateResource(vertexName, v); err != nil {
			return err
		}
	}
	return nil
}

func (r *lcnc) validateVar(vertexName string, v lcncsyntax.LcncVar) error {
	if v.LcncQuery.Query != nil {
		if err := r.validateContext(&OriginContext{Origin: OriginVariable, Query: true}, vertexName, *v.LcncQuery.Query); err != nil {
			return err
		}
	}
	if v.For.Range != nil {
		if err := r.validateContext(&OriginContext{Origin: OriginVariable, ForLoop: true}, vertexName, *v.For.Range); err != nil {
			return err
		}
		if v.For.Map != nil && v.For.Map.Value.Query != nil {
			if err := r.validateContext(&OriginContext{Origin: OriginVariable, ForLoop: true}, vertexName, *v.For.Range); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *lcnc) validateResource(vertexName string, v lcncsyntax.LcncResource) error {
	forloop := false
	if v.For.Range != nil {
		forloop = true
		if err := r.validateContext(&OriginContext{Origin: Originfunction, ForLoop: forloop}, vertexName, *v.For.Range); err != nil {
			return err
		}
	}
	for _, v := range v.Function.Vars {
		if err := r.validateContext(&OriginContext{Origin: Originfunction, ForLoop: forloop, Query: true}, vertexName, *v.Query); err != nil {
			return err
		}
	}
	for _, v := range v.Function.Input {
		if err := r.validateContext(&OriginContext{Origin: Originfunction, ForLoop: forloop, Input: true}, vertexName, v); err != nil {
			return err
		}
	}
	for _, v := range v.Function.Output {
		if err := r.validateContext(&OriginContext{Origin: Originfunction, ForLoop: forloop, Output: true}, vertexName, v); err != nil {
			return err
		}
	}
	return nil
}
*/

/*
func (r *lcnc) resolve() error {
	if err := r.d.AddVertex("root", nil); err != nil {
		return err
	}
	for vertexName, v := range r.lcncCfg.For {
		fmt.Printf("for: %s gvr: %s \n", vertexName, v.Gvr)
		if err := r.d.AddVertex(vertexName, v); err != nil {
			return err
		}
	}
	// we add the variables in the dag vertex list already
	// but we still check the resolution to see if they resolve or not
	for vertexName, v := range r.lcncCfg.Vars {
		if err := r.d.AddVertex(vertexName, v); err != nil {
			return err
		}
	}
	// copy the vars, so we keep the original variables
	unresolvedVars := r.lcncCfg.DeepCopyVariables()
	// if not all defined veriables resolve we fail
	if err := r.resolveUnresolvedVars(unresolvedVars); err != nil {
		return err
	}

	fmt.Println("all variables resolved")

	// we add the Resources in the dag vertex list already
	// but we still check the resolution to see if they resolve or not
	// Also the output is added in the dag vertex list
	for vertexName, v := range r.lcncCfg.Resources {
		if err := r.d.AddVertex(vertexName, v); err != nil {
			return err
		}
		for k := range v.Function.Output {
			// for output we generate a new output mapping
			if err := r.AddOutputMapping(k, vertexName); err != nil {
				return err
			}
		}
	}
	unresolvedResources := r.lcncCfg.DeepCopyResources()
	// if not all defined veriables resolve we fail
	if err := r.resolveUnresolvedResources(unresolvedResources); err != nil {
		return err
	}
	fmt.Println("all resources resolved")
	r.PrintOutputMappings()
	return nil
}
*/

/*
func (r *lcnc) resolveUnresolvedVars(unresolved map[string]lcncsyntax.LcncVar) error {
	totalUnresolved := len(unresolved)
	for vertexName, v := range unresolved {
		if r.resolveVariables(vertexName, v) {
			delete(unresolved, vertexName)
		}
	}
	// when the new unresolved is 0 we are done and all variabled are resolved
	newUnresolved := len(unresolved)
	if newUnresolved == 0 {
		return nil
	}
	if newUnresolved < totalUnresolved {
		// recursively continue resolution
		if err := r.resolveUnresolvedVars(unresolved); err != nil {
			return err
		}
	}
	return fmt.Errorf("not all variables could be resolved: %v", unresolved)
}
*/

/*
func (r *lcnc) resolveVariables(vertexName string, v lcncsyntax.LcncVar) bool {
	if v.LcncQuery.Query != nil {
		if !r.isResolved(&OriginContext{Origin: OriginVariable, Query: true}, vertexName, *v.LcncQuery.Query) {
			return false
		}
	}
	if v.For.Range != nil {
		if !r.isResolved(&OriginContext{Origin: OriginVariable, ForLoop: true}, vertexName, *v.For.Range) {
			return false
		}
		if v.For.Map != nil && v.For.Map.Value.Query != nil {
			if !r.isResolved(&OriginContext{Origin: OriginVariable, ForLoop: true}, vertexName, *v.For.Range) {

			}
		}
	}
	return true
}
*/

/*
func (r *lcnc) resolveUnresolvedResources(unresolved map[string]lcncsyntax.LcncResource) error {
	totalUnresolved := len(unresolved)
	for vertexName, v := range unresolved {
		if r.resolveResources(vertexName, v) {
			delete(unresolved, vertexName)
		}
	}
	// when the new unresolved is 0 we are done and all variabled are resolved
	newUnresolved := len(unresolved)
	if newUnresolved == 0 {
		return nil
	}
	if newUnresolved < totalUnresolved {
		// recursively continue resolution
		if err := r.resolveUnresolvedResources(unresolved); err != nil {
			return err
		}
	}
	return fmt.Errorf("not all resource could be resolved: %v", unresolved)
}
*/

/*
func (r *lcnc) resolveResources(vertexName string, v lcncsyntax.LcncResource) bool {
	r.initLocalVariables()
	forloop := false
	if v.For.Range != nil {
		forloop = true
		if !r.isResolved(&OriginContext{Origin: Originfunction, ForLoop: forloop}, vertexName, *v.For.Range) {
			return false
		}
	}
	for localVarName, v := range v.Function.Vars {
		// TODO how to handle this error better
		if err := r.AddLocalVariable(localVarName, v); err != nil {
			return false
		}
		if !r.isResolved(&OriginContext{Origin: Originfunction, ForLoop: forloop, Query: true}, vertexName, *v.Query) {
			return false
		}
	}
	for _, v := range v.Function.Input {
		if !r.isResolved(&OriginContext{Origin: Originfunction, ForLoop: forloop, Input: true}, vertexName, v) {
			return false
		}
	}
	// Output does not need to be resolved
	return true
}
*/

/*
func (r *lcnc) addDependencies() error {
	for vertexName := range r.lcncCfg.For {
		r.d.AddDownEdge("root", vertexName)
	}
	for vertexName, v := range r.lcncCfg.Vars {
		if err := r.addDependenciesVar(vertexName, v); err != nil {
			return err
		}
	}
	for vertexName, v := range r.lcncCfg.Resources {
		if err := r.addDependenciesResource(vertexName, v); err != nil {
			return err
		}
	}
	return nil
}
*/

/*
func (r *lcnc) addDependenciesVar(vertexName string, v lcncsyntax.LcncVar) error {
	if v.LcncQuery.Query != nil {
		if err := r.addEdge(vertexName, *v.LcncQuery.Query); err != nil {
			return err
		}
	}
	if v.For.Range != nil {
		if err := r.addEdge(vertexName, *v.For.Range); err != nil {
			return err
		}
		if v.For.Map != nil {
			if v.For.Map.Value.Query != nil {
				if err := r.addEdge(vertexName, *v.For.Map.Value.Query); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
*/

/*
func (r *lcnc) addDependenciesResource(vertexName string, v lcncsyntax.LcncResource) error {
	r.initLocalVariables()
	if v.For.Range != nil {
		// output resolves to the output dependency resource
		value, err := lcncsyntax.GetValue(*v.For.Range)
		if err != nil {
			return err
		}
		if r.HasOutputMapping(value.Variable[0]) {
			r.d.AddDownEdge(r.GetOutputMapping(value.Variable[0]), vertexName)
			r.d.AddUpEdge(vertexName, r.GetOutputMapping(value.Variable[0]))
		} else {
			if err := r.addEdge(vertexName, *v.For.Range); err != nil {
				return err
			}
		}
	}
	for localVarName, v := range v.Function.Vars {
		// TODO how to handle this error better
		if err := r.AddLocalVariable(localVarName, v); err != nil {
			return err
		}
		if v.LcncQuery.Query != nil {
			// output resolves to the output dependency resource
			value, err := lcncsyntax.GetValue(*v.LcncQuery.Query)
			if err != nil {
				return err
			}
			if r.HasOutputMapping(value.Variable[0]) {
				r.d.AddDownEdge(r.GetOutputMapping(value.Variable[0]), vertexName)
				r.d.AddUpEdge(vertexName, r.GetOutputMapping(value.Variable[0]))
			} else {
				if err := r.addEdge(vertexName, *v.LcncQuery.Query); err != nil {
					return err
				}
			}
		}
	}
	for _, v := range v.Function.Input {
		// dont add dependencies for local variables
		value, err := lcncsyntax.GetValue(v)
		if err != nil {
			return err
		}
		if r.GetLocalVariable(value.Variable[0]) {
			return nil
		}
		// output resolves to the output dependency resource
		if r.HasOutputMapping(value.Variable[0]) {
			r.d.AddDownEdge(r.GetOutputMapping(value.Variable[0]), vertexName)
			r.d.AddUpEdge(vertexName, r.GetOutputMapping(value.Variable[0]))
		} else {
			if err := r.addEdge(vertexName, v); err != nil {
				return err
			}
		}
	}

	//for k := range v.Function.Output {
	//	r.d.AddDownEdge(vertexName, k)
	//}
	return nil
}
*/

/*
func (r *lcnc) validateContext(o *OriginContext, ctxName, s string) error {
	value, err := lcncsyntax.GetValue(s)
	if err != nil {
		return err
	}
	fmt.Printf("validate ctxName: %s, value: %s, kind: %s, variable: %v\n", ctxName, s, value.Kind, value.Variable)
	switch value.Kind {
	case lcncsyntax.GVRKind:
		// only allowed for variables and output
		if o.Origin == Originfunction && !o.Output {
			return fmt.Errorf("cannot use gvr encoding syntax in resource/function statements, use variables instead")
		}
	case lcncsyntax.ChildVariableReferenceKind, lcncsyntax.RootVariableReferenceKind:
		if o.Origin == OriginFor {
			return fmt.Errorf("can only use GVR resources in for statements")
		}
	case lcncsyntax.KeyVariableReferenceKind:
		if o.Origin == OriginFor {
			return fmt.Errorf("can only use GVR resources in for statements")
		}
		if !o.ForLoop {
			return fmt.Errorf("cannot use Key variabales without a for statement")
		}
	case lcncsyntax.VariableKind:
		if o.Origin == OriginFor {
			return fmt.Errorf("can only use GVR resources in for statements")
		}
	default:
		return fmt.Errorf("unknown variable, got: %s", s)
	}
	return nil
}
*/

/*
func (r *lcnc) isResolved(o *OriginContext, vertexName, s string) bool {
	// we dont handle the validation here, since is handled before
	value, err := lcncsyntax.GetValue(s)
	if err != nil {
		// this should never happen since validation should have happened before
		return false
	}
	resolved := false
	switch value.Kind {
	case lcncsyntax.GVRKind:
		// resolution is global, so the only resolution we can validate is if the resource exists
		// on the api server
		resolved = true
		break
	case lcncsyntax.ChildVariableReferenceKind, lcncsyntax.RootVariableReferenceKind:
		// input of a function can resolve to a local variable
		// if so we should be ok and dont have to add an edge since the variable has already been
		// resolved to handle the dependency
		if o.Origin == Originfunction && o.Input && r.GetLocalVariable(value.Variable[0]) {
			resolved = true
			break
		}
		// a fucntion can be dependent on another fn based on the output
		if o.Origin == Originfunction && r.HasOutputMapping(value.Variable[0]) {
			resolved = true
			break
		}
		// check if the global variable/output exists
		if r.d.GetVertex(value.Variable[0]) {
			resolved = true
			break
		}
	case lcncsyntax.KeyVariableReferenceKind:
		if o.ForLoop {
			resolved = true
			break
		}
	}
	fmt.Printf("isResolved vertexName: %s, value: %s, kind: %s, variable: %v, resolved: %t\n", vertexName, s, value.Kind, value.Variable, resolved)
	return resolved
}
*/

func (r *lcnc) addEdge(vertexName, s string) error {
	value, err := lcncsyntax.GetValue(s)
	if err != nil {
		return err
	}
	switch value.Kind {
	case lcncsyntax.GVRKind:
		r.d.AddDownEdge(r.rootVertexName, vertexName)
		r.d.AddUpEdge(vertexName, r.rootVertexName)
	case lcncsyntax.ChildVariableReferenceKind, lcncsyntax.RootVariableReferenceKind:
		r.d.AddDownEdge(value.Variable[0], vertexName)
		r.d.AddUpEdge(vertexName, value.Variable[0])
	case lcncsyntax.KeyVariableReferenceKind:
	default:
		return fmt.Errorf("cannot add edge: from %s, to: %s\n", value.Variable[0], vertexName)
	}
	return nil
}
