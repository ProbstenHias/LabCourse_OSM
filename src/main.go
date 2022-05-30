package main

import (
	"OSM/src/backend/helpers"
	pre2 "OSM/src/backend/pre"
	"os"
)

const N = 1e6

func main() {
	path := os.Getenv("PBF")

	wayNodes := pre2.GenerateCoastlines(path)
	points := pre2.GenerateSpherePoints(N)
	classification := pre2.TopLevel(wayNodes, points)
	graph := pre2.GenerateGraphFromPoints(N, points, classification)
	helpers.CreateFileFromGraph(graph, "./out/graph.fmi")

}
