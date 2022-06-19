package main

import (
	"OSM/src/backend/shortestPath"
	"fmt"
	"log"
)

func main() {
	outPath := "C:/Users/Matthias/go/src/OSM/out/"
	pathToFmiFile := "D:/OneDrive - stud.uni-stuttgart.de/Uni/10. Semester/FP-OSM/pbf files/oceanfmi.sec"
	var N int = 1e2

	//log.Println("Starting to bench dijkstra.")
	//outPathDijkstra := fmt.Sprintf(outPath + "dijkstra_bench.csv")
	//shortestPath.BenchDijkstra(outPathDijkstra, pathToFmiFile, N)
	//
	//log.Println("Starting to bench astar")
	//outPathAStar := fmt.Sprintf(outPath + "astar_bench.csv")
	//shortestPath.BenchAStar(outPathAStar, pathToFmiFile, N)

	log.Println("Starting to bench bidijkstra")
	outPathBiDijkstra := fmt.Sprintf(outPath + "bidijkstra_bench.csv")
	shortestPath.BenchBiDijkstra(outPathBiDijkstra, pathToFmiFile, N)

}
