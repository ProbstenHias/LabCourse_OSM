package main

import (
	"OSM/src/backend/helpers"
	"OSM/src/backend/pre"
	"log"
	"os"
	"path"
)

const N = 1e6

func main() {
	pathToPBF := os.Args[1]
	pathToFmi := path.Dir(pathToPBF)
	pathToFmi = pathToFmi + "/graph.fmi"
	log.Println(pathToFmi)
	wayNodes := pre.GenerateCoastlines(pathToPBF)
	points := pre.GenerateSpherePoints(N)
	classification := pre.TopLevel(wayNodes, points)
	graph := pre.GenerateGraphFromPoints(N, points, classification)
	helpers.CreateFileFromGraph(graph, pathToFmi)
}
