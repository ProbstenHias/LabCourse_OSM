package shortestPath

import (
	"OSM/src/backend/datastructures"
	"container/heap"
	"log"
	"math"
	"sync"
	"time"
)

var visited map[int32]bool

func dijkstra(start, end int32, prev, dist []int32, graph datastructures.Graph, common *int32, numberOfHeapPulls *int, wg *sync.WaitGroup, m *sync.Mutex) {
	defer wg.Done()
	for i := 0; i < len(dist); i++ {
		dist[i] = math.MaxInt32
	}
	pq := make(datastructures.PriorityQueue, 0)

	heap.Push(&pq, &datastructures.Item{
		Id:   start,
		Prio: 0,
		Prev: start,
	})
	for pq.Len() > 0 {
		*numberOfHeapPulls++
		node := heap.Pop(&pq).(*datastructures.Item)
		if node.Prio >= dist[node.Id] {
			continue
		}

		dist[node.Id] = node.Prio
		prev[node.Id] = node.Prev
		m.Lock()
		_, exists := visited[node.Id]
		if exists {
			*common = node.Id
			m.Unlock()
			return
		}
		visited[node.Id] = true
		m.Unlock()
		if node.Id == end {
			*common = node.Id
			return
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
}

func BiDijkstra(start int32, end int32, graph datastructures.Graph) (int32, []int32, int) {
	startTime := time.Now()

	var heapPullsFor = 0
	var heapPullsBack = 0
	visited = make(map[int32]bool)
	prevFor := make([]int32, len(graph.Nodes))
	prevBack := make([]int32, len(graph.Nodes))
	distFor := make([]int32, len(graph.Nodes))
	distBack := make([]int32, len(graph.Nodes))
	var commonFor int32 = -1
	var commonBack int32 = -1
	var wg sync.WaitGroup
	var m sync.Mutex
	wg.Add(1)
	go dijkstra(start, end, prevFor, distFor, graph, &commonFor, &heapPullsFor, &wg, &m)
	wg.Add(1)
	go dijkstra(end, start, prevBack, distBack, graph, &commonBack, &heapPullsBack, &wg, &m)
	wg.Wait()

	//if commonFor == -1 || commonBack == -1 {
	//	return math.MaxInt32, prevFor
	//}
	//dist1 := distFor[commonFor] + distBack[commonFor]
	//dist2 := distFor[commonBack] + distBack[commonBack]
	//var common int32 = -1
	//var dist int32 = -1
	//if dist1 < dist2 {
	//	common = commonFor
	//	dist = dist1
	//
	//} else {
	//	common = commonBack
	//	dist = dist2
	//}
	//prev := generatePrev(start, end, prevFor, prevBack, common)
	log.Printf("Time to calculate bidijkstra: %s\n", time.Since(startTime))
	return 0, []int32{}, heapPullsFor + heapPullsBack

}

func generatePrev(start, end int32, prevFor, prevBack []int32, common int32) []int32 {
	prev := make([]int32, len(prevBack))
	var currIdx = start
	for currIdx != common {
		prev[currIdx] = prevFor[currIdx]
		currIdx = prevFor[currIdx]
	}
	prev[common] = prevFor[common]
	for currIdx != end {
		prev[prevBack[currIdx]] = currIdx
		currIdx = prevBack[currIdx]
	}
	return prev
}
