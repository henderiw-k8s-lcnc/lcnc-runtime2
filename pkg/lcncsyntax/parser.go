package lcncsyntax

import (
	"github.com/henderiw-k8s-lcnc/lcnc-runtime2/pkg/dag"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type Parser interface {
	GetExternalResources() ([]schema.GroupVersionResource, []Result)
	Parse() (dag.DAG, string, []Result)
}

func NewParser(cfg *LcncConfig) (Parser, []Result) {
	p := &lcncparser{
		lcncCfg: cfg,
		//d:       dag.NewDAG(),
		//output: map[string]string{},
	}
	// add the callback function to record validation results results
	result := p.ValidateSyntax()
	p.rootVertexName = cfg.GetRootVertexName()

	return p, result
}

type lcncparser struct {
	lcncCfg *LcncConfig
	//d              dag.DAG
	rootVertexName string
}

func (r *lcncparser) Parse() (dag.DAG, string, []Result) {
	// validate the config when creating the dag
	d := dag.NewDAG()
	// resolves the dependencies in the dag
	// step1. check if all dependencies resolve
	// step2. add the dependencies in the dag
	result := r.Resolve(d)
	if len(result) != 0 {
		return nil, "", result
	}
	//d.GetDependencyMap(r.rootVertexName)
	// optimizes the dependncy graph based on transit reduction
	// techniques
	d.TransitiveReduction()
	//d.GetDependencyMap(r.rootVertexName)
	return d, r.rootVertexName, nil
}
