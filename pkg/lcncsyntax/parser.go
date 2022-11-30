package lcncsyntax

import (
	"sync"

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
		output: map[string]string{},
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
	// localVariable is used to store local variables of a function
	// key is the local variable
	ml            sync.RWMutex
	localVariable map[string]interface{}
	// output is used to store output to function mapping for
	// lookup resolution
	// key is the outputKey
	mo     sync.RWMutex
	output map[string]string

	syntaxValidationResult []Result
}

func (r *lcncparser) Parse() (dag.DAG, string, []Result) {
	// validate
	d := dag.NewDAG()
	result := r.Resolve(d)
	if len(result) != 0 {
		return nil, "", result
	}
	newd := d.TransitiveReduction()
	// transitive reduction
	return newd, r.rootVertexName, nil
}
