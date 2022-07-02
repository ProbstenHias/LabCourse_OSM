package shortestPath

import (
	"OSM/src/backend/helpers"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestDijkstra(t *testing.T) {

	graph := helpers.CreateGraphFromFile("D:/OneDrive - stud.uni-stuttgart.de/Uni/10. Semester/FP-OSM/pbf files/oceanfmi.sec")
	log.Println("Generated graph")
	var times []time.Duration
	for i := 0; i < 10; i++ {
		randomStart := rand.Intn(len(graph.Nodes) - 1)
		randomDest := rand.Intn(len(graph.Nodes) - 1)
		start := time.Now()
		Dijkstra(randomStart, randomDest, graph)
		diff := time.Since(start)
		times = append(times, diff)
		log.Printf("Dijkstra took %s seconds\n", diff)
	}
}
