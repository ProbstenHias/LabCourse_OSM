package shortestPath

import "math"

const R = 6_371_000

func haversine(p1 []float64, p2 []float64) int {
	var dLat = deg2grad(p2[0] - p1[0])
	var dLon = deg2grad(p2[1] - p1[1])
	var a = math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(deg2grad(p1[0]))*math.Cos(deg2grad(p2[0]))*math.Sin(dLon/2)*math.Sin(dLon/2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	var d = R * c
	return int(d)
}

func deg2grad(deg float64) float64 {
	return deg * (math.Pi / 180)
}
