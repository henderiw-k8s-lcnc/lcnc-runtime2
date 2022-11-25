package dag

import (
	"errors"
	"fmt"
	"sync"
)

type DAG interface {
	AddVertex(s string, v interface{}) error
	AddEdge(from, to string) error
	GetVertex(s string) bool
	GetVertices() map[string]interface{}
	GetEdges() []Edge
	Walk(vertexName string)
}

// used for returning
type Edge struct {
	From string
	To   string
}

type dag struct {
	mv       sync.RWMutex
	vertices map[string]interface{}
	me       sync.RWMutex
	// 1st key is from, 2nd key is to
	downEdges map[string]map[string]struct{}
}

func NewDAG() DAG {
	return &dag{
		vertices:  make(map[string]interface{}),
		downEdges: make(map[string]map[string]struct{}),
	}
}

func (d *dag) AddVertex(s string, v interface{}) error {
	d.mv.Lock()
	defer d.mv.Unlock()

	// validate duplicate entry
	if _, ok := d.vertices[s]; !ok {
		return errors.New("duplicate vertex entry")
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

func (d *dag) AddEdge(from, to string) error {
	d.me.Lock()
	defer d.me.Unlock()

	fmt.Printf("addEdge: from: %s, to: %s\n", from, to)

	if t, ok := d.downEdges[from]; ok {
		if _, ok := t[to]; ok {
			//  down edge entry already exists
			return nil
		}
	} else {
		d.downEdges[from] = make(map[string]struct{})
	}
	d.downEdges[from][to] = struct{}{}
	return nil
}

func (d *dag) GetEdges() []Edge {
	d.me.RLock()
	defer d.me.RUnlock()

	edges := make([]Edge, 0)
	for from, tos := range d.downEdges {
		for to := range tos {
			edges = append(edges, Edge{From: from, To: to})
		}
	}
	return edges
}

func (d *dag) Walk(from string) {
	d.me.RLock()
	defer d.me.RUnlock()

	if _, ok := d.downEdges[from]; !ok {
		// no downstream edge -> we are done
		return
	}
	downEdges := make([]string, 0, len(d.downEdges[from]))
	for to := range d.downEdges[from] {
		fmt.Printf("edge %s -> %s\n", from, to)
		downEdges = append(downEdges, to)
	}
	for _, downEdge := range downEdges {
		go d.Walk(downEdge)
	}
}



