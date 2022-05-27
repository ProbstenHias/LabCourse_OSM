package shortestPath

import (
	"OSM/src/datastructures"
	"container/heap"
	"log"
	"math"
	"time"
)

func Dijkstra(start int32, end int32, graph datastructures.Graph) (int32, []int32) {
	startTime := time.Now()
	dist := make([]int32, len(graph.Nodes))
	for i := 0; i < len(dist); i++ {
		dist[i] = math.MaxInt32
	}
	prev := make([]int32, len(graph.Nodes))
	pq := make(datastructures.PriorityQueue, 0)

	heap.Push(&pq, &datastructures.Item{
		Id:   start,
		Dist: 0,
		Prev: start,
	})
	for pq.Len() > 0 {
		node := heap.Pop(&pq).(*datastructures.Item)
		if node.Dist >= dist[node.Id] {
			continue
		}
		dist[node.Id] = node.Dist
		prev[node.Id] = node.Prev

		if node.Id == end {
			log.Printf("Time to calculate dijkstra: %s\n", time.Since(startTime))
			return node.Dist, prev
		}
		for _, e := range graph.GetAllOutgoingEdgesOfNode(node.Id) {

			var alt = node.Dist + graph.Distance[e]
			if alt >= dist[graph.Edges[e]] {
				continue
			}
			heap.Push(&pq, &datastructures.Item{
				Id:   graph.Edges[e],
				Dist: alt,
				Prev: node.Id,
			})
		}
	}
	log.Printf("No path from %d to %d could be found.\n", start, end)
	return dist[end], prev
}
