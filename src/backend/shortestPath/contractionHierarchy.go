package shortestPath

import (
	"OSM/src/backend/datastructures"
	"container/heap"
	"log"
	"math"
	"time"
)

func ContractGraph(graph datastructures.Graph) datastructures.Graph {

	// this is the order in which nodes were contracted
	chg := graph.ToCHGraph()
	startTime := time.Now()
	nodeOrder := make([]int, len(chg.Nodes))
	for i := 0; i < len(chg.Nodes); i++ {
		nodeOrder[i] = -1
	}
	// here we hold the shortcut edges of a node
	// looking up a node in the map returns a slice of all shortcut
	// a shortcut is an array where shortcut[0] returns the dest node and shortcut[1] holds the corresponding distance
	shortcuts := make(map[int][][]int)

	// this queue holds the nodes that are yet to be contracted and sorts them according to the heuristic
	pq := make(datastructures.PriorityQueue, 0)
	for i := 0; i < len(chg.Nodes); i++ {
		prio := calcPrioOfNode(i, shortcuts, chg)
		heap.Push(&pq, &datastructures.Item{
			Id:   i,
			Prio: prio,
		})
	}

	orderIndex := 0
	log.Println("Starting to contract nodes.")

	for pq.Len() > 0 {
		currNode := getNextNode(&pq, shortcuts, chg)
		nodeOrder[currNode.Id] = orderIndex
		orderIndex++

		directNeighbours := chg.GetAllNeighboursOfNode(currNode.Id)
		printStatus(orderIndex, directNeighbours, chg)

		// in this loop from each neighbour to all other neighbours of currNode shortcuts are inserted if and only if the direct path
		// via currNode is the shortest path from one neighbour to another
		for k := 0; k < len(directNeighbours)-1; k++ {
			from := directNeighbours[k]

			// exclude all path searches between nodes that already have a direct edge
			destinationNodes := removeExistingPaths(directNeighbours, k, chg, from)

			pathLengths := calcDirectPathLengthsOneToN(from, destinationNodes, currNode.Id, chg)

			// restrict the search space of the dijkstra
			maxPathLength := calcDijkstraLimit(from, currNode.Id, destinationNodes, chg)

			// use buckets as described in the paper
			buckets := setUpBuckets(destinationNodes, chg)

			dist, _, prioQueue := setUpDijkstra(from, chg)
			shortcutable := make(map[int]bool)
			for _, node := range destinationNodes {
				shortcutable[node] = true
			}

			for true {
				// when we reach a node x we scan its bucket entries
				// for each entry we know there is a path from u to w with length c(u, x) + c(x,w)
				// with that info we make shortcuts or not
				x := limitedDijkstra(buckets, maxPathLength, dist, &prioQueue, nodeOrder, chg)
				handleFoundNode(from, x, dist, pathLengths, buckets, destinationNodes, shortcuts, currNode, shortcutable, chg)
				if x == -1 {
					break
				}
				shouldContinue := false
				for _, value := range shortcutable {
					if value {
						shouldContinue = true
						break
					}
				}
				if !shouldContinue {
					break
				}
			}
		}
		chg.DeleteNode(currNode.Id)

	}
	graph = graph.CombineEdges(shortcuts)
	chg = graph.ToCHGraph()
	chg.NodeOrder = nodeOrder
	chg = removeUnnecessaryEdges(chg)
	graph = chg.ToGraph()

	graph = graph.AddShortcuts(shortcuts)
	graph.Order = nodeOrder
	log.Printf("Time to contract the graph: %s", time.Since(startTime))

	return graph

}

