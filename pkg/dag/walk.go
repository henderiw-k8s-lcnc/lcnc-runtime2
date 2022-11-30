package dag

import (
	"fmt"
	"sync"
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

func (d *dag) Walk(from string, wc WalkConfig) {
	// walk initialization
	d.initWalk(wc)
	start := time.Now()
	d.walk(from, true, 1, wc)
	// add total as a last entry in the result
	d.recordResult(&ResultEntry{
		vertexName: "total",
		duration:   time.Now().Sub(start),
	})
}

func (d *dag) walk(from string, init bool, depth int, wc WalkConfig) {
	wg := new(sync.WaitGroup)
	if wc.WalkInitFn != nil {
		if init {
			wc.WalkInitFn()
		}
	}
	if wc.WalkEntryFn != nil {
		wc.WalkEntryFn(from, depth)
	}
	wCtx := d.getWalkContext(from)
	// avoid scheduling a vertex that is already visted
	if !wCtx.isVisted() {
		wg.Add(1)
		wCtx.visited = time.Now()
		// execute the vertex function
		fmt.Printf("%s scheduled\n", wCtx.vertexName)
		go func() {
			defer wg.Done()
			if wc.Dep {
				if !d.dependenciesFinished(wCtx.depChs) {
					fmt.Printf("%s not finished\n", from)
				}
				if !wCtx.waitDependencies() {
					// TODO gather info why the failure occured
					return
				}
			}
			// execute the vertex function
			wCtx.run()
		}()
	}
	// continue walking the graph
	downEdges := d.GetDownVertexes(from)
	if len(downEdges) == 0 {
		return
	}
	// increment the depth
	depth++
	for _, downEdge := range downEdges {
		// continue the recursive walk
		//wg.Add(1)
		go func(downEdge string) {
			d.walk(downEdge, false, depth, wc)
		}(downEdge)
		//wg.Done()
	}
	wg.Wait()
}

func (d *dag) initWalk(wc WalkConfig) {
	//d.wg = new(sync.WaitGroup)
	d.cancelCh = make(chan struct{})
	d.result = []*ResultEntry{}
	d.walkMap = map[string]*vertexcontext{}
	for vertexName := range d.GetVertices() {
		//fmt.Printf("init vertexName: %s\n", vertexName)
		d.walkMap[vertexName] = &vertexcontext{
			dep:        wc.Dep,
			vertexName: vertexName,
			//wg:         d.wg,
			cancelCh: d.cancelCh,
			doneChs:  make(map[string]chan bool), //snd
			depChs:   make(map[string]chan bool), //rcv
			// callback to gather the result
			recordResult: d.recordResult,
		}
	}
	if wc.Dep {
		// build the channel matrix to signal dependencies through channels
		// for every dependency (upstream relationship between verteces)
		// create a channel
		// add the channel to the upstreamm vertex doneCh map ->
		// usedto signal/send the vertex finished the function/job
		// add the channel to the downstream vertex depCh map ->
		// used to wait for the upstream vertex to signal the fn/job is done
		for vertexName, wCtx := range d.walkMap {
			for _, depVertexName := range d.GetUpVertexes(vertexName) {
				//fmt.Printf("vertexName: %s, depBVertexName: %s\n", vertexName, depVertexName)
				depCh := make(chan bool, 1)
				d.walkMap[depVertexName].AddDoneCh(vertexName, depCh) // send when done
				wCtx.AddDepCh(depVertexName, depCh)                   // rcvr when done
			}
		}
	}
}

func (d *dag) getWalkContext(s string) *vertexcontext {
	d.mw.RLock()
	defer d.mw.RUnlock()
	return d.walkMap[s]
}

func (d *dag) dependenciesFinished(dep map[string]chan bool) bool {
	for vertexName := range dep {
		if !d.getWalkContext(vertexName).isFinished() {
			return false
		}
	}
	return true
}
