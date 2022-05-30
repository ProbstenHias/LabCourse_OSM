package helpers

import (
	"OSM/src/backend/datastructures"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func CreateGraphFromFile(pathToFile string) datastructures.Graph {
	readFile, err := os.Open(pathToFile)
	if err != nil {
		log.Println(err)
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

	err = readFile.Close()
	if err != nil {
		log.Fatalf("Got error while clsing a filereader. Err: %s", err.Error())
	}
	return graph
}

func CreateFileFromGraph(graph datastructures.Graph, pathToFile string) {
	f, _ := os.Create(pathToFile)
	w := bufio.NewWriter(f)
	linesToWrite := []string{"# a comment for good measure\n",
		"\n",
		fmt.Sprintln(len(graph.Nodes)),
		fmt.Sprintln(len(graph.Edges))}

	for _, line := range linesToWrite {
		_, err := w.WriteString(line)
		if err != nil {
			log.Fatalf("Got error while writing to a file. Err: %s", err.Error())
		}
	}

	for i := 0; i < len(graph.Nodes); i++ {
		line := fmt.Sprintf("%d %f %f\n", i, graph.Nodes[i][0], graph.Nodes[i][1])
		_, err := w.WriteString(line)
		if err != nil {
			log.Fatalf("Got error while writing to a file. Err: %s", err.Error())

		}
	}
	var currFrom = 0
	for i := 0; i < len(graph.Edges); i++ {
		for currFrom+1 < len(graph.Nodes) && int32(i) >= graph.Offset[currFrom+1] {
			currFrom++
		}
		line := fmt.Sprintf("%d %d %d\n", currFrom, graph.Edges[i], graph.Distance[i])
		_, err := w.WriteString(line)
		if err != nil {
			log.Fatalf("Got error while writing to a file. Err: %s", err.Error())
		}
	}
	err := w.Flush()
	if err != nil {
		log.Fatalf("Got error while flushing a filewriter. Err: %s", err.Error())
	}
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
		// for all nodes with no edges add offset as well
		for j := from + 1; j <= fr; j++ {
			graph.Offset[j] = int32(i)

		}
		from = fr
	}
}
