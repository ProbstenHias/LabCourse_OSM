package spherePoints

import (
	"OSM/src/datastructures"
	"OSM/src/helpers"
	"log"
	"math"
	"sort"
	"time"
)

const LowerBound = -85

func GenerateSpherePoints(n int) [][]float64 {
	startTime := time.Now()
	var points [][]float64
	var a = 4 * math.Pi / float64(n)
	var d = math.Sqrt(a)
	var mTheta = math.Round(math.Pi / d)
	var dTheta = math.Pi / mTheta
	var dPhi = a / dTheta
	for m := 0.; m < mTheta; m++ {
		var theta = math.Pi * (m + 0.5) / mTheta
		var lat = thetaToLat(theta)
		if lat < LowerBound {
			continue
		}
		var mPhi = math.Round(2 * math.Pi * math.Sin(theta) / dPhi)
		for n := 0.; n < mPhi; n++ {
			var phi = 2 * math.Pi * n / mPhi
			lat, long := anglesToLatLong(theta, phi)
			point := []float64{lat, long}
			points = append(points, point)
		}
	}
	log.Printf("Time to generate equidistant points on spere: %s\n", time.Since(startTime))
	return points
}

func GenerateGraphFromPoints(n int64, points [][]float64, classification []bool) datastructures.Graph {
	start := time.Now()

	// here we store the pointsInWater that were generated and are in fact in water
	var pointsInWater [][]float64

	// here we store the edges between the pointsInWater
	var edges []Edge

	// this is a helper variable that saves the index of a special point
	// this point is the first one that was created in a new latitude row
	// in other words, the point where n = 0
	var firstPointOnCurrLat = 0

	// this is a helper variable that saves the index of a special point
	// this point is the "firstPointOnCurrLat" of the previous iteration
	var firstPointOnPrevLat = -1

	// this variable saves the maximum longitude difference in the current iteration
	var maxLongDiff = 0.

	var currPointIdx = 0

	// these things are all for the calculation of the coordinates
	var a = 4 * math.Pi / float64(n)
	var d = math.Sqrt(a)
	var mTheta = math.Round(math.Pi / d)
	var dTheta = math.Pi / mTheta
	var dPhi = a / dTheta

	for m := 0.; m < mTheta; m++ {
		// this is the lat coordinate
		var theta = math.Pi * (m + 0.5) / mTheta

		// skip all points below certain threshold
		var lat = thetaToLat(theta)
		if lat < LowerBound {
			continue
		}
		var mPhi = math.Round(2 * math.Pi * math.Sin(theta) / dPhi)
		maxLongDiff = 2 * math.Pi / mPhi
		for n := 0.; n < mPhi; n++ {
			var phi = 2 * math.Pi * n / mPhi
			point := points[currPointIdx]
			if !classification[currPointIdx] {
				currPointIdx++
				continue
			}
			pointsInWater = append(pointsInWater, point)

			var currPointInWaterIdx = len(pointsInWater) - 1

			// in this case we started with a new lat
			if currPointInWaterIdx > 0 && point[0] != pointsInWater[currPointInWaterIdx-1][0] {
				firstPointOnPrevLat = firstPointOnCurrLat
				firstPointOnCurrLat = currPointInWaterIdx
			}

			// create edge to the left
			// this will be done in case it is not the first point on this lat
			// and if the last point was in water
			if n != 0 && classification[currPointIdx-1] {
				var leftPointIdx = currPointInWaterIdx - 1
				var distance = helpers.Haversine(pointsInWater[leftPointIdx], pointsInWater[currPointInWaterIdx])

				edge1, edge2 := createForwardAndBackwardEdge(currPointInWaterIdx, leftPointIdx, distance)
				edges = append(edges, edge1, edge2)
			}

			// create wrap around edge to the left if this is the last point on this latitude line
			if n+1 >= mPhi {
				// check if the wrap around point was removed
				if pointsInWater[firstPointOnCurrLat][1] == -180 {
					var distance = helpers.Haversine(pointsInWater[firstPointOnCurrLat], pointsInWater[currPointInWaterIdx])

					edge1, edge2 := createForwardAndBackwardEdge(firstPointOnCurrLat, currPointInWaterIdx, distance)

					edges = append(edges, edge1, edge2)
				}
			}

			// create edge down
			// check if there are pointsInWater below
			if firstPointOnPrevLat != -1 {
				// search in the range from pointsInWater[firstPointOnPrevLat] to pointsInWater[firstPoints - 1] for and point that is
				// close enough to be connected to the current point
				// if there are multiple pointsInWater that are close enough choose the closest of them

				// longitude the furthest to the left
				var toLeft = 0.
				// we have to make a differentiation for the case that the current point is the first on this lat
				// in this case there is always a node directly below it, so we can just use -180 as toLeft
				if n == 0 {
					toLeft = point[1]
				} else {
					toLeft = phiToLong(phi - (maxLongDiff / 2))
				}

				// the long the farthest to the right
				// we never run into the case that toRight gets over 180 because of how the point calculation works
				toRight := phiToLong(phi + maxLongDiff/2)

				// longitude the furthest to the right
				var index = -1
				for i := firstPointOnPrevLat; i < firstPointOnCurrLat; i++ {
					var long = pointsInWater[i][1]

					if long < toLeft {
						continue
					}
					if long >= toRight {
						continue
					}

					index = i
					break

				}
				if index == -1 {
					currPointIdx++
					continue
				}
				var distance = helpers.Haversine(pointsInWater[index], pointsInWater[currPointInWaterIdx])

				// it could be the case that the next point is even closer
				// but we make sure that this point is not on the curr lat line
				var distanceNext = math.MaxInt
				if index+1 < firstPointOnCurrLat {
					distanceNext = helpers.Haversine(pointsInWater[index+1], pointsInWater[currPointInWaterIdx])
				}
				if distanceNext < distance {
					index++
					distance = distanceNext
				}

				edge1, edge2 := createForwardAndBackwardEdge(index, currPointInWaterIdx, distance)

				edges = append(edges, edge1, edge2)

			}
			currPointIdx++

		}
	}

	// now that we have all edges and nodes we can create a graph struct
	graph := createGraphFromPointsAndEdges(pointsInWater, edges)
	log.Printf("Time to generate graph from pointsInWater: %s \n", time.Since(start))
	return graph
}

func createGraphFromPointsAndEdges(pointsInWater [][]float64, edges []Edge) datastructures.Graph {
	graph := datastructures.Graph{
		Nodes:    pointsInWater,
		Edges:    make([]int32, len(edges)),
		Distance: make([]int32, len(edges)),
		Offset:   make([]int32, len(pointsInWater)),
	}
	graph.Nodes = pointsInWater
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].from < edges[j].from
	})

	graph.Offset[0] = 0
	var from = 0
	for i, edge := range edges {
		graph.Edges[i] = int32(edge.to)
		graph.Distance[i] = int32(edge.distance)
		// for all nodes with no edges add offset as well
		for j := from + 1; j <= edge.from; j++ {
			graph.Offset[j] = int32(i)
		}

		from = edge.from
	}
	return graph
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
