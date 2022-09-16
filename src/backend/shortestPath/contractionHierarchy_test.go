package shortestPath

import (
	"OSM/src/backend/helpers"
	"testing"
)

func TestFindPath(t *testing.T) {
	pathToFMIFile := "../../../out/graph1k.fmi"
	graph := helpers.CreateGraphFromFile(pathToFMIFile)
	graphContracted := helpers.CreateGraphFromFile(pathToFMIFile)
	graphContracted = ContractGraph(graphContracted)
	randomRoutes := helpers.CreateRandomRoutes(100, len(graph.Nodes))
	for i := 0; i < len(randomRoutes); i++ {
		dijkstraLength, _, _ := Dijkstra(randomRoutes[i][0], randomRoutes[i][1], graph)
		contractionLength, _, _ := CHDijkstra(randomRoutes[i][0], randomRoutes[i][1], graphContracted)
		if dijkstraLength != contractionLength {
			t.Errorf("The shortest path length of CH is not equal to Dijkstra.\n"+
				"Routet from node %d to %d.\n"+
				"Got a length of %d but should have been %d", randomRoutes[i][0], randomRoutes[i][1], contractionLength, dijkstraLength)
			return
		}

	}
}
