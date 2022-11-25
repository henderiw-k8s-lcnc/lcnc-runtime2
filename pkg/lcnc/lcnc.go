package lcnc

import (
	"fmt"

	"github.com/henderiw-k8s-lcnc/lcnc-runtime2/pkg/dag"
	"github.com/henderiw-k8s-lcnc/lcnc-runtime2/pkg/lcncsyntax"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// TODO
// add output to the graph
// handle unresolved edges
// deal with $VALUE/$KEY

type Lcnc interface {
	GetExternalResources() ([]schema.GroupVersionResource, error)
	Transform() error
}

func New(cfg *lcncsyntax.LcncConfig) (Lcnc, error) {
	rootVertexName, err := cfg.GetRootVertexName()
	if err != nil {
		return nil, err
	}
	return &lcnc{
		lcncCfg:        cfg,
		d:              dag.NewDAG(),
		rootVertexName: rootVertexName,
	}, nil
}

type lcnc struct {
	lcncCfg        *lcncsyntax.LcncConfig
	d              dag.DAG
	rootVertexName string
}

func (r *lcnc) GetExternalResources() ([]schema.GroupVersionResource, error) {
	return r.lcncCfg.GetExternalResources()
}

func (r *lcnc) Transform() error {
	// process for
	for vertexName, v := range r.lcncCfg.For {
		if err := r.d.AddVertex(vertexName, v); err != nil {
			return err
		}
		if err := r.d.AddEdge("root", vertexName); err != nil {
			return err
		}
	}

	// process vars
	unresolvedVars := map[string]interface{}{}
	for vertexName, v := range r.lcncCfg.Vars {
		// debug -> to be removed
		if v.LcncQuery.Query != nil {
			fmt.Printf("var: %s query: %s \n", vertexName, *v.LcncQuery.Query)
		}
		// debug -> to be removed
		if v.For.Range != nil {
			fmt.Printf("var: %s forrange: %s \n", vertexName, *v.For.Range)
		}

		if v.LcncQuery.Query != nil {
			unresolved, err := r.addEdgeFromConfig(vertexName, *v.LcncQuery.Query, true, false)
			if err != nil {
				return err
			}
			if unresolved {
				unresolvedVars[vertexName] = v
				continue
			}
		}
		if v.For.Range != nil {
			unresolved, err := r.addEdgeFromConfig(vertexName, *v.For.Range, true, false)
			if err != nil {
				return err
			}
			if unresolved {
				unresolvedVars[vertexName] = v
				continue
			}
			if v.For.Map != nil {
				if v.For.Map.Value.Query != nil {
					unresolved, err := r.addEdgeFromConfig(vertexName, *v.For.Map.Value.Query, true, false)
					if err != nil {
						return err
					}
					if unresolved {
						unresolvedVars[vertexName] = v
						continue
					}
				}
			}
		}
		if _, ok := unresolvedVars[vertexName]; !ok {
			if err := r.d.AddVertex(vertexName, v); err != nil {
				return err
			}
		}
	}
	fmt.Printf("unresolved vars: %v \n", unresolvedVars)
	// resources, fucntions
	for vertexName, v := range r.lcncCfg.Resources {
		if v.For.Range != nil {
			unresolved, err := r.addEdgeFromConfig(vertexName, *v.For.Range, false, true)
			if err != nil {
				return err
			}
			if unresolved {
				unresolvedVars[vertexName] = v
			}
		}
		for localVarVertexName, v := range v.Function.Vars {
			// add the local var vertex name for resolution purposes
			if err := r.d.AddVertex(localVarVertexName, v); err != nil {
				return err
			}
			if v.LcncQuery.Query != nil {
				unresolved, err := r.addEdgeFromConfig(vertexName, *v.LcncQuery.Query, false, true)
				if err != nil {
					return err
				}
				if unresolved {
					unresolvedVars[vertexName] = v
				}
			}
		}
		for _, v := range v.Function.Input {
			unresolved, err := r.addEdgeFromConfig(vertexName, v, false, true)
			if err != nil {
				return err
			}
			if unresolved {
				unresolvedVars[vertexName] = v
			}
		}
		for k, v := range v.Function.Output {
			// for output we generate a new vertex and link ot to the root
			if err := r.d.AddVertex(k, v); err != nil {
				return err
			}
			if err := r.d.AddEdge(vertexName, k); err != nil {
				return err
			}
		}

	}
	fmt.Printf("unresolved resources/functions: %v \n", unresolvedVars)
	r.d.Walk("root")
	return nil
}

func (r *lcnc) addEdgeFromConfig(vertexName, s string, variable bool, function bool) (bool, error) {
	value, err := lcncsyntax.GetValue(s)
	if err != nil {
		return false, err
	}
	fmt.Printf("addEdgeFromConfig vertexName: %s, value: %s, kind: %s, variable: %v\n", vertexName, s, value.Kind, value.Variable)
	switch value.Kind {
	case lcncsyntax.GVRKind:
		// this would only be used for variables
		if !variable {
			return false, fmt.Errorf("cannot use gvr encoding in lcnc syntax other than variables and cluster applied output")
		}
		if err := r.d.AddEdge(r.rootVertexName, vertexName); err != nil {
			return false, fmt.Errorf("cannot add edge: from %s, to: %s, error: %s\n", r.rootVertexName, vertexName, err.Error())
		}

	case lcncsyntax.ChildVariableReferenceKind, lcncsyntax.RootVariableReferenceKind:
		// check of the reference exists
		if r.d.GetVertex(value.Variable[0]) {
			if err := r.d.AddEdge(value.Variable[0], vertexName); err != nil {
				return false, fmt.Errorf("cannot add edge: from %s, to: %s, error: %s\n", value.Variable[0], vertexName, err.Error())
			}
		} else {
			/* TODO Handle unresolved resolutions
			if function {
				return false, fmt.Errorf("a function cannot have an unresolved reference: got %s", value.Variable[0])
			}
			*/
			// TBD what to do if unresolved ????
			return true, nil
		}
	case lcncsyntax.KeyVariableReferenceKind:
		return false, nil
	default:
		return false, fmt.Errorf("cannot add edge: from %s, to: %s, error: %s\n", value.Variable[0], vertexName, err.Error())
	}
	return false, nil
}
