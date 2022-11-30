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

func New(d dag.DAG, root string) Walker {
	r := &walker{
		root:     root,
		d:        d,
		wg:       new(sync.WaitGroup),
		cancelCh: make(chan struct{}),

		result: []*ResultEntry{},
	}
	r.initWalkerContext()
	return r
}

type walker struct {
	root      string
	d         dag.DAG
	m         sync.RWMutex
	walkerMap map[string]*walkerContext

	wg       *sync.WaitGroup
	cancelCh chan struct{}

	mr     sync.RWMutex
	result []*ResultEntry
}

func (r *walker) initWalkerContext() {
	r.walkerMap = map[string]*walkerContext{}
	for vertexName := range r.d.GetVertices() {
		r.walkerMap[vertexName] = &walkerContext{
			vertexName: vertexName,
			wg:         r.wg,
			cancelCh:   r.cancelCh,
			doneChs:    make(map[string]chan bool), //snd
			depChs:     make(map[string]chan bool), //rcv
			// callback to gather the result
			recordResult: r.recordResult,
		}
	}
	// build the channel matrix to signal dependencies through channels
	// for every dependency (upstream relationship between verteces)
	// create a channel
	// add the channel to the upstreamm vertex doneCh map ->
	// usedto signal/send the vertex finished the function/job
	// add the channel to the downstream vertex depCh map ->
	// used to wait for the upstream vertex to signal the fn/job is done
	for vertexName, wCtx := range r.walkerMap {
		for _, depVertexName := range r.d.GetUpVertexes(vertexName) {
			depCh := make(chan bool)
			r.walkerMap[depVertexName].AddDoneCh(vertexName, depCh) // send when done
			wCtx.AddDepCh(depVertexName, depCh)                     // rcvr when done
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
	r.walk(r.root)
	// add total as a last entry in the result
	r.result = append(r.result, &ResultEntry{vertexName: "total", duration: time.Now().Sub(start)})
	return nil
}

func (r *walker) walk(from string) error {
	wCtx := r.getWalkerContext(from)
	// avoid scheduling a vertex that is already running
	if !wCtx.isScheduled() {
		wCtx.scheduled = time.Now()
		r.wg.Add(1)
		// execute the vertex function
		fmt.Printf("%s scheduled\n", wCtx.vertexName)
		go func() {
			if !r.dependenciesFinished(wCtx.deps) {
				fmt.Printf("%s not finished\n", from)
			}
			if !wCtx.waitDependencies() {
				// TODO gather info why the failure occured
				return
			}
			// execute the vertex function
			wCtx.run()
		}()
	}

	// continue walking the graph
	downEdges := r.d.GetDownVertexes(from)
	if len(downEdges) == 0 {
		// if no more downstream edges -> we are done
		return nil
	}
	for _, downEdge := range downEdges {
		// statement is here to keep go func happy
		downEdge := downEdge
		// continue the recursive walk
		go func() {
			r.walk(downEdge)
		}()

	}
	r.wg.Wait()
	return nil
}

func (r *walker) recordResult(re *ResultEntry) {
	r.mr.Lock()
	defer r.mr.Unlock()
	r.result = append(r.result, re)
}

func (r *walker) GetResult() {
	for i, result := range r.result {
		fmt.Printf("result order: %d vertex: %s, duration %s\n", i, result.vertexName, result.duration)
	}
}
