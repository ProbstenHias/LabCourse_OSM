package datastructures

import (
	"container/heap"
	"math/rand"
	"testing"
	"time"
)

func TestPriorityQueue_Pop(t *testing.T) {
	rand.Seed(time.Now().Unix())
	input := rand.Perm(1000)
	pq := make(PriorityQueue, 0)
	for i := 0; i < len(input); i++ {
		heap.Push(&pq, &Item{
			Id:   input[i],
			Prio: input[i],
		})
	}
	var lastElem = -1
	for i := 0; i < len(input); i++ {
		elem := heap.Pop(&pq).(*Item)
		if lastElem > elem.Prio {
			t.Errorf("Failed! The pq did not order the elements right")
			return
		}
		lastElem = elem.Prio
	}
	t.Log("PQ was Successful")
}
