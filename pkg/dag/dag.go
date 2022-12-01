package dag

import (
	"fmt"
	"os"
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
	Walk(vertexName string)
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
	//mfd       sync.RWMutex
	fnDoneMap map[string]chan bool
	mr        sync.RWMutex
	result    []*ResultEntry
}

func NewDAG() DAG {
	return &dag{
		vertices:  make(map[string]interface{}),
		downEdges: make(map[string]map[string]struct{}),
		upEdges:   make(map[string]map[string]struct{}),
		//wg:        new(sync.WaitGroup),
		fnDoneMap: make(map[string]chan bool),
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
	fmt.Println("######### dependency map verteces start ###########")
	for vertexName := range d.GetVertices() {
		fmt.Printf("%s\n", vertexName)
	}
	fmt.Println("######### dependency map verteces end ###########")
	fmt.Println("######### dependency map start ###########")
	d.getDependencyMap(from, 0)
	fmt.Println("######### dependency map end   ###########")
}

func (d *dag) getDependencyMap(from string, indent int) {
	fmt.Printf("%s:\n", from)
	for _, upVertex := range d.GetUpVertexes(from) {
		found := d.checkVertex(upVertex)
		if !found {
			fmt.Printf("upVertex %s no found in vertices\n", upVertex)
			os.Exit(1)
		}
		fmt.Printf("-> %s\n", upVertex)
	}
	indent++
	for _, downVertex := range d.GetDownVertexes(from) {
		found := d.checkVertex(downVertex)
		if !found {
			fmt.Printf("upVertex %s no found in vertices\n", downVertex)
			os.Exit(1)
		}
		d.getDependencyMap(downVertex, indent)
	}
}

func (d *dag) checkVertex(s string) bool {
	found := false
	for vertexName := range d.GetVertices() {
		if vertexName == s {
			return true
		}
	}
	return found
}
