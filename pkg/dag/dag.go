package dag

import (
	"fmt"
	"sync"
)

type DAG interface {
	AddVertex(s string, v interface{}) error
	Connect(from, to string)
	AddDownEdge(from, to string)
	AddUpEdge(from, to string)
	GetVertex(s string) bool
	GetVertices() map[string]interface{}
	GetDownVertexes(from string) []string
	GetUpVertexes(from string) []string

	GetDependencyMap(vertexName string)
	Walk(vertexName string, wc WalkConfig)
	GetWalkResult()
	TransitiveReduction()
}

// used for returning
type Edge struct {
	From string
	To   string
}

type dag struct {
	// vertices first key is the vertexName
	mv       sync.RWMutex
	vertices map[string]interface{}
	// downEdges/upEdges
	// 1st key is from, 2nd key is to
	mde       sync.RWMutex
	downEdges map[string]map[string]struct{}
	mue       sync.RWMutex
	upEdges   map[string]map[string]struct{}
	// used for transit reduction
	mvd         sync.RWMutex
	vertexDepth map[string]int
	// used for the walker
	//wg       *sync.WaitGroup
	cancelCh chan struct{}
	mw       sync.RWMutex
	walkMap  map[string]*vertexcontext
	mr       sync.RWMutex
	result   []*ResultEntry
}

func NewDAG() DAG {
	return &dag{
		vertices:  make(map[string]interface{}),
		downEdges: make(map[string]map[string]struct{}),
		upEdges:   make(map[string]map[string]struct{}),
		//wg:        new(sync.WaitGroup),
	}
}

func (d *dag) AddVertex(s string, v interface{}) error {
	d.mv.Lock()
	defer d.mv.Unlock()

	// validate duplicate entry
	if _, ok := d.vertices[s]; ok {
		return fmt.Errorf("duplicate vertex entry: %s", s)
	}
	d.vertices[s] = v
	return nil
}

func (d *dag) GetVertices() map[string]interface{} {
	d.mv.RLock()
	defer d.mv.RUnlock()
	return d.vertices
}

func (d *dag) GetVertex(s string) bool {
	d.mv.RLock()
	defer d.mv.RUnlock()
	_, ok := d.vertices[s]
	return ok
}

func (d *dag) Connect(from, to string) {
	d.AddDownEdge(from, to)
	d.AddUpEdge(to, from)
}

func (d *dag) Disconnect(from, to string) {
	d.DeleteDownEdge(from, to)
	d.DeleteUpEdge(to, from)
}

func (d *dag) AddDownEdge(from, to string) {
	d.mde.Lock()
	defer d.mde.Unlock()

	//fmt.Printf("addDownEdge: from: %s, to: %s\n", from, to)

	// initialize the from entry if it does not exist
	if _, ok := d.downEdges[from]; !ok {
		d.downEdges[from] = make(map[string]struct{})
	}
	if _, ok := d.downEdges[from][to]; ok {
		//  down edge entry already exists
		return
	}
	// add entry
	d.downEdges[from][to] = struct{}{}
}

func (d *dag) DeleteDownEdge(from, to string) {
	d.mde.Lock()
	defer d.mde.Unlock()

	//fmt.Printf("deleteDownEdge: from: %s, to: %s\n", from, to)
	if de, ok := d.downEdges[from]; ok {
		if _, ok := d.downEdges[from][to]; ok {
			delete(de, to)
		}
	}
}

func (d *dag) GetDownVertexes(from string) []string {
	d.mde.RLock()
	defer d.mde.RUnlock()

	edges := make([]string, 0)
	if fromVertex, ok := d.downEdges[from]; ok {
		for to := range fromVertex {
			edges = append(edges, to)
		}
	}
	return edges
}

func (d *dag) AddUpEdge(from, to string) {
	d.mue.Lock()
	defer d.mue.Unlock()

	//fmt.Printf("addUpEdge: from: %s, to: %s\n", from, to)

	// initialize the from entry if it does not exist
	if _, ok := d.upEdges[from]; !ok {
		d.upEdges[from] = make(map[string]struct{})
	}
	if _, ok := d.upEdges[from][to]; ok {
		//  up edge entry already exists
		return
	}
	// add entry
	d.upEdges[from][to] = struct{}{}
}

func (d *dag) DeleteUpEdge(from, to string) {
	d.mue.Lock()
	defer d.mue.Unlock()

	//fmt.Printf("deleteUpEdge: from: %s, to: %s\n", from, to)
	if ue, ok := d.upEdges[from]; ok {
		if _, ok := d.upEdges[from][to]; ok {
			delete(ue, to)
		}
	}
}

func (d *dag) GetUpEdges(from string) []Edge {
	d.mue.RLock()
	defer d.mue.RUnlock()

	edges := make([]Edge, 0)
	if fromVertex, ok := d.upEdges[from]; ok {
		for to := range fromVertex {
			edges = append(edges, Edge{From: from, To: to})
		}
	}
	return edges
}

func (d *dag) GetUpVertexes(from string) []string {
	d.mue.RLock()
	defer d.mue.RUnlock()

	upVerteces := []string{}
	if fromVertex, ok := d.upEdges[from]; ok {
		for to := range fromVertex {
			upVerteces = append(upVerteces, to)
		}
	}
	return upVerteces
}

func (d *dag) GetDependencyMap(from string) {
	fmt.Println("######### dependency map start ###########")
	d.getDependencyMap(from, 0)
	fmt.Println("######### dependency map end   ###########")
}

