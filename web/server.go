package web

import (
	"OSM/src/datastructures"
	"OSM/src/helpers"
	"OSM/src/shortestPath"
	"log"
	"math"
	"net/http"
	"strconv"
)

func Main(pathToFmiFile string) {
	graph := helpers.CreateGraphFromFile(pathToFmiFile)

	fileServer := http.FileServer(http.Dir("web/static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/route", routeHandler(graph))
	http.HandleFunc("/point", pointHandler(graph))

	log.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
func pointHandler(graph datastructures.Graph) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/route" {
			http.Error(w, "404 not found.", http.StatusNotFound)
			return
		}

		if r.Method != "GET" {
			http.Error(w, "Method is not supported.", http.StatusNotFound)
			return
		}

		query := r.URL.Query()
		lat, _ := strconv.ParseFloat(query["lat"][0], 64)
		lng, _ := strconv.ParseFloat(query["lng"][0], 64)
		idx, node := helpers.GetClosestNodeInGraph([]float64{lat, lng}, graph)
		rawJson := helpers.NodeToPoint(node, idx)
		w.WriteHeader(http.StatusOK)
		w.Write(rawJson)
	}

}
func routeHandler(graph datastructures.Graph) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/route" {
			http.Error(w, "404 not found.", http.StatusNotFound)
			return
		}

		if r.Method != "GET" {
			http.Error(w, "Method is not supported.", http.StatusNotFound)
			return
		}
		query := r.URL.Query()
		startLat, _ := strconv.ParseFloat(query["startLat"][0], 64)
		startLng, _ := strconv.ParseFloat(query["startLng"][0], 64)
		endLat, _ := strconv.ParseFloat(query["endLat"][0], 64)
		endLng, _ := strconv.ParseFloat(query["endLng"][0], 64)

		start := []float64{startLat, startLng}
		startIdx, _ := helpers.GetClosestNodeInGraph(start, graph)
		end := []float64{endLat, endLng}
		endIdx, _ := helpers.GetClosestNodeInGraph(end, graph)

		distance, prev := shortestPath.Dijkstra(startIdx, endIdx, graph)
		if distance == math.MaxInt32 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		lineNodes := helpers.CreatePathFromPrev(startIdx, endIdx, prev, graph)
		rawJson := helpers.NodesToLineString(lineNodes, distance)

		w.WriteHeader(http.StatusOK)
		w.Write(rawJson)
	}
}
