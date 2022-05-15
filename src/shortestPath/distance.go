package shortestPath

import "math"

const earthDiameter = 12_742_000

func haversine(p1 []float64, p2 []float64) uint32 {
	const p = math.Pi / 180
	var a = 0.5 - math.Cos((p2[0]-p1[0])*p)/2 + math.Cos(p1[0]*p)*math.Cos(p2[0])*(1-math.Cos((p2[1]-p1[1])*p))/2
	var d = earthDiameter * math.Asin(math.Sqrt(a))
	return uint32(d)
}
