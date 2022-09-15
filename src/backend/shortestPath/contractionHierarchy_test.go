package shortestPath

import (
	"OSM/src/backend/datastructures"
	"OSM/src/backend/helpers"
	"fmt"
	"log"
	"testing"
)

func TestContractionHierarchy_ContractGraph(t *testing.T) {
	//build a grid graph and test contraction

	graph := datastructures.Graph{}
	// coordinates of nodes do not matter just make 20 nodes
	graph.Nodes = [][]float64{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}
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
		graph.Distance = append(graph.Distance, 1)
	}

	// test if graph was created correctly
	if graph.GetAllNeighboursOfNode(12)[1] != 11 {
		t.Errorf("The graph was not set up correctly.\n Neighbour of 12 should have been 11 but was %d.\n", graph.GetAllNeighboursOfNode(12)[1])
		return
	}
	if graph.GetAllNeighboursOfNode(19)[0] != 14 {
		log.Println(graph.GetAllNeighboursOfNode(19))

		t.Errorf("The graph was not set up correctly.\n Neighbour of 19 should have been 14 but was %d.\n", graph.GetAllNeighboursOfNode(19)[0])
		return

	}
	t.Log("Graph was set up correctly")

	graph = ContractGraph(graph)

}

func TestFindPath(t *testing.T) {
	pathToFMIFile := "../../../out/graph100k.fmi"
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

func TestSetUpAllDestinations(t *testing.T) {
	graph := helpers.CreateGraphFromFile("../../../out/graph1h.fmi")
	chg := graph.ToCHGraph()
	destinationNodes := []int{0, 1, 2, 50, 40}
	dests := setUpBuckets(destinationNodes, chg)
	fmt.Println(dests)
}
