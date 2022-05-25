package helpers

import (
	"OSM/src/datastructures"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func CreateGraphFromFile(pathToFile string) datastructures.Graph {
	readFile, err := os.Open(pathToFile)
	if err != nil {
		fmt.Println(err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	//skip comments
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if line == "" {
			break
		}
	}
	fileScanner.Scan()
	numberNodes, _ := strconv.Atoi(fileScanner.Text())
	fileScanner.Scan()
	numberEdges, _ := strconv.Atoi(fileScanner.Text())

	graph := datastructures.Graph{
		Nodes:    make([][]float64, numberNodes),
		Edges:    make([]int32, numberEdges),
		Distance: make([]int32, numberEdges),
		Offset:   make([]int32, numberNodes),
	}
	idToIdx := make(map[int]int)
	createNodesFromFile(numberNodes, graph, idToIdx, fileScanner)
	createEdgesFromFile(numberEdges, graph, idToIdx, fileScanner)

	readFile.Close()
	return graph
}

func createFileFromGraph(graph datastructures.Graph, pathToFile string) {

}

func createNodesFromFile(nodeCount int, graph datastructures.Graph, idToIdx map[int]int, fileScanner *bufio.Scanner) {
	for i := 0; i < nodeCount; i++ {
		fileScanner.Scan()
		line := strings.Split(fileScanner.Text(), " ")
		id, _ := strconv.Atoi(line[0])
		lat, _ := strconv.ParseFloat(line[1], 64)
		long, _ := strconv.ParseFloat(line[2], 64)
		node := []float64{lat, long}
		graph.Nodes[i] = node
		idToIdx[id] = i
	}
}

func createEdgesFromFile(edgeCount int, graph datastructures.Graph, idToIdx map[int]int, fileScanner *bufio.Scanner) {
	graph.Offset[0] = 0
	var from = 0
	for i := 0; i < edgeCount; i++ {
		fileScanner.Scan()
		line := strings.Split(fileScanner.Text(), " ")
		fr, _ := strconv.Atoi(line[0])
		fr = idToIdx[fr]
		to, _ := strconv.Atoi(line[1])
		to = idToIdx[to]
		dist, _ := strconv.Atoi(line[2])
		graph.Edges[i] = int32(to)
		graph.Distance[i] = int32(dist)
		if fr == from {
			continue
		}
		graph.Offset[fr] = int32(i)
		from = fr
	}
}
