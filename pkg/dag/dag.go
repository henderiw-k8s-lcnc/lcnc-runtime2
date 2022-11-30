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

	Walk(vertexName string)
	TransitiveReduction() DAG
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
}

func NewDAG() DAG {
	return &dag{
		vertices:  make(map[string]interface{}),
		downEdges: make(map[string]map[string]struct{}),
		upEdges:   make(map[string]map[string]struct{}),
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

func (d *dag) TransitiveReduction() DAG {
	// initialize a new dag
	newdag := NewDAG()
	for vertexName, v := range d.GetVertices() {
		newdag.AddVertex(vertexName, v)
	}

	for vertexName := range d.GetVertices() {
		fmt.Printf("############## ORIGIN VERTEX: %s ###############\n", vertexName)
		// we initialize the vertexdeptch map as 1 since 0 is used for uninitialized verteces
		// 0 is also used to avoid adding the vertex back in the graph
		d.initVertexDepth(vertexName, 1)
		d.transitiveReduction(vertexName, newdag)
	}
	return newdag

}

func (d *dag) initVertexDepth(from string, depth int) {
	//d.mvd.Lock()
	//defer d.mvd.Unlock()
	if depth == 1 {
		d.vertexDepth = map[string]int{}
	}
	d.vertexDepth[from] = depth
	// increase the depth for further recursive transitiveReduction
	depth++

	downEdges := d.GetDownVertexes(from)
	if len(downEdges) == 0 {
		// no downstream edge -> we are done
		return
	}
	wg := new(sync.WaitGroup)
	wg.Add(len(downEdges))
	for _, downEdge := range downEdges {
		downEdge := downEdge
		go func() {
			defer wg.Done()
			d.initVertexDepth(downEdge, depth)
		}()
	}
	wg.Wait()
}

func (d *dag) transitiveReduction(from string, newdag DAG) {
	//d.mvd.RLock()
	//defer d.mvd.RUnlock()
	fmt.Printf("from: %s, upVerteces: %v\n", from, d.GetUpVertexes(from))
	bestVertexDepth := 0
	for _, upVertex := range d.GetUpVertexes(from) {
		if d.vertexDepth[upVertex] > bestVertexDepth {
			bestVertexDepth = d.vertexDepth[upVertex]
		}
	}
	fmt.Printf("from: %s, bestVertexDepth: %v\n", from, bestVertexDepth)
	for _, upVertex := range d.GetUpVertexes(from) {
		// if bestVertexDepth == 0 it means we refer to an uninitialized vertex and we dont need
		// to process this.
		if bestVertexDepth != 0 {
			if d.vertexDepth[upVertex] == bestVertexDepth {
				newdag.Connect(upVertex, from)
				fmt.Printf("connect: from %s, to %s\n", upVertex, from)
				//newdag.AddUpEdge(from, upVertex)
			} else {
				fmt.Printf("transitive reduction: from %s, to %s\n", upVertex, from)
			}
		}
	}

	downEdges := d.GetDownVertexes(from)
	if len(downEdges) == 0 {
		// no downstream edge -> we are done
		return
	}
	wg := new(sync.WaitGroup)
	wg.Add(len(downEdges))
	for _, downEdge := range downEdges {
		/*
			if from == "root" {
				newdag.AddDownEdge("root", downEdge)
			}
		*/
		downEdge := downEdge
		go func() {
			defer wg.Done()
			d.transitiveReduction(downEdge, newdag)
		}()
	}
	wg.Wait()
}
