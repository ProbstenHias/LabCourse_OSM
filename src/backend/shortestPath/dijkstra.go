package shortestPath

import (
	"OSM/src/backend/datastructures"
	"container/heap"
	"log"
	"math"
	"time"
)

func Dijkstra(start int, end int, graph datastructures.Graph) (int, []int, int) {
	var numberOfHeapPulls = 0
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
	})
	for pq.Len() > 0 {
		numberOfHeapPulls++
		node := heap.Pop(&pq).(*datastructures.Item)
		if node.Prio >= dist[node.Id] {
			continue
		}
		dist[node.Id] = node.Prio
		prev[node.Id] = node.Prev

		if node.Id == end {
			log.Printf("Time to calculate dijkstra: %s\n", time.Since(startTime))
			return node.Prio, prev, numberOfHeapPulls
		}
		for _, e := range graph.GetAllOutgoingEdgesOfNode(node.Id) {
			var to = graph.Edges[e]
			var alt = node.Prio + graph.Distance[e]
			if alt >= dist[to] {
				continue
			}
			heap.Push(&pq, &datastructures.Item{
				Id:   to,
				Prio: alt,
				Prev: node.Id,
			})
		}
	}
	log.Printf("No path from %d to %d could be found.\n", start, end)
	return dist[end], prev, numberOfHeapPulls
}

func DijkstraOneToN(start int, end []int, graph datastructures.Graph) ([]int, []int, int) {

	endLookup := make(map[int]bool)
	for i := 0; i < len(end); i++ {
		endLookup[end[i]] = false
	}
	var numberOfHeapPulls = 0
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
	})
	for pq.Len() > 0 {
		numberOfHeapPulls++
		node := heap.Pop(&pq).(*datastructures.Item)
		if node.Prio >= dist[node.Id] {
			continue
		}
		dist[node.Id] = node.Prio
		prev[node.Id] = node.Prev

		_, ok := endLookup[node.Id]
		if ok {
			endLookup[node.Id] = true
		}
		if allVisited(endLookup) {
			return ConstructReturnDists(end, dist), prev, numberOfHeapPulls
		}
		for _, e := range graph.GetAllOutgoingEdgesOfNode(node.Id) {
			var to = graph.Edges[e]
			var alt = node.Prio + graph.Distance[e]
			if alt >= dist[to] {
				continue
			}
			heap.Push(&pq, &datastructures.Item{
				Id:   to,
				Prio: alt,
				Prev: node.Id,
			})
		}
	}
	return ConstructReturnDists(end, dist), prev, numberOfHeapPulls
}

func allVisited(v map[int]bool) bool {
	var returnValue = true
	for _, value := range v {
		returnValue = returnValue && value
	}
	return returnValue
}

func ConstructReturnDists(end []int, dist []int) []int {
	returnDists := make([]int, len(end))
	for i := 0; i < len(end); i++ {
		returnDists[i] = dist[end[i]]
	}
	return returnDists
}
