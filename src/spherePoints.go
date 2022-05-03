package main

import (
	"fmt"
	geojson "github.com/paulmach/go.geojson"
	"math"
	"time"
)

func generatePointsOnSphere(n int64) [][]float64 {
	start := time.Now()
	var points [][]float64
	var nCount = 0
	var a = 4 * math.Pi / float64(n)
	var d = math.Sqrt(a)
	var mTheta = math.Round(math.Pi / d)
	var dTheta = math.Pi / mTheta
	var dPhi = a / dTheta
	for m := 0.; m < mTheta; m++ {
		var theta = math.Pi * (m + 0.5) / mTheta
		var mPhi = math.Round(2 * math.Pi * math.Sin(theta) / dPhi)
		for n := 0.; n < mPhi; n++ {
			var phi = 2 * math.Pi * n / mPhi
			long, lat := anglesToLatLong(theta, phi)
			longlat := []float64{long, lat}
			points = append(points, longlat)
			nCount++
		}
	}
	end := time.Now()
	diff := end.Sub(start)
	fmt.Printf("time to create points: %s", diff)
	return points
}

func anglesToLatLong(theta, phi float64) (float64, float64) {
	// min max scaling
	// phi from [0,2pi] to -180 to 180
	var long = -180 + ((phi * 360) / (2 * math.Pi))
	// theta from [0, pi] to -45 to 45
	var lat = -90 + ((theta * 180) / math.Pi)
	return long, lat
}

func pointsToGeoJson(points [][]float64) []byte {
	fc := geojson.NewFeatureCollection()
	for _, elem := range points {
		feature := geojson.NewPointFeature(elem)
		feature.SetProperty("", 0)
		fc.AddFeature(feature)
	}
	rawJson, _ := fc.MarshalJSON()
	return rawJson
}
