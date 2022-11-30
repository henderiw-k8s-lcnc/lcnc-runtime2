package dag

import (
	"fmt"
	"time"
)

// ResultFunc is the callback used for gathering the
// result during graph execution.
type ResultFunc func(*ResultEntry)

type ResultEntry struct {
	vertexName string
	duration   time.Duration
}

func (r *dag) recordResult(re *ResultEntry) {
	r.mr.Lock()
	defer r.mr.Unlock()
	r.result = append(r.result, re)
}

func (r *dag) GetWalkResult() {
	r.mr.RLock()
	defer r.mr.RUnlock()
	for i, result := range r.result {
		fmt.Printf("result order: %d vertex: %s, duration %s\n", i, result.vertexName, result.duration)
	}
}
