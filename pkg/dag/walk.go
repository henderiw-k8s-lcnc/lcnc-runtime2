package dag

import (
	"fmt"
	"time"
)

/*
func (d *dag) Walk(from string) {
	d.mde.RLock()
	defer d.mde.RUnlock()
	//fmt.Printf("walk from %s\n", from)

	var wg sync.WaitGroup

	if _, ok := d.downEdges[from]; !ok {
		fmt.Printf("walk from %s -> done\n", from)
		// no downstream edge -> we are done
		return
	}
	downEdges := make([]string, 0, len(d.downEdges[from]))
	for to := range d.downEdges[from] {
		downEdges = append(downEdges, to)
	}
	for _, downEdge := range downEdges {
		wg.Add(1)
		downEdge := downEdge
		fmt.Printf("walk from %s -> %s\n", from, downEdge)
		go func() {
			defer wg.Done()
			d.Walk(downEdge)

		}()
	}

	wg.Wait()
}
*/

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
	r.walkMap = map[string]*vertexcontext{}
	for vertexName := range r.GetVertices() {
		//fmt.Printf("init vertexName: %s\n", vertexName)
		r.walkMap[vertexName] = &vertexcontext{
			//dep:        wc.Dep,
			vertexName: vertexName,
			//wg:         d.wg,
			cancelCh: r.cancelCh,
			doneChs:  make(map[string]chan bool), //snd
			depChs:   make(map[string]chan bool), //rcv
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
		//if wc.Dep {
		for _, depVertexName := range r.GetUpVertexes(vertexName) {
			//fmt.Printf("vertexName: %s, depBVertexName: %s\n", vertexName, depVertexName)
			depCh := make(chan bool, 1)
			r.walkMap[depVertexName].AddDoneCh(vertexName, depCh) // send when done
			wCtx.AddDepCh(depVertexName, depCh)                   // rcvr when done
		}
		//}
		// add a done channel for all the fn
		doneFnCh := make(chan bool, 1)
		wCtx.doneFnCh = doneFnCh
		r.fnDoneMap[vertexName] = doneFnCh
	}
}

func (r *dag) walk(from string, init bool, depth int) {
	//wg := new(sync.WaitGroup)
	/*
		if wc.WalkInitFn != nil {
			if init {
				wc.WalkInitFn()
			}
		}
		if wc.WalkEntryFn != nil {
			wc.WalkEntryFn(from, depth)
		}
	*/
	wCtx := r.getWalkContext(from)
	// avoid scheduling a vertex that is already visted
	if !wCtx.isVisted() {
		//wg.Add(1)
		wCtx.visited = time.Now()
		// execute the vertex function
		fmt.Printf("%s scheduled\n", wCtx.vertexName)
		go func() {
			//defer wg.Done()
			//if wc.Dep {
			if !r.dependenciesFinished(wCtx.depChs) {
				fmt.Printf("%s not finished\n", from)
			}
			if !wCtx.waitDependencies() {
				// TODO gather info why the failure occured
				return
			}
			//	}
			// execute the vertex function
			wCtx.run()
		}()
	}
	// continue walking the graph
	downEdges := r.GetDownVertexes(from)
	if len(downEdges) == 0 {
		return
	}
	// increment the depth
	depth++
	for _, downEdge := range downEdges {
		// continue the recursive walk
		//wg.Add(1)
		go func(downEdge string) {
			r.walk(downEdge, false, depth)
		}(downEdge)
		//wg.Done()
	}
	//wg.Wait()
	if init {
		r.waitFunctionCompletion()
	}
}

func (r *dag) getWalkContext(s string) *vertexcontext {
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
	fmt.Printf("function completion waiting...\n")
DepSatisfied:
	for vertexName, doneFnCh := range r.fnDoneMap {
		// for transitive reduction we dont schedule all the vertices so we should
		// only validate the once that are scheduled to run
		//fmt.Printf("%s vertexdepth: %d\n", vertexName, r.getVertexDepth(vertexName))
		//if wc.TransitiveReduction && r.getVertexDepth(vertexName) == 0 {
		//	fmt.Printf("%s not waiting for fn completion while doing transitive reduction since vertexmap is: %d\n", vertexName, r.getVertexDepth(vertexName))
		//	continue DepSatisfied
		//}
		//fmt.Printf("waitDependencies %s -> %s\n", depVertexName, r.vertexName)
		//DepSatisfied:
		for {
			select {
			case d, ok := <-doneFnCh:
				fmt.Printf("%s -> walk rcvd fn done, d: %t, ok: %t\n", vertexName, d, ok)

				continue DepSatisfied
			case <-r.cancelCh:
				// we can return, since someone cancelled the operation
				return
			case <-time.After(time.Second * 5):
				fmt.Printf("wait timeout main walk is waiting for %s\n", vertexName)
			}
		}
	}
	fmt.Printf("function completion waiting finished - bye !\n")
}
