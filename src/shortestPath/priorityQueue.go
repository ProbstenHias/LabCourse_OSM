// This example demonstrates a dist queue built using the heap interface.
package shortestPath

import (
	"container/heap"
)

// An Item is something we manage in a dist queue.
type Item struct {
	id   int32   // The id of the item; arbitrary.
	dist float64 // The dist of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	prev  int32
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the lowest, not highest, dist, so we use greater than here.
	return pq[i].dist < pq[j].dist
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the dist and id of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, nodeId int32, dist float64) {
	item.id = nodeId
	item.dist = dist
	heap.Fix(pq, item.index)
}
