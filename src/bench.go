package main

import (
	"OSM/src/backend/datastructures"
	"OSM/src/backend/helpers"
	"OSM/src/backend/shortestPath"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	outPath := "C:/Users/Matthias/go/src/OSM/out/"
	pathToFmiFile := "D:/OneDrive - stud.uni-stuttgart.de/Uni/10. Semester/FP-OSM/pbf files/oceanfmi.sec"
	const N = 10
	graph := helpers.CreateGraphFromFile(pathToFmiFile)
	randomRoutes := helpers.CreateRandomRoutes(N, int32(len(graph.Nodes)))

	log.Println("Starting to bench dijkstra.")
	outPathDijkstra := fmt.Sprintf(outPath + "dijkstra_bench.csv")
	benchShortestPathWithFunction(shortestPath.Dijkstra, graph, randomRoutes, outPathDijkstra)

	log.Println("Starting to bench astar")
	outPathAStar := fmt.Sprintf(outPath + "astar_bench.csv")
	benchShortestPathWithFunction(shortestPath.AStar, graph, randomRoutes, outPathAStar)

	log.Println("Starting to bench bidijkstra")
	outPathBiDijkstra := fmt.Sprintf(outPath + "bidijkstra_bench.csv")
	benchShortestPathWithFunction(shortestPath.BiDijkstra, graph, randomRoutes, outPathBiDijkstra)

}

func benchShortestPathWithFunction(shortestPathFunction func(int32, int32, datastructures.Graph) (int32, []int32, int), graph datastructures.Graph, randomRoutes [][]int32, outPath string) {
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
