package shortestPath

import (
	"OSM/src/backend/datastructures"
	"OSM/src/backend/helpers"
	"container/heap"
	"log"
	"math"
	"time"
)

// try with caching heuristic results
// try 2d distance as heuristic

func manhattenDistance(point1, point2 []float64) int {
	return helpers.Haversine(point1, point2)
}

func AStar(start int, end int, graph datastructures.Graph) (int, []int, int) {
	var endCoordinates = graph.Nodes[end]
	var numberOfHeapPulls = 0
	heuristics := make(map[int]int)
	startTime := time.Now()
	dist := make([]int, len(graph.Nodes))
	for i := 0; i < len(dist); i++ {
		dist[i] = math.MaxInt
	}
	prev := make([]int, len(graph.Nodes))
	pq := make(datastructures.PriorityQueue, 0)
	heap.Push(&pq, &datastructures.Item{
		Id:   start,
		Prio: 0,
		Prev: start,
		Dist: 0,
	})
	for pq.Len() > 0 {
		numberOfHeapPulls++
		node := heap.Pop(&pq).(*datastructures.Item)
		if node.Dist >= dist[node.Id] {
			continue
		}
		dist[node.Id] = node.Dist
		prev[node.Id] = node.Prev

		if node.Id == end {
			log.Printf("Time to calculate astar: %s\n", time.Since(startTime))
			return node.Dist, prev, numberOfHeapPulls
		}
		for _, e := range graph.GetAllOutgoingEdgesOfNode(node.Id) {
			var to = graph.Edges[e]
			var alt = node.Dist + graph.Distance[e]
			if alt >= dist[to] {
				continue
			}
			h, ok := heuristics[to]
			if !ok {
				h = manhattenDistance(graph.Nodes[to], endCoordinates)
				heuristics[to] = h
			}
			heap.Push(&pq, &datastructures.Item{
				Id:   to,
				Prio: alt + h,
				Prev: node.Id,
				Dist: alt,
			})
		}
	}
	log.Printf("No path from %d to %d could be found.\n", start, end)
	return dist[end], prev, numberOfHeapPulls
}
