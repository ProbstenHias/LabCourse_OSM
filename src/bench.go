package main

import "OSM/src/backend/shortestPath"

func main() {
	outPath := "C:/Users/Matthias/go/src/OSM/out/dijkstra_bench.csv"
	pathToFmiFile := "D:/OneDrive - stud.uni-stuttgart.de/Uni/10. Semester/FP-OSM/pbf files/oceanfmi.sec"
	shortestPath.BenchDijkstra(outPath, pathToFmiFile, 1e3)
}
