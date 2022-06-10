package shortestPath

import (
	datastructures2 "OSM/src/backend/datastructures"
	"container/heap"
	"log"
	"math"
	"time"
)

func Dijkstra(start int32, end int32, graph datastructures2.Graph) (int32, []int32) {
	startTime := time.Now()
	dist := make([]int32, len(graph.Nodes))
	for i := 0; i < len(dist); i++ {
		dist[i] = math.MaxInt32
	}
	prev := make([]int32, len(graph.Nodes))
	pq := make(datastructures2.PriorityQueue, 0)

	heap.Push(&pq, &datastructures2.Item{
		Id:   start,
		Dist: 0,
		Prev: start,
	})
	for pq.Len() > 0 {
		node := heap.Pop(&pq).(*datastructures2.Item)
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
			var to = graph.Edges[e]
			var alt = node.Dist + graph.Distance[e]
			if alt >= dist[to] {
				continue
			}
			heap.Push(&pq, &datastructures2.Item{
				Id:   to,
				Dist: alt,
				Prev: node.Id,
			})
		}
	}
	log.Printf("No path from %d to %d could be found.\n", start, end)
	return dist[end], prev
}
