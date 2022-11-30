package walker

import (
	"fmt"
	"sync"
	"time"
)

type walkerContext struct {
	vertexName string
	wg         *sync.WaitGroup
	cancelCh   chan struct{}
	deps       []string

	m       sync.RWMutex
	doneChs map[string]chan bool
	depChs  map[string]chan bool

	scheduled time.Time
	start     time.Time
	finished  time.Time

	// callback
	recordResult ResultFunc
}

func (r *walkerContext) AddDoneCh(n string, c chan bool) {
	r.m.Lock()
	defer r.m.Unlock()
	r.doneChs[n] = c
}

func (r *walkerContext) AddDepCh(n string, c chan bool) {
	r.m.Lock()
	defer r.m.Unlock()
	r.depChs[n] = c
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

	// callback function to capture the reuslt
	r.recordResult(&ResultEntry{vertexName: r.vertexName, duration: r.getDuration()})

	// inform all dependent/downstream verteces that the job is finished
	for vertexName, doneCh := range r.doneChs {
		doneCh <- true
		close(doneCh)
		fmt.Printf("%s -> %s send done\n", r.vertexName, vertexName)
	}
}

func (r *walkerContext) waitDependencies() bool {
	// for each dependency wait till a it completed, either through
	// the dependency Channel or cancel or
	//fmt.Printf("runcontext vertex: %s waitDependencies depChs: %v\n", r.vertexName, r.depChs)
DepSatisfied:
	for depVertexName, depCh := range r.depChs {
		//fmt.Printf("waitDependencies %s -> %s\n", depVertexName, r.vertexName)
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
