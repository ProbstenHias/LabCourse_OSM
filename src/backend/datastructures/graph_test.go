package datastructures

import (
	"reflect"
	"testing"
)

func buildTestGraph() Graph {
	graph := Graph{}
	// coordinates of nodes do not matter just make 20 nodes
	graph.Nodes = [][]float64{{1, 2}, {3, 1}, {5, 1}, {1, 6}, {1, 1}, {1, 1}, {8, 1}, {1, 1}, {9, 1}, {1, 1}, {534, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}
	graph.Edges = []int{1, 5,
		0, 6, 2,
		1, 7, 3,
		2, 8, 4,
		3, 9,
		0, 6, 10,
		1, 5, 7, 11,
		2, 6, 8, 12,
		3, 7, 9, 13,
		4, 8, 14,
		5, 11, 15,
		6, 10, 12, 16,
		7, 11, 13, 17,
		8, 12, 14, 18,
		9, 13, 19,
		10, 16,
		11, 15, 17,
		12, 16, 18,
		13, 17, 19,
		14, 18}
	graph.Offset = []int{0, 2, 5, 8, 11, 13, 16, 20, 24, 28, 31, 34, 38, 42, 46, 48, 51, 54, 57, 60, 62}
	graph.Distance = []int{}
	for i := 0; i < len(graph.Edges); i++ {
		graph.Distance = append(graph.Distance, i)
	}
	return graph
}

func TestGraph_ToCHGraph(t *testing.T) {
	graph := buildTestGraph()
	chg := graph.ToCHGraph()
	graph2 := chg.ToGraph()
	if !reflect.DeepEqual(graph, graph2) {
		t.Errorf("The conversion between graphs types did not work correctly.")
	}
}
