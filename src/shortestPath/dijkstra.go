package shortestPath

import (
	"container/heap"
	"log"
	"math"
)

func dijkstra(start int32, end int32, graph Graph) (float64, []int32) {
	dist := make([]float64, len(graph.Nodes))
	for i := 0; i < len(dist); i++ {
		dist[i] = math.MaxFloat64
	}
	prev := make([]int32, len(graph.Nodes))
	pq := make(PriorityQueue, 0)

	heap.Push(&pq, &Item{
		id:   start,
		dist: 0,
		prev: start,
	})
	for pq.Len() > 0 {
		node := heap.Pop(&pq).(*Item)
		if node.dist > dist[node.id] {
			continue
		}
		dist[node.id] = node.dist
		prev[node.id] = node.prev

		if node.id == end {
			return node.dist, prev
		}
		for _, e := range graph.getAllOutgoingEdgesOfNode(node.id) {
			var alt = node.dist + graph.Distance[e]
			if alt >= dist[graph.Edges[e]] {
				continue
			}
			heap.Push(&pq, &Item{
				id:   graph.Edges[e],
				dist: alt,
				prev: node.id,
			})
		}
	}
	log.Fatalf("No path from %d to %d could be found.", start, end)

	return dist[end], prev
}
