package helpers

import (
	"OSM/src/backend/datastructures"
)

func CreateCoordinatesPathFromPrev(start, end int, prev []int, graph datastructures.Graph) [][]float64 {
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

func CreateIndexPathFromPrev(start, end int, prev []int) []int {
	var path []int
	currNode := end
	for {
		path = append(path, currNode)
		if currNode == start {
			return path
		}
		currNode = prev[currNode]
	}
}
