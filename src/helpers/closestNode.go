package helpers

import (
	"OSM/src/datastructures"
	"log"
	"math"
	"time"
)

func GetClosestNodeInGraph(node []float64, graph datastructures.Graph) int32 {
	startTime := time.Now()
	minDistance := math.MaxInt32
	var minIdx = 0
	for i, elem := range graph.Nodes {
		var distance = Haversine(node, elem)
		if distance >= minDistance {
			continue
		}
		minDistance = distance
		minIdx = i
	}
	log.Printf("Time to calculate closest node: %s\n", time.Since(startTime))
	return int32(minIdx)
}