func (d *dag) getDependencyMap(from string, indent int) {
	fmt.Printf("%s:\n", from)
	for _, upVertex := range d.GetUpVertexes(from) {
		fmt.Printf("-> %s\n", upVertex)
	}
	indent++
	for _, downVertex := range d.GetDownVertexes(from) {
		d.getDependencyMap(downVertex, indent)
	}
}

func (d *dag) TransitiveReduction() {
	// initialize a new dag
	//newdag := NewDAG()
	//for vertexName, v := range d.GetVertices() {
	//	newdag.AddVertex(vertexName, v)
	//}

	for vertexName := range d.GetVertices() {
		fmt.Printf("##### TRANSIT REDUCTION VERTEX START: %s ###############\n", vertexName)
		// we initialize the vertexdeptch map as 1 since 0 is used for uninitialized verteces
		// 0 is also used to avoid adding the vertex back in the graph
		
		// initVertexeDepthMap
		//d.wg = new(sync.WaitGroup)
		wc := WalkConfig{
			Dep:         false,
			WalkInitFn:  d.initVertexDepthMap,
			WalkEntryFn: d.addVertexDepth,
		}
		d.initWalk(wc)
		d.walk(vertexName, true, 1, wc)
		//d.transitiveReduction(vertexName, true)
		//d.wg = new(sync.WaitGroup)
		wc = WalkConfig{
			Dep:         false,
			WalkEntryFn: d.processTransitiveReducation,
		}
		d.initWalk(wc)
		d.walk(vertexName, true, 1, wc)
		fmt.Printf("##### TRANSIT REDUCTION VERTEX ENDED: %s ###############\n", vertexName)
	}
}

func (d *dag) processTransitiveReducation(from string, depth int) {
	//fmt.Printf("from: %s, upVerteces: %v\n", from, d.GetUpVertexes(from))
	bestVertexDepth := d.getbestVertexDepth(from)
	//fmt.Printf("from: %s, bestVertexDepth: %v\n", from, bestVertexDepth)
	for _, upVertex := range d.GetUpVertexes(from) {
		// if bestVertexDepth == 0 it means we refer to an uninitialized vertex and we dont need
		// to process this.
		if bestVertexDepth != 0 {
			if d.getVertexDepth(upVertex) != bestVertexDepth {
				fmt.Printf("transitive reduction %s -> %s\n", upVertex, from)
				d.Disconnect(upVertex, from)
			}
		}
	}
}

/*
func (d *dag) initVertexDepth(from string, init bool, depth int) {
	if init {
		depth = 1
		d.wg = new(sync.WaitGroup)
		d.initVertexDepthMap()
	}
	d.addVertexDepth(from, depth)
	// increase the depth of the dag

	fmt.Printf("initvertexDepth %s, downvertices: %v depth: %d\n", from, d.GetDownVertexes(from), depth)
	depth++
	downEdges := d.GetDownVertexes(from)
	if len(downEdges) == 0 {
		return
	}

	// continue walk the dag
	for _, downEdge := range downEdges {
		d.wg.Add(1)
		downEdge := downEdge
		go func() {
			defer d.wg.Done()
			d.initVertexDepth(downEdge, false, depth)
		}()
	}
	d.wg.Wait()
}
*/

/*
func (d *dag) transitiveReduction(from string, init bool) {
	fmt.Printf("from: %s, upVerteces: %v\n", from, d.GetUpVertexes(from))
	bestVertexDepth := d.getbestVertexDepth(from)
	fmt.Printf("from: %s, bestVertexDepth: %v\n", from, bestVertexDepth)
	for _, upVertex := range d.GetUpVertexes(from) {
		// if bestVertexDepth == 0 it means we refer to an uninitialized vertex and we dont need
		// to process this.
		if bestVertexDepth != 0 {
			if d.getVertexDepth(upVertex) != bestVertexDepth {
				d.Disconnect(upVertex, from)
			}

			//	if d.vertexDepth[upVertex] == bestVertexDepth {
			//		newdag.Connect(upVertex, from)
			//		fmt.Printf("connect: from %s, to %s\n", upVertex, from)
			//		//newdag.AddUpEdge(from, upVertex)
			//	} else {
			//		fmt.Printf("transitive reduction: from %s, to %s\n", upVertex, from)
			//	}
		}
	}

	// this retuns to the main loop
	downEdges := d.GetDownVertexes(from)
	if len(downEdges) == 0 {
		return
	}

	for _, downEdge := range downEdges {
		d.wg.Add(1)
		downEdge := downEdge
		go func() {
			defer d.wg.Done()
			d.transitiveReduction(downEdge, false)
		}()
	}
	d.wg.Wait()
}
*/

func (d *dag) initVertexDepthMap() {
	d.mvd.Lock()
	defer d.mvd.Unlock()
	d.vertexDepth = map[string]int{}
}

func (d *dag) getVertexDepth(n string) int {
	d.mvd.RLock()
	defer d.mvd.RUnlock()
	if depth, ok := d.vertexDepth[n]; ok {
		return depth
	}
	return 0
}

func (d *dag) addVertexDepth(n string, depth int) {
	d.mvd.Lock()
	defer d.mvd.Unlock()
	d.vertexDepth[n] = depth
}

func (d *dag) getbestVertexDepth(from string) int {
	bestVertexDepth := 0
	for _, upVertex := range d.GetUpVertexes(from) {
		upVertexDepth := d.getVertexDepth(upVertex)
		if upVertexDepth > bestVertexDepth {
			bestVertexDepth = upVertexDepth
		}
	}
	return bestVertexDepth
}
