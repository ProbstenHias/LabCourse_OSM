package shortestPath

import (
	"OSM/src/backend/helpers"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const outPath = "C:/Users/Matthias/GolandProjects/LabCourse_OSM/out/dijkstra_bench.csv"
const PathToFmiFile = "C:/Users/Matthias/GolandProjects/LabCourse_OSM/in/oceanfmi.sec"
const N = 1000

func BenchDijkstra() {
	graph := helpers.CreateGraphFromFile(PathToFmiFile)
	var times []string
	for i := 0; i < N; i++ {
		randomStart := rand.Intn(len(graph.Nodes) - 1)
		randomDest := rand.Intn(len(graph.Nodes) - 1)
		start := time.Now()
		Dijkstra(int32(randomStart), int32(randomDest), graph)
		diff := time.Since(start).Microseconds()
		times = append(times, fmt.Sprint(diff))
	}
	f, e := os.Create(outPath)
	if e != nil {
		fmt.Println(e)
	}
	writer := csv.NewWriter(f)

	e = writer.Write(times)
	if e != nil {
		fmt.Println(e)
	}
	writer.Flush()

}
