package lcnc

import (
	"errors"
	"fmt"
)

func (r *lcnc) AddOutputMapping(outVar, fnName string) error {
	r.mo.Lock()
	defer r.mo.Unlock()

	// validate duplicate entry
	if _, ok := r.output[outVar]; ok {
		return errors.New("duplicate outputVariable entry")
	}
	r.output[outVar] = fnName
	return nil
}

func (r *lcnc) HasOutputMapping(s string) bool {
	r.mo.RLock()
	defer r.mo.RUnlock()
	_, ok := r.output[s]
	return ok
}

func (r *lcnc) GetOutputMapping(s string) string {
	r.ml.RLock()
	defer r.ml.RUnlock()
	if fnName, ok := r.output[s]; ok {
		return fnName
	}
	return ""
}

func (r *lcnc) PrintOutputMappings() {
	for k, v := range r.output {
		fmt.Printf("output: key: %s, fnName: %s\n", k, v)
	}
}
