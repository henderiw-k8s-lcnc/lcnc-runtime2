package dag

import (
	"fmt"
	"time"
)

func (r *dag) Walk(from string) {
	// walk initialization
	r.initWalk()
	start := time.Now()
	r.walk(from, true, 1)
	// add total as a last entry in the result
	r.recordResult(&ResultEntry{
		vertexName: "total",
		duration:   time.Since(start),
	})
}

func (r *dag) initWalk() {
	//d.wg = new(sync.WaitGroup)
	r.cancelCh = make(chan struct{})
	r.result = []*ResultEntry{}
	r.walkMap = map[string]*vertexContext{}
	for vertexName := range r.GetVertices() {
		//fmt.Printf("init vertexName: %s\n", vertexName)
		r.walkMap[vertexName] = &vertexContext{
			vertexName: vertexName,
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
	for vertexName, wCtx := range r.walkMap {
		// only run these channels when we want to add dependency validation
		for _, depVertexName := range r.GetUpVertexes(vertexName) {
			//fmt.Printf("vertexName: %s, depBVertexName: %s\n", vertexName, depVertexName)
			depCh := make(chan bool, 1)
			r.walkMap[depVertexName].AddDoneCh(vertexName, depCh) // send when done
			wCtx.AddDepCh(depVertexName, depCh)                   // rcvr when done
		}
		doneFnCh := make(chan bool, 1)
		wCtx.doneFnCh = doneFnCh
		r.fnDoneMap[vertexName] = doneFnCh
	}
}

func (r *dag) walk(from string, init bool, depth int) {
	wCtx := r.getWalkContext(from)
	// avoid scheduling a vertex that is already visted
	if !wCtx.isVisted() {
		wCtx.visited = time.Now()
		// execute the vertex function
		fmt.Printf("%s scheduled\n", wCtx.vertexName)
		go func() {
			if !r.dependenciesFinished(wCtx.depChs) {
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
	depth++
	for _, downEdge := range r.GetDownVertexes(from) {
		go func(downEdge string) {
			r.walk(downEdge, false, depth)
		}(downEdge)
	}
	if init {
		r.waitFunctionCompletion()
	}
}

func (r *dag) getWalkContext(s string) *vertexContext {
	r.mw.RLock()
	defer r.mw.RUnlock()
	return r.walkMap[s]
}

func (r *dag) dependenciesFinished(dep map[string]chan bool) bool {
	for vertexName := range dep {
		if !r.getWalkContext(vertexName).isFinished() {
			return false
		}
	}
	return true
}

func (r *dag) waitFunctionCompletion() {
	fmt.Printf("main walk wait waiting for function completion...\n")
DepSatisfied:
	for vertexName, doneFnCh := range r.fnDoneMap {
		for {
			select {
			case d, ok := <-doneFnCh:
				fmt.Printf("main walk wait rcvd fn done from %s, d: %t, ok: %t\n", vertexName, d, ok)

				continue DepSatisfied
			case <-r.cancelCh:
				// we can return, since someone cancelled the operation
				return
			case <-time.After(time.Second * 5):
				fmt.Printf("main walk wait timeout, waiting for %s\n", vertexName)
			}
		}
	}
	fmt.Printf("main walk wait function completion waiting finished - bye !\n")
}
