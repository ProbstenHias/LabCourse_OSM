package piptest

import (
	"OSM/src/coastlines"
	"testing"
)

func TestTopLevel(t *testing.T) {
	var coordinates = [][]float64{
		{-5.997549029583585, -63.63281250000001},
		{-13.394988587960974, -38.25439453125001},
		{-19.50786449872789, -15.820312500000002},
		{-52.89929456807523, -59.1855812072754},
		//{-79.6674383780146, -56.25000000000001},
	}
	var correctClassification = []bool{false, true, true, false, true}

	wayNodes := coastlines.Main("D:/OneDrive - stud.uni-stuttgart.de/Uni/10. Semester/FP-OSM/pbf files/planet-coastlinespbf.sec")
	results := TopLevel(wayNodes, coordinates)

	for i, elem := range results {
		if elem == correctClassification[i] {
			continue
		}
		t.Errorf("Failed! A coodinate was not classified correctly\n"+
			"Coordinate: Lat %f Long%f\n"+
			"should have been %v but was %v", coordinates[i][0], coordinates[i][1], correctClassification[i], results[i])
		return
	}
	t.Log("PIP test was Successful")
}
