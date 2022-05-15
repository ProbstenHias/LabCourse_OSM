package helpers

import (
	"OSM/src/shortestPath"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func CreateGraphFromFile(pathToFile string) shortestPath.Graph {
	readFile, err := os.Open(pathToFile)
	if err != nil {
		fmt.Println(err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	//skip comments
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if line[0] != '#' {
			break
		}
	}
	fileScanner.Scan()
	numberNodes, _ := strconv.Atoi(fileScanner.Text())
	fileScanner.Scan()
	numberEdges, _ := strconv.Atoi(fileScanner.Text())

	graph := shortestPath.Graph{
		Nodes:    make([][]float64, numberNodes),
		Edges:    make([]int32, numberEdges),
		Distance: make([]float64, numberEdges),
		Offset:   make([]int32, numberNodes),
	}

	createNodesFromFile(numberNodes, graph, fileScanner)
	createEdgesFromFile(numberEdges, graph, fileScanner)

	readFile.Close()
	return graph
}

func createFileFromGraph(graph shortestPath.Graph, pathToFile string) {

}

func createNodesFromFile(nodeCount int, graph shortestPath.Graph, fileScanner *bufio.Scanner) {
	for i := 0; i < nodeCount; i++ {
		fileScanner.Scan()
		line := strings.Split(fileScanner.Text(), " ")
		lat, _ := strconv.ParseFloat(line[1], 64)
		long, _ := strconv.ParseFloat(line[2], 64)
		node := []float64{lat, long}
		graph.Nodes[i] = node
	}
}

func createEdgesFromFile(edgeCount int, graph shortestPath.Graph, fileScanner *bufio.Scanner) {
	graph.Offset[0] = 0
	var from = 0
	for i := 0; i < edgeCount; i++ {
		fileScanner.Scan()
		line := strings.Split(fileScanner.Text(), " ")
		fr, _ := strconv.Atoi(line[0])
		to, _ := strconv.Atoi(line[1])
		dist, _ := strconv.ParseFloat(line[3], 64)
		graph.Edges[i] = int32(to)
		graph.Distance[i] = dist
		if fr == from {
			continue
		}
		graph.Offset[fr] = int32(i)
		from = fr
	}
}
