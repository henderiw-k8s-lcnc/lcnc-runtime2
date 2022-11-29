package lcncsyntax

import (
	"sync"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (r *lcncparser) GetExternalResources() ([]schema.GroupVersionResource, []Result) {
	er := &er{
		result:    []Result{},
		resources: []schema.GroupVersionResource{},
	}
	er.resultFn = er.recordResult
	er.addResourceFn = er.addResource

	fnc := WalkConfig{
		lcncCfgPreHookFn: nil,
		lcncGvrObjectFn:  er.validateLcncGvrObjectFn,
		lcncBlockFn:      er.validateBlockFn,
		lcncVarFn:        er.validateVarFn,
		lcncFunctionFn:   er.validateFunctionFn,
		lcncServiceFn:    nil,
	}

	// validate the external resources
	r.walkLcncConfig(fnc)
	return er.resources, er.result
}

type er struct {
	mr            sync.RWMutex
	result        []Result
	resultFn      recordResultFn
	mrs           sync.RWMutex
	resources     []schema.GroupVersionResource
	addResourceFn erAddResourceFn
}

type erAddResourceFn func(schema.GroupVersionResource)

func (r *er) recordResult(result Result) {
	r.mr.Lock()
	defer r.mr.Unlock()
	r.result = append(r.result, result)
}

func (r *er) addResource(er schema.GroupVersionResource) {
	r.mrs.Lock()
	defer r.mrs.Unlock()
	found := false
	for _, resource := range r.resources {
		if resource.Group == er.Group &&
			resource.Version == er.Version &&
			resource.Resource == er.Resource {
			return
		}
	}
	if !found {
		r.resources = append(r.resources, er)
	}
}

func (r *er) validateLcncHookNopFn(o Origin, v map[string]LcncGvrObject) {}

func (r *er) validateLcncGvrObjectFn(o Origin, idx int, n string, v LcncGvrObject) {
	gvr, err := GetGVR(v.Gvr)
	if err != nil {
		r.recordResult(Result{
			Origin: o,
			Index:  idx,
			Name:   n,
			Error:  err.Error(),
		})
	}
	r.addResource(*gvr)
}

func (r *er) validateBlockFn(o Origin, idx int, v LcncBlock) {
	value, err := GetValue(*v.For.Range)
	if err != nil {
		r.recordResult(Result{
			Origin: o,
			Index:  idx,
			Name:   "",
			Error:  err.Error(),
		})
	}
	if value.Kind == GVRKind {
		r.addResource(*value.Gvr)
	}
}

func (r *er) validateVarFn(o Origin, block bool, idx int, vertexName string, v LcncVar) {
	if v.Slice != nil {
		r.validateValue(o, block, idx, vertexName, v.Slice.LcncValue)
	}
	if v.Map != nil {
		r.validateValue(o, block, idx, vertexName, v.Map.LcncValue)
	}
}

func (r *er) validateValue(o Origin, block bool, idx int, vertexName string, v LcncValue) {
	dv := ""
	if v.String != nil {
		dv = *v.String
	}
	if v.Query != nil {
		dv = *v.Query
	}
	value, err := GetValue(dv)
	if err != nil {
		r.recordResult(Result{
			Origin: o,
			Index:  idx,
			Name:   "",
			Error:  err.Error(),
		})
	}
	if value.Kind == GVRKind {
		r.addResource(*value.Gvr)
	}
}

func (r *er) validateFunctionFn(o Origin, block bool, idx int, vertexName string, v LcncFunction) {
	for _, v := range v.Output {
		value, err := GetValue(v)
		if err != nil {
			r.recordResult(Result{
				Origin: o,
				Index:  idx,
				Name:   vertexName,
				Error:  err.Error(),
			})
		}
		if value.Kind == GVRKind {
			r.addResource(*value.Gvr)
		}
	}
}

func (r *er) validateServiceNopFn(o Origin, block bool, idx int, vertexName string, v LcncFunction) {}
