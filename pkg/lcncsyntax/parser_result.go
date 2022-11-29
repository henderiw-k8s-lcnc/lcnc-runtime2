package lcncsyntax

import "sync"

type Result struct {
	Origin Origin `json:"origin,omitempty" yaml:"origin,omitempty"`
	Index  int    `json:"index,omitempty" yaml:"index,omitempty"`
	Name   string `json:"name,omitempty" yaml:"name,omitempty"`
	Error  string  `json:"error,omitempty" yaml:"error,omitempty"`
}

type recordResultFn func(Result)

func (r *lcncparser) recordSyntaxValidationResult(result Result) {
	var l sync.Mutex
	l.Lock()
	defer l.Unlock()
	r.syntaxValidationResult = append(r.syntaxValidationResult, result)
}