func handleFoundNode(from, x int, dist, pathLengths map[int]int, buckets map[int][][]int, destinationNodes []int, shortcuts map[int][][]int, currNode *datastructures.Item, shortcutable map[int]bool, chg datastructures.CHGraph) {
	if x != -1 {
		// we iterate over all bucket entries belonging to x
		for _, bucket := range buckets[x] {
			w := bucket[1]
			if value, _ := shortcutable[w]; !value {
				continue
			}
			// here we get the path length that was calculated
			foundPathLength := dist[x] + bucket[0]

			// we know there is a path to w with the found length
			// if this length is shorter than the direct path via currNode we know there will never be a shortcut
			if pathLengths[w] >= foundPathLength {
				shortcutable[w] = false
			}

		}
	}

	// if we reached no more nodes that means for the remaining nodes we can add shortcuts
	if x == -1 {
		// all shortcuts can be set since no path could be found between the remaining nodes
		for _, node := range destinationNodes {
			// if node was declared non shortcutable skip it
			if value, _ := shortcutable[node]; !value {
				continue
			}

			shortcuts[from] = append(shortcuts[from], []int{node, pathLengths[node], currNode.Id})
			chg.AddEdge(from, node, pathLengths[node])
			shortcuts[node] = append(shortcuts[node], []int{from, pathLengths[node], currNode.Id})
			chg.AddEdge(node, from, pathLengths[node])
		}
	}
}

func printStatus(nodesContracted int, directNeighbours []int, chg datastructures.CHGraph) {
	if nodesContracted%1000 == 0 {
		log.Printf("Contracted %d nodes\n", nodesContracted)
		log.Printf("Neighbours %d\n", len(directNeighbours))
		log.Printf("Nodes left: %d", len(chg.Nodes))
	}
}

func getNextNode(priorityQueue *datastructures.PriorityQueue, shortcuts map[int][][]int, chg datastructures.CHGraph) *datastructures.Item {
	// this method returns a new node from the priority queue
	// this method uses lazy update to make sure the priority of each node is up to date
	currNode := heap.Pop(priorityQueue).(*datastructures.Item)

	for currNode.Prio != calcPrioOfNode(currNode.Id, shortcuts, chg) {
		heap.Push(priorityQueue, &datastructures.Item{
			Id:   currNode.Id,
			Prio: calcPrioOfNode(currNode.Id, shortcuts, chg),
		})
		currNode = heap.Pop(priorityQueue).(*datastructures.Item)
	}
	return currNode
}

func removeExistingPaths(directNeighbours []int, k int, chg datastructures.CHGraph, from int) []int {
	destinationNodes := make([]int, 0, len(directNeighbours[k+1:]))
	for _, neighbour := range directNeighbours[k+1:] {
		if chg.GetEdgeLength(from, neighbour) != -1 {
			continue
		}
		destinationNodes = append(destinationNodes, neighbour)
	}
	return destinationNodes
}

func calcDijkstraLimit(from, via int, destinationNodes []int, chg datastructures.CHGraph) int {

	// the limit of a dijkstra run is c(from, via) + max_w ( c(v, w)) - min_x ( c(x,w))
	// with x being all the direct neighbours of all destinationNodes w
	max := -1
	min := math.MaxInt
	for _, node := range destinationNodes {
		l := chg.GetEdgeLength(via, node)
		if l > max {
			max = l
		}
		for _, edge := range chg.Nodes[node].OutboundEdges {
			if edge.Distance < min {
				min = edge.Distance
			}
		}
	}
	return chg.GetEdgeLength(from, via) + max - min
}

func setUpBuckets(destinationNodes []int, chg datastructures.CHGraph) map[int][][]int {
	destinationsMap := make(map[int][][]int)
	for _, node := range destinationNodes {
		for _, edge := range chg.Nodes[node].OutboundEdges {
			entry, _ := destinationsMap[edge.To]
			entry = append(entry, []int{edge.Distance, node})
			destinationsMap[edge.To] = entry
		}
	}
	return destinationsMap
}

func calcDirectPathLengthsOneToN(from int, to []int, via int, chg datastructures.CHGraph) map[int]int {
	lengths := make(map[int]int, len(to))
	for _, element := range to {
		fromVia := chg.GetEdgeLength(from, via)
		viaTo := chg.GetEdgeLength(via, element)
		lengths[element] = fromVia + viaTo
	}
	return lengths
}

