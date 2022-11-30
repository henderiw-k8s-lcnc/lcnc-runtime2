package dag

import (
	"fmt"
	"sync"
	"time"
)

type vertexcontext struct {
	vertexName string
	dep        bool
	//wg         *sync.WaitGroup
	cancelCh   chan struct{}

	m       sync.RWMutex
	doneChs map[string]chan bool
	depChs  map[string]chan bool

	visited  time.Time
	start    time.Time
	finished time.Time

	// callback
	recordResult ResultFunc
}

func (r *vertexcontext) AddDoneCh(n string, c chan bool) {
	r.m.Lock()
	defer r.m.Unlock()
	r.doneChs[n] = c
}

func (r *vertexcontext) AddDepCh(n string, c chan bool) {
	r.m.Lock()
	defer r.m.Unlock()
	r.depChs[n] = c
}

func (r *vertexcontext) isFinished() bool {
	return !r.finished.IsZero()
}

func (r *vertexcontext) hasStarted() bool {
	return !r.start.IsZero()
}

func (r *vertexcontext) isVisted() bool {
	if r == nil {
		return true
	}
	return !r.visited.IsZero()
}

func (r *vertexcontext) getDuration() time.Duration {
	return r.finished.Sub(r.start)
}

func (r *vertexcontext) run() {
	//fmt.Printf("runcontext vertex: %s run\n", r.vertexName)
	

	r.start = time.Now()
	// todo execute the function
	fmt.Printf("%s fn executed, doneChs: %v\n", r.vertexName, r.doneChs)
	r.finished = time.Now()

	// callback function to capture the result
	r.recordResult(&ResultEntry{vertexName: r.vertexName, duration: r.getDuration()})

	if r.dep {
		// inform all dependent/downstream verteces that the job is finished
		for vertexName, doneCh := range r.doneChs {
			doneCh <- true
			close(doneCh)
			fmt.Printf("%s -> %s send done\n", r.vertexName, vertexName)
		}
	}
}

func (r *vertexcontext) waitDependencies() bool {
	// for each dependency wait till a it completed, either through
	// the dependency Channel or cancel or
	//fmt.Printf("runcontext vertex: %s waitDependencies depChs: %v\n", r.vertexName, r.depChs)

	fmt.Printf("%s wait dependencies: %v\n", r.vertexName, r.depChs)
	DepSatisfied:
	for depVertexName, depCh := range r.depChs {
		//fmt.Printf("waitDependencies %s -> %s\n", depVertexName, r.vertexName)
		//DepSatisfied:
		for {
			select {
			case d, ok := <-depCh:
				fmt.Printf("%s -> %s rcvd done, d: %t, ok: %t\n", depVertexName, r.vertexName, d, ok)
				if ok {
					continue DepSatisfied
				}
				if !d {
					// dependency failed
					return false
				}
				continue DepSatisfied
			case <-r.cancelCh:
				// we can return, since someone cancelled the operation
				return false
			case <-time.After(time.Second * 5):
				fmt.Printf("wait timeout vertex: %s is waiting for %s\n", r.vertexName, depVertexName)
			}
		}
	}
	fmt.Printf("%s finished waiting\n", r.vertexName)
	return true
}
