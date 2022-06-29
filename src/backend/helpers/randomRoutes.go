package helpers

import "math/rand"

func CreateRandomRoutes(numberOfRoutes int, numberOfNodes int32) [][]int32 {
	routes := make([][]int32, numberOfRoutes)
	for i := 0; i < numberOfRoutes; i++ {
		randomStart := rand.Int31n(numberOfNodes)
		randomDest := rand.Int31n(numberOfNodes)
		routes[i] = []int32{randomStart, randomDest}
	}
	return routes
}