func limitedDijkstra(end map[int][][]int, maxDist int, dist map[int]int, pq *datastructures.PriorityQueue, nodeOrder []int, chg datastructures.CHGraph) int {
	for pq.Len() > 0 {
		node := heap.Pop(pq).(*datastructures.Item)
		if node.Prio > maxDist {
			return -1
		}
		if entry, ok := dist[node.Id]; ok && node.Prio >= entry {
			continue
		}
		dist[node.Id] = node.Prio

		for _, e := range chg.Nodes[node.Id].OutboundEdges {
			var to = e.To
			if nodeOrder[to] != -1 {
				continue
			}
			var alt = node.Prio + e.Distance

			if entry, exists := dist[to]; exists && alt >= entry {
				continue
			}
			heap.Push(pq, &datastructures.Item{
				Id:   to,
				Prio: alt,
			})
		}
		_, ok := end[node.Id]
		if ok {
			return node.Id
		}

	}
	return -1
}

func calcPrioOfNode(node int, shortcuts map[int][][]int, chg datastructures.CHGraph) int {
	// Worst possible shortcuts number through the vertex is: NumWorstShortcuts = NumIncomingEdges*NumOutcomingEdges
	shortcutCover := 2 * len(chg.Nodes[node].OutboundEdges)
	// Number of total incident edges is: NumIncomingEdges+NumOutcomingEdges
	incidentEdgesNum := 2 * len(chg.Nodes[node].OutboundEdges)
	// Edge difference is between NumWorstShortcuts and TotalIncidentEdgesNum
	edgeDiff := shortcutCover - incidentEdgesNum

	// Spatial diversity heuristic: for each vertex count the count of the number of neighbors that have already been contracted, and add this to the summary importance
	spacialDiv := calcSpatialDiversity(node, shortcuts)

	// Bidirection edges heuristic: for each vertex check how many bidirected incident edges vertex has. Sub that from importance
	//bidirection := len(chg.Nodes[node].OutboundEdges)

	importance := edgeDiff + incidentEdgesNum + spacialDiv
	return importance
}

func calcSpatialDiversity(node int, shortcuts map[int][][]int) int {
	delNeighbours := make(map[int]bool)
	for _, shortcut := range shortcuts[node] {
		delNeighbours[shortcut[2]] = true
	}
	return len(delNeighbours)
}

func CHDijkstra(start, end int, graph datastructures.Graph) (int, []int, int) {
	// setup a forward and a backward dijkstra
	startTime := time.Now()
	shortestPathLength := math.MaxInt
	shortestPathVia := -1
	var numberOfHeapPulls = 0
	abortFor := false
	abortBack := false
	distFor := make([]int, len(graph.Nodes))
	distBack := make([]int, len(graph.Nodes))
	prevFor := make([]int, len(graph.Nodes))
	prevBack := make([]int, len(graph.Nodes))
	for i := 0; i < len(distFor); i++ {
		distFor[i] = math.MaxInt
		distBack[i] = math.MaxInt
		prevFor[i] = -1
		prevBack[i] = -1

	}
	pqFor := make(datastructures.PriorityQueue, 0)
	pqBack := make(datastructures.PriorityQueue, 0)

	heap.Push(&pqFor, &datastructures.Item{
		Id:   start,
		Prio: 0,
		Prev: start,
	})
	heap.Push(&pqBack, &datastructures.Item{
		Id:   end,
		Prio: 0,
		Prev: end,
	})

	for !(abortFor && abortBack) {
		nodeFor, nodeBack := -1, -1

		// iterate between a forward and a backward dijkstra step
		if !abortFor {
			nodeFor, abortFor = dijkstraStep(end, &pqFor, distFor, prevFor, &numberOfHeapPulls, graph)
		}
		if !abortBack {
			nodeBack, abortBack = dijkstraStep(start, &pqBack, distBack, prevBack, &numberOfHeapPulls, graph)
		}

		if nodeFor != -1 {
			// if the settled node in forward direction has a greater distance to the start node than the current shortest path
			// we can abort the forward run
			if distFor[nodeFor] >= shortestPathLength {
				abortFor = true
			}
			// if we found a node that was also visited in the backward run we have a shortest path candidate
			if distBack[nodeFor] != math.MaxInt {
				pathLengthAlt := distFor[nodeFor] + distBack[nodeFor]
				// if the candidate is shorter than the current shortest path we update the curr shortest path
				if pathLengthAlt < shortestPathLength {
					shortestPathLength = pathLengthAlt
					shortestPathVia = nodeFor
				}
			}
		}
		// we do the same as in the forward run
		if nodeBack != -1 {
			if distBack[nodeBack] >= shortestPathLength {
				abortBack = true
			}
			if distFor[nodeBack] != math.MaxInt {
				pathLengthAlt := distFor[nodeBack] + distBack[nodeBack]
				if pathLengthAlt < shortestPathLength {
					shortestPathLength = pathLengthAlt
					shortestPathVia = nodeBack
				}
			}
		}

		nodeFor, nodeBack = -1, -1

	}
	// if the run found a path set up the prev array
	if shortestPathLength < math.MaxInt {

		prevNode := shortestPathVia
		for prevNode != end {
			prevFor[prevBack[prevNode]] = prevNode
			prevNode = prevBack[prevNode]
		}

		currNode := end
		for currNode != start {
			to := prevFor[currNode]
			via, ok := graph.Shortcuts[to][currNode]
			if ok {
				prevFor[currNode] = via
				prevFor[via] = to
			} else {
				currNode = to
			}

		}
	}
	log.Printf("Time to calculate CHDijkstra: %s", time.Since(startTime))
	return shortestPathLength, prevFor, numberOfHeapPulls

}

