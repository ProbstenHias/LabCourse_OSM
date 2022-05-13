package main

type Graph struct {
	// holds the nodes of the Graph with its lat and long coordinates
	nodes [][]float32
	// edge list that indicate to which node an edge points, its distance is saved in the 'offset' array
	edges []int32
	// distance of and edge
	distance []float64
	// for each node in nodes it holds the offset where to look for its outgoing edges in the 'edges' array
	offset []int32
}

func (g Graph) getAllOutgoingEdgesOfNode(node int32) []int32 {
	var outEdges []int32
	for i := g.offset[node]; i < g.offset[node+1]; i++ {
		outEdges = append(outEdges, i)
	}
	return outEdges

}
