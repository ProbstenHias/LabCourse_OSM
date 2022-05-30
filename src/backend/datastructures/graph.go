package datastructures

type Graph struct {
	// holds the nodes of the Graph with its lat and long coordinates
	Nodes [][]float64
	// edge list that indicate to which node an edge points, its distance is saved in the 'offset' array
	Edges []int32
	// distance of and edge
	Distance []int32
	// for each node in nodes it holds the offset where to look for its outgoing edges in the 'edges' array
	Offset []int32
}

func (g Graph) GetAllOutgoingEdgesOfNode(node int32) []int32 {
	var outEdges []int32
	var nextOffset = 0
	if node == int32(len(g.Nodes)-1) {
		nextOffset = len(g.Edges)
	} else {
		nextOffset = int(g.Offset[node+1])
	}
	for i := g.Offset[node]; i < int32(nextOffset); i++ {
		outEdges = append(outEdges, i)
	}
	return outEdges

}