func dijkstraStep(end int, queue *datastructures.PriorityQueue, dist, prev []int, hp *int, graph datastructures.Graph) (int, bool) {
	// if there are no more entries in heap left, then abort execution
	if queue.Len() <= 0 {
		return -1, true
	}
	node := heap.Pop(queue).(*datastructures.Item)
	*hp++
	// skip all nodes that are where already visited
	for node.Prio >= dist[node.Id] && queue.Len() > 0 {
		node = heap.Pop(queue).(*datastructures.Item)
		*hp++
	}
	// if the curr node was already visited we abort
	if node.Prio >= dist[node.Id] {
		return -1, true
	}
	dist[node.Id] = node.Prio
	prev[node.Id] = node.Prev

	if node.Id == end {
		return node.Id, true
	}
	for _, e := range graph.GetAllOutgoingEdgesOfNode(node.Id) {
		var to = graph.Edges[e]

		if graph.Order[node.Id] >= graph.Order[to] {
			continue
		}

		var alt = node.Prio + graph.Distance[e]
		if alt >= dist[to] {
			continue
		}
		heap.Push(queue, &datastructures.Item{
			Id:   to,
			Prio: alt,
			Prev: node.Id,
		})
	}

	return node.Id, false
}

func setUpDijkstra(start int, chg datastructures.CHGraph) (map[int]int, []int, datastructures.PriorityQueue) {
	dist := map[int]int{}

	prev := make([]int, len(chg.Nodes))
	pq := make(datastructures.PriorityQueue, 0)
	heap.Push(&pq, &datastructures.Item{
		Id:   start,
		Prio: 0,
		Prev: start,
	})

	return dist, prev, pq

}

func removeUnnecessaryEdges(chg datastructures.CHGraph) datastructures.CHGraph {
	// routing only need edges that go from nodes with lower order to nodes with higher order
	// this method removes all other edges
	for _, node := range chg.Nodes {
		toDelete := make([]int, 0, len(node.OutboundEdges))
		for _, edge := range node.OutboundEdges {
			from := edge.From
			to := edge.To
			if chg.NodeOrder[to] < chg.NodeOrder[from] {
				toDelete = append(toDelete, to)
			}
		}
		for _, to := range toDelete {
			delete(node.OutboundEdges, to)
		}
	}
	return chg
}
