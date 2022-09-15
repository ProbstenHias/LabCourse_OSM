package shortestPath

import (
	"OSM/src/backend/helpers"
	"fmt"
	"os"
	"testing"
)

func TestAStar(t *testing.T) {
	fmt.Println(os.Getwd())
	graph := helpers.CreateGraphFromFile("../../../out/graph1m.fmi")
	randomRoutes := helpers.CreateRandomRoutes(100, len(graph.Nodes))
	for i := 0; i < len(randomRoutes); i++ {
		dijkstraLength, _, _ := Dijkstra(randomRoutes[i][0], randomRoutes[i][1], graph)
		aStarLength, _, _ := AStar(randomRoutes[i][0], randomRoutes[i][1], graph)
		if dijkstraLength != aStarLength {
			t.Errorf("The shortest path length of AStar is not equal to Dijkstra.\n"+
				"Routet from node %d to %d.\n"+
				"Got a length of %d but should have been %d", randomRoutes[i][0], randomRoutes[i][1], aStarLength, dijkstraLength)
			return
		}
	}
}

func TestSpecificNode(t *testing.T) {
	graph := helpers.CreateGraphFromFile("../../../out/graph100k.fmi")
	start := 28565
	end := 46981
	dijkstraL, dijkstraPrev, _ := Dijkstra(start, end, graph)
	astarL, astarPrev, _ := AStar(start, end, graph)
	dijkstraRoute := helpers.CreateIndexPathFromPrev(start, end, dijkstraPrev)
	astarRoute := helpers.CreateIndexPathFromPrev(start, end, astarPrev)
	for i := 0; i < len(dijkstraRoute); i++ {
		if dijkstraRoute[i] != astarRoute[i] {
			t.Errorf("Dijkstra calculated a different route than Astar\n"+
				" %d vs. %d\n"+
				"AStar Length: %d vs. Dijkstra Length: %d\n", astarRoute[i], dijkstraRoute[i], astarL, dijkstraL)
			return

		}
	}

}
