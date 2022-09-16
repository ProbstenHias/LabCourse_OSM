package main

import (
	"OSM/src/backend/datastructures"
	"OSM/src/backend/helpers"
	"OSM/src/backend/shortestPath"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

func main() {
	pathToFmiFile := os.Args[1]
	pathToCHFmiFile := os.Args[2]
	outPath := path.Dir(pathToFmiFile)
	const N = 100
	graph := helpers.CreateGraphFromFile(pathToFmiFile)
	chg := helpers.CreateContractedGraphFromFile(pathToCHFmiFile)
	randomRoutes := helpers.CreateRandomRoutes(N, len(graph.Nodes))

	log.Println("Starting to bench Dijkstra.")
	outPathDijkstra := fmt.Sprintf(outPath + "/dijkstra_bench.csv")
	benchShortestPathWithFunction(shortestPath.Dijkstra, graph, randomRoutes, outPathDijkstra)

	log.Println("Starting to bench A Star")
	outPathAStar := fmt.Sprintf(outPath + "/astar_bench.csv")
	benchShortestPathWithFunction(shortestPath.AStar, graph, randomRoutes, outPathAStar)

	log.Println("Starting to bench Contraction Hierarchy with Dijkstra")
	outPathCHDijkstra := fmt.Sprintf(outPath + "/chdijkstra_bench.csv")
	benchShortestPathWithFunction(shortestPath.CHDijkstra, chg, randomRoutes, outPathCHDijkstra)

}

func benchShortestPathWithFunction(shortestPathFunction func(int, int, datastructures.Graph) (int, []int, int), graph datastructures.Graph, randomRoutes [][]int, outPath string) {
	var timesAndPulls [][]string
	for i := 0; i < len(randomRoutes); i++ {
		randomStart := randomRoutes[i][0]
		randomDest := randomRoutes[i][1]
		start := time.Now()
		_, _, numberOfHeapPulls := shortestPathFunction(randomStart, randomDest, graph)
		diff := time.Since(start).Microseconds()
		timesAndPulls = append(timesAndPulls, []string{fmt.Sprint(diff), fmt.Sprint(numberOfHeapPulls)})
	}
	file, e := os.Create(outPath)
	if e != nil {
		fmt.Println(e)
	}
	writer := csv.NewWriter(file)

	e = writer.WriteAll(timesAndPulls)
	if e != nil {
		fmt.Println(e)
	}
	writer.Flush()
}
