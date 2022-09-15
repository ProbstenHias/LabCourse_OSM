package shortestPath

import (
	"OSM/src/backend/helpers"
	"testing"
)

func TestBiDijkstra(t *testing.T) {
	graph := helpers.CreateGraphFromFile("../../../out/graph1m.fmi")
	randomRoutes := helpers.CreateRandomRoutes(10, len(graph.Nodes))
	for i := 0; i < len(randomRoutes); i++ {
		dijkstraLength, _, _ := Dijkstra(randomRoutes[i][0], randomRoutes[i][1], graph)
		biDijkstraLength, _, _ := BiDijkstra(randomRoutes[i][0], randomRoutes[i][1], graph)
		if dijkstraLength != biDijkstraLength {
			t.Errorf("The shortest path length of Bidijkstra is not equal to Dijkstra.\n"+
				"Routet from node %d to %d.\n"+
				"Got a length of %d but should have been %d", randomRoutes[i][0], randomRoutes[i][1], biDijkstraLength, dijkstraLength)

		}
	}
}
