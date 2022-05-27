package main

import (
	"OSM/src/coastlines"
	"OSM/src/piptest"
	"OSM/src/spherePoints"
)

func main() {
	//path := "E:/Classes Infotech/4th Semster/Fachpraktika/Code/data/planet-coastlines.pbf" //change to own path
	//nodes, ways := coastlines.GenerateCoastlines(path)
	//
	//var no_of_nodes int64 = 10000
	//
	//tested_p_array := piptest.TopLevel(nodes, ways, no_of_nodes)
	//fmt.Printf("Size of mapped points: %d", len(tested_p_array))
	//
	//json_file := spherePoints.PointsToGeoJson(tested_p_array)
	//if err := os.WriteFile("../out/out_pip_test.json", json_file, 06666); err != nil {
	//
	//	log.Fatal(err)
	//}

	//web.GenerateCoastlines("D:/OneDrive - stud.uni-stuttgart.de/Uni/10. Semester/FP-OSM/pbf files/oceanfmi.sec")

	wayNodes := coastlines.GenerateCoastlines("D:/OneDrive - stud.uni-stuttgart.de/Uni/10. Semester/FP-OSM/pbf files/planet-coastlinespbf.sec")
	points := spherePoints.GenerateSpherePoints(1e6)
	piptest.TopLevel(wayNodes, points)

}
