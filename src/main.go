package main

import (
	"OSM/src/coastlines"
	"OSM/src/piptest"
	"OSM/src/spherePoints"
	"fmt"
	"log"
	"os"
)

func main() {
	path := "E:/Classes Infotech/4th Semster/Fachpraktika/Code/data/planet-coastlines.pbf" //change to own path
	nodes, ways := coastlines.Main(path)

	var no_of_nodes int64 = 10000

	tested_p_array := piptest.Top_level(nodes, ways, no_of_nodes)
	fmt.Printf("Size of mapped points: %d", len(tested_p_array))

	json_file := spherePoints.PointsToGeoJson(tested_p_array)
	if err := os.WriteFile("../out/out_pip_test.json", json_file, 06666); err != nil {

		log.Fatal(err)
	}
}
