package walker

import (
	"sync"
	"time"
)

type walkerContext struct {
	vertexName string
	wg         *sync.WaitGroup
	cancelCh   chan struct{}
	deps       []string

	scheduled time.Time
	start     time.Time
	finished  time.Time

	// callback
	recordResult ResultFunc
}

func (r *walkerContext) isFinished() bool {
	return !r.finished.IsZero()
}

func (r *walkerContext) hasStarted() bool {
	return !r.start.IsZero()
}

func (r *walkerContext) isScheduled() bool {
	return !r.scheduled.IsZero()
}

func (r *walkerContext) getDuration() time.Duration {
	return r.finished.Sub(r.start)
}

func (r *walkerContext) run() {
	//fmt.Printf("runcontext vertex: %s run\n", r.vertexName)
	defer r.wg.Done()

	r.start = time.Now()
	// todo execute the function
	r.finished = time.Now()

	// callback function
	r.recordResult(&ResultEntry{vertexName: r.vertexName, duration: r.getDuration()})

}



/*
func (r *runContext) waitDependencies() bool {
	// for each dependency wait till a it completed, either through
	// the dependency Channel or cancel or
	fmt.Printf("runcontext vertex: %s waitDependencies depChs: %v\n", r.vertexName, r.depChs)
	for depVertexName, depCh := range r.depChs {

	DepSatisfied:
		for {
			select {
			case d := <-depCh:
				if !d {
					// dependency failed
					return false
				}
				break DepSatisfied
			case <-r.cancelCh:
				// we can return, since someone cancelled the operation
				return false
			case <-time.After(time.Second * 5):
				fmt.Printf("runcontext vertex: %s is waiting for %s\n", r.vertexName, depVertexName)
			}
		}
	}
	return true
}
*/
