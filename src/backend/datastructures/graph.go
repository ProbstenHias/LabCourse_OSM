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
