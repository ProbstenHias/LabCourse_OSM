package datastructures

import (
	"errors"
)

type edge struct {
	From     int
	To       int
	Distance int
}

type Node struct {
	Id            int
	Lat           float64
	Long          float64
	OutboundEdges map[int]edge
}

type CHGraph struct {
	Nodes         map[int]Node
	NodeOrder     []int
	NumberOfEdges int
}

func NewCHGraph() CHGraph {
	return CHGraph{map[int]Node{}, []int{}, 0}
}

func (chg *CHGraph) AddNode(id int, lat, long float64) error {
	_, ok := chg.Nodes[id]
	if ok {
		return errors.New("this node already exists")
	}
	chg.Nodes[id] = Node{
		Id:            id,
		Lat:           lat,
		Long:          long,
		OutboundEdges: map[int]edge{},
	}
	return nil
}

func (chg *CHGraph) AddEdge(from, to, distance int) {
	if entry, ok := chg.Nodes[from]; ok {

		entry.OutboundEdges[to] = edge{
			From:     from,
			To:       to,
			Distance: distance,
		}

		chg.Nodes[from] = entry
		chg.NumberOfEdges++
	}
}

func (chg *CHGraph) ToGraph() Graph {
	g := Graph{
		Nodes:    make([][]float64, len(chg.Nodes)),
		Edges:    make([]int, 0),
		Distance: make([]int, 0),
		Offset:   make([]int, len(chg.Nodes)+1),
		Order:    nil,
	}
	lastOffset := 0
	for i := 0; i < len(chg.Nodes); i++ {
		currNode := chg.Nodes[i]
		g.Nodes[i] = []float64{currNode.Lat, currNode.Long}
		g.Offset[i] = lastOffset
		for _, e := range currNode.OutboundEdges {
			g.Edges = append(g.Edges, e.To)
			g.Distance = append(g.Distance, e.Distance)
			lastOffset++
		}
	}
	g.Offset[len(chg.Nodes)] = lastOffset
	g.Order = chg.NodeOrder
	return g

}
func (chg *CHGraph) GetAllNeighboursOfNode(node int) []int {
	neighbours := make([]int, 0, len(chg.Nodes[node].OutboundEdges))
	for to := range chg.Nodes[node].OutboundEdges {
		neighbours = append(neighbours, to)
	}
	return neighbours
}

func (chg *CHGraph) GetEdgeLength(from int, to int) int {
	if entry, ok := chg.Nodes[from].OutboundEdges[to]; ok {
		return entry.Distance
	}
	return -1
}

func (chg *CHGraph) DeleteNode(node int) {
	for key := range chg.Nodes[node].OutboundEdges {
		delete(chg.Nodes[key].OutboundEdges, node)
		chg.NumberOfEdges--
	}
	chg.NumberOfEdges -= len(chg.Nodes[node].OutboundEdges)
	delete(chg.Nodes, node)
}
