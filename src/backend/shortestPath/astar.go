package shortestPath

import (
	datastructures2 "OSM/src/backend/datastructures"
	"OSM/src/backend/helpers"
	"container/heap"
	"log"
	"math"
	"time"
)

// try with caching heuristic results
// try 2d distance as heuristic

func manhattenDistance(point1, point2 []float64) int32 {
	return int32(helpers.Haversine(point1, point2))
}

func AStar(start int32, end int32, graph datastructures2.Graph) (int32, []int32, int) {
	var endCoordinates = graph.Nodes[end]
	var numberOfHeapPulls = 0
	heuristics := make(map[int32]int32)
	startTime := time.Now()
	dist := make([]int32, len(graph.Nodes))
	for i := 0; i < len(dist); i++ {
		dist[i] = math.MaxInt32
	}
	prev := make([]int32, len(graph.Nodes))
	pq := make(datastructures2.PriorityQueue, 0)
	heap.Push(&pq, &datastructures2.Item{
		Id:   start,
		Prio: 0,
		Prev: start,
		Dist: 0,
	})
	for pq.Len() > 0 {
		numberOfHeapPulls++
		node := heap.Pop(&pq).(*datastructures2.Item)
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
			heap.Push(&pq, &datastructures2.Item{
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
