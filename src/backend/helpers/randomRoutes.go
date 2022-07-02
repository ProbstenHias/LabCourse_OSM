package helpers

import "math/rand"

func CreateRandomRoutes(numberOfRoutes int, numberOfNodes int) [][]int {
	routes := make([][]int, numberOfRoutes)
	for i := 0; i < numberOfRoutes; i++ {
		randomStart := rand.Intn(numberOfNodes)
		randomDest := rand.Intn(numberOfNodes)
		routes[i] = []int{randomStart, randomDest}
	}
	return routes
}
