package datastructures

type Graph struct {
	// holds the nodes of the Graph with its lat and long coordinates
	Nodes [][]float64
	// edge list that indicate to which node an edge points, its distance is saved in the 'offset' array
	Edges []int
	// distance of and edge
	Distance []int
	// for each node in nodes it holds the offset where to look for its outgoing edges in the 'edges' array
	Offset []int
	// for contraction hierarchy this array saves the order in which nodes were contracted
	Order []int
	// saves information about which two edges make up a shortcut
	Shortcuts map[int]map[int]int
}

func (g Graph) GetAllOutgoingEdgesOfNode(node int) []int {
	var outEdges []int
	var nextOffset = 0
	if node == (len(g.Nodes) - 1) {
		nextOffset = len(g.Edges)
	} else {
		nextOffset = g.Offset[node+1]
	}
	for i := g.Offset[node]; i < nextOffset; i++ {
		outEdges = append(outEdges, i)
	}
	return outEdges

}

func (g Graph) GetAllOutgoingEdgeDistancesOfNode(node int) []int {
	var outDistances []int
	var nextOffset = 0
	if node == (len(g.Nodes) - 1) {
		nextOffset = len(g.Edges)
	} else {
		nextOffset = g.Offset[node+1]
	}
	for i := g.Offset[node]; i < nextOffset; i++ {
		outDistances = append(outDistances, g.Distance[i])
	}
	return outDistances
}

func (g Graph) GetAllNeighboursOfNode(node int) []int {
	edges := g.GetAllOutgoingEdgesOfNode(node)
	neighbours := make([]int, len(edges))
	for i, e := range edges {
		neighbours[i] = g.Edges[e]
	}
	return neighbours
}

func (g Graph) GetEdgeLength(from, to int) int {
	edges := g.GetAllOutgoingEdgesOfNode(from)
	for _, edge := range edges {
		if g.Edges[edge] == to {
			return g.Distance[edge]
		}
	}
	return -1
}

func (g Graph) ToCHGraph() CHGraph {
	chg := NewCHGraph()
	for i, node := range g.Nodes {
		_ = chg.AddNode(i, node[0], node[1])
		neighbours := g.GetAllNeighboursOfNode(i)
		distances := g.GetAllOutgoingEdgeDistancesOfNode(i)
		for j, neighbour := range neighbours {
			chg.AddEdge(i, neighbour, distances[j])
		}
	}
	chg.NodeOrder = g.Order
	return chg
}

func (g Graph) CombineEdges(addedEdges map[int][][]int) Graph {
	var totalEdgesAdded = 0
	var newEdges []int
	var newDistances []int
	for i := 0; i < len(g.Nodes); i++ {
		oldEdges := g.GetAllNeighboursOfNode(i)

		oldDistances := g.GetAllOutgoingEdgeDistancesOfNode(i)
		newEdges = append(newEdges, oldEdges...)
		newDistances = append(newDistances, oldDistances...)
		g.Offset[i] += totalEdgesAdded
		toAdd, ok := addedEdges[i]
		if !ok {
			continue
		}
		totalEdgesAdded += len(toAdd)

		for _, item := range toAdd {
			newEdges = append(newEdges, item[0])
			newDistances = append(newDistances, item[1])
		}
	}
	g.Edges = newEdges
	g.Distance = newDistances
	return g
}

func (g Graph) AddShortcuts(shortcuts map[int][][]int) Graph {
	sc := make(map[int]map[int]int)
	for from, ar := range shortcuts {
		sc[from] = make(map[int]int)
		for _, e := range ar {
			to := e[0]
			via := e[2]
			sc[from][to] = via
		}
	}
	g.Shortcuts = sc
	return g
}
