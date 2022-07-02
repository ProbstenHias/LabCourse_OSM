package datastructures

// An Item is something we manage in a Prio queue.
type Item struct {
	Id   int // The Id of the item; arbitrary.
	Prio int // The Prio of the item in the queue.
	// The Index is needed by update and is maintained by the heap.Interface methods.
	Prev  int
	Index int // The Index of the item in the heap.
	Dist  int // the distance we need for a star
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the lowest, not highest, Prio, so we use greater than here.
	return pq[i].Prio < pq[j].Prio
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
