package helpers

import (
	"OSM/src/backend/datastructures"
)

func CreatePathFromPrev(start, end int32, prev []int32, graph datastructures.Graph) [][]float64 {
	var path [][]float64
	currNode := end
	for {
		path = append(path, graph.Nodes[currNode])
		if currNode == start {
			return path
		}
		currNode = prev[currNode]
	}
}
