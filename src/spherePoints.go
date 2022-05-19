package main

import (
	"OSM/src/helpers"
	"fmt"
	geojson "github.com/paulmach/go.geojson"
	"math"
	"time"
)

const LowerBound = -85

func generatePointsOnSphere(n int64) ([][]float64, []Edge) {
	start := time.Now()

	// here we store the points that were generated and are in fact in water
	var points [][]float64

	// here we store the edges between the points
	var edges []Edge

	// this is a helper variable that saves the index of a special point
	// this point is the first one that was created in a new latitude row
	// in other words, the point where n = 0
	var firstPoint = -1

	// this is a helper variable that saves the index of a special point
	// this point is the "firstPoint" of the previous iteration
	var firstPointBefore = -1

	// this variable saves the maximum longitude difference in the current iteration
	var maxLongDiff = 0.

	// these things are all for the calculation of the coordinates
	var a = 4 * math.Pi / float64(n)
	var d = math.Sqrt(a)
	var mTheta = math.Round(math.Pi / d)
	var dTheta = math.Pi / mTheta
	var dPhi = a / dTheta
	for m := 0.; m < mTheta; m++ {
		var isLastPointInWater = false
		var theta = math.Pi * (m + 0.5) / mTheta
		var lat = thetaToLat(theta)
		if lat < LowerBound {
			continue
		}
		var mPhi = math.Round(2 * math.Pi * math.Sin(theta) / dPhi)
		maxLongDiff = 2 * math.Pi / mPhi
		for n := 0.; n < mPhi; n++ {
			var phi = 2 * math.Pi * n / mPhi
			lat, long := anglesToLatLong(theta, phi)
			point := []float64{lat, long}
			if isInPoly(point) {
				isLastPointInWater = false
				continue
			}
			points = append(points, point)

			var currPointIdx = len(points) - 1

			if phi == 0 {
				firstPointBefore = firstPoint
				firstPoint = currPointIdx
			}

			// create edge to the left
			if n != 0 && isLastPointInWater {
				var leftPointIdx = len(points) - 2
				var distance = helpers.Haversine(points[leftPointIdx], points[currPointIdx])

				edge1, edge2 := createForwardAndBackwardEdge(currPointIdx, leftPointIdx, distance)
				edges = append(edges, edge1, edge2)
			}

			// create wrap around edge to the left if this is the last point on this latitude line
			if n+1 >= mPhi {
				// check if the wrap around point was removed
				if points[firstPoint][1] == -180 {
					var distance = helpers.Haversine(points[firstPoint], points[currPointIdx])

					edge1, edge2 := createForwardAndBackwardEdge(firstPoint, currPointIdx, distance)

					edges = append(edges, edge1, edge2)
				}
			}

			// create edge down
			// check if there are points below
			if firstPointBefore != -1 {
				// search in the range from points[firstPointBefore] to points[firstPoints - 1] for and point that is
				// close enough to be connected to the current point
				// if there are multiple points that are close enough choose the closest of them

				// longitude the furthest to the left
				var toLeft = 0.
				if n == 0 {
					toLeft = phiToLong(2*math.Pi - maxLongDiff/2)
				} else {
					toLeft = phiToLong(phi - (maxLongDiff / 2))
				}
				var toRight = 0.
				if n+1 >= mPhi {
					toRight = phiToLong((phi + maxLongDiff/2) - 2*math.Pi)
				} else {
					toRight = phiToLong(phi + maxLongDiff/2)
				}

				// longitude the furthest to the right
				var index = -1
				for i := firstPointBefore; i < firstPoint; i++ {
					var long = points[i][1]
					// wrap around case
					if toLeft > toRight {
						if long < toLeft && long >= toRight {
							continue
						}
					} else {
						// non-wrap around case
						if long < toLeft {
							continue
						}
						if long >= toRight {
							continue
						}
					}

					index = i
					break

				}
				// heyo
				if index == -1 {
					continue
				}
				var distance = helpers.Haversine(points[index], points[currPointIdx])

				// it could be the case that the next point is even closer
				// but we make sure that this point is not the current point itself
				var distanceNext = math.MaxInt
				if index+1 != currPointIdx {
					distanceNext = helpers.Haversine(points[index+1], points[currPointIdx])
				}
				if distanceNext < distance {
					index += 1
					distance = distanceNext
				}

				edge1, edge2 := createForwardAndBackwardEdge(index, currPointIdx, distance)

				edges = append(edges, edge1, edge2)

			}

		}
	}
	end := time.Now()
	diff := end.Sub(start)
	fmt.Printf("time to create points: %s", diff)
	return points, edges
}

func anglesToLatLong(theta, phi float64) (float64, float64) {
	// min max scaling
	// phi from [0,2pi] to -180 to 180
	var long = -180 + ((phi * 360) / (2 * math.Pi))
	// theta from [0, pi] to -90 to 90
	var lat = -90 + ((theta * 180) / math.Pi)
	return lat, long
}

func thetaToLat(theta float64) float64 {
	return -90 + ((theta * 180) / math.Pi)

}

func phiToLong(phi float64) float64 {
	return -180 + ((phi * 360) / (2 * math.Pi))
}

func pointsToGeoJson(points [][]float64) []byte {
	fc := geojson.NewFeatureCollection()
	for _, elem := range points {
		feature := geojson.NewPointFeature([]float64{elem[1], elem[0]})
		feature.SetProperty("", 0)
		fc.AddFeature(feature)
	}
	rawJson, _ := fc.MarshalJSON()
	return rawJson
}

func isInPoly(point []float64) bool {
	return false
}

func createEdge(from, to, distance int) Edge {
	return Edge{
		from:     from,
		to:       to,
		distance: distance,
	}
}

func createForwardAndBackwardEdge(node1, node2, distance int) (Edge, Edge) {
	edge1 := createEdge(node1, node2, distance)
	edge2 := createEdge(node2, node1, distance)
	return edge1, edge2
}

type Edge struct {
	from     int
	to       int
	distance int
}
