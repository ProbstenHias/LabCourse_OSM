package shortestPath

import (
	"OSM/src/backend/helpers"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func BenchBiDijkstra(outPath, pathToFmiFile string, N int) {
	graph := helpers.CreateGraphFromFile(pathToFmiFile)
	var timesAndPulls [][]string
	for i := 0; i < N; i++ {
		randomStart := rand.Intn(len(graph.Nodes) - 1)
		randomDest := rand.Intn(len(graph.Nodes) - 1)
		start := time.Now()
		_, _, numberOfHeapPulls := BiDijkstraWithNumberOfHeapPulls(int32(randomStart), int32(randomDest), graph)
		diff := time.Since(start).Microseconds()
		timesAndPulls = append(timesAndPulls, []string{fmt.Sprint(diff), fmt.Sprint(numberOfHeapPulls)})
	}
	f, e := os.Create(outPath)
	if e != nil {
		fmt.Println(e)
	}
	writer := csv.NewWriter(f)

	e = writer.WriteAll(timesAndPulls)
	if e != nil {
		fmt.Println(e)
	}
	writer.Flush()

}
