package shortestPath

import (
	datastructures2 "OSM/src/backend/datastructures"
	"container/heap"
	"log"
	"math"
	"time"
)

func Dijkstra(start int32, end int32, graph datastructures2.Graph) (int32, []int32, int) {
	var numberOfHeapPulls = 0
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
	})
	for pq.Len() > 0 {
		numberOfHeapPulls++
		node := heap.Pop(&pq).(*datastructures2.Item)
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
			heap.Push(&pq, &datastructures2.Item{
				Id:   to,
				Prio: alt,
				Prev: node.Id,
			})
		}
	}
	log.Printf("No path from %d to %d could be found.\n", start, end)
	return dist[end], prev, numberOfHeapPulls
}

func DijkstraOneToN(start int32, end []int32, graph datastructures2.Graph) ([]int32, []int32, int) {

	endLookup := make(map[int32]bool)
	for i := 0; i < len(end); i++ {
		endLookup[end[i]] = false
	}
	var numberOfHeapPulls = 0
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
	})
	for pq.Len() > 0 {
		numberOfHeapPulls++
		node := heap.Pop(&pq).(*datastructures2.Item)
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
			return constructReturnDists(end, dist), prev, numberOfHeapPulls
		}
		for _, e := range graph.GetAllOutgoingEdgesOfNode(node.Id) {
			var to = graph.Edges[e]
			var alt = node.Prio + graph.Distance[e]
			if alt >= dist[to] {
				continue
			}
			heap.Push(&pq, &datastructures2.Item{
				Id:   to,
				Prio: alt,
				Prev: node.Id,
			})
		}
	}
	return constructReturnDists(end, dist), prev, numberOfHeapPulls
}

func allVisited(v map[int32]bool) bool {
	var returnValue = true
	for _, value := range v {
		returnValue = returnValue && value
	}
	return returnValue
}

func constructReturnDists(end []int32, dist []int32) []int32 {
	returnDists := make([]int32, len(end))
	for i := 0; i < len(end); i++ {
		returnDists[i] = dist[end[i]]
	}
	return returnDists
}
