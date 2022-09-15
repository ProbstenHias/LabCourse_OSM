package main

import (
	"OSM/src/backend/helpers"
	"OSM/src/backend/shortestPath"
	"os"
	"path"
)

func main() {
	pathToFmi := os.Args[1]
	outPath := path.Dir(pathToFmi) + "/graphContracted.fmi"
	graph := helpers.CreateGraphFromFile(pathToFmi)
	chg := shortestPath.ContractGraph(graph)
	helpers.CreateFileFromContractedGraph(chg, outPath)
}
