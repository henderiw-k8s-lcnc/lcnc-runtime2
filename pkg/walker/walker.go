package walker

import (
	"fmt"
	"sync"
	"time"

	"github.com/henderiw-k8s-lcnc/lcnc-runtime2/pkg/dag"
)

type Walker interface {
	Walk() error
	GetResult() 
}

func New(d dag.DAG) Walker {
	r := &walker{
		d:        d,
		wg:       new(sync.WaitGroup),
		cancelCh: make(chan struct{}),
		result:   []*ResultEntry{},
	}
	r.initWalkerContext()
	return r
}

type walker struct {
	d         dag.DAG
	m         sync.RWMutex
	walkerMap map[string]*walkerContext

	wg       *sync.WaitGroup
	cancelCh chan struct{}

	result []*ResultEntry
}

func (r *walker) initWalkerContext() {
	r.walkerMap = map[string]*walkerContext{}
	for vertexName := range r.d.GetVertices() {
		r.walkerMap[vertexName] = &walkerContext{
			vertexName: vertexName,
			wg:         r.wg,
			cancelCh:   r.cancelCh,
			deps:       r.d.GetUpVertexes(vertexName),
			// callback
			recordResult: r.recordResult,
		}
	}
}

func (r *walker) dependenciesFinished(dep []string) bool {
	for _, vertexName := range dep {
		if !r.getWalkerContext(vertexName).isFinished() {
			return false
		}
	}
	return true
}

func (r *walker) getWalkerContext(s string) *walkerContext {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.walkerMap[s]
}

func (r *walker) Walk() error {
	start := time.Now()
	r.walk("root")
	// add total as a last entry in the result
	r.result = append(r.result, &ResultEntry{vertexName: "total", duration: time.Now().Sub(start)})
	return nil
}

func (r *walker) walk(from string) error {
	downEdges := r.d.GetDownVertexes(from)
	// if no more downstream edges -> we are done
	if len(downEdges) == 0 {
		return nil
	}
	for _, downEdge := range downEdges {
		// statement is here to keep go func happy
		downEdge := downEdge

		wCtx := r.getWalkerContext(downEdge)
		if !wCtx.isScheduled() {
			wCtx.scheduled = time.Now()
			r.wg.Add(1)
			// execute the vertex
			go func() {
				if !r.dependenciesFinished(wCtx.deps) {
					fmt.Printf("vertex: %s not finished\n", downEdge)
				}
				wCtx.run()
			}()
			// continue the recursive walk
			go func() {
				r.walk(downEdge)
			}()
		}
	}
	r.wg.Wait()
	return nil
}

func (r *walker) recordResult(re *ResultEntry) {
	var l sync.Mutex
	l.Lock()
	defer l.Unlock()
	r.result = append(r.result, re)
}

func (r *walker) GetResult() {
	for i, result := range r.result {
		fmt.Printf("result order: %d vertex: %s, duration %s\n", i, result.vertexName, result.duration)
	}
}