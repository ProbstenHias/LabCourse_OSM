package web

import (
	"OSM/src/backend/datastructures"
	helpers2 "OSM/src/backend/helpers"
	"OSM/src/backend/shortestPath"
	"log"
	"math"
	"net/http"
	"strconv"
)

func Main(pathToFmiFile string) {
	graph := helpers2.CreateGraphFromFile(pathToFmiFile)

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
		if r.URL.Path != "/point" {
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
		idx, node := helpers2.GetClosestNodeInGraph([]float64{lat, lng}, graph)
		rawJson := helpers2.NodeToPointGeoJson(node, idx)
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

		startIdx, _ := strconv.Atoi(query["startIdx"][0])
		startIdx32 := int32(startIdx)
		endIdx, _ := strconv.Atoi(query["endIdx"][0])
		endIdx32 := int32(endIdx)

		distance, prev := shortestPath.Dijkstra(startIdx32, endIdx32, graph)
		if distance == math.MaxInt32 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		lineNodes := helpers2.CreatePathFromPrev(startIdx32, endIdx32, prev, graph)
		rawJson := helpers2.NodesToLineStringGeoJson(lineNodes, distance)

		w.WriteHeader(http.StatusOK)
		w.Write(rawJson)
	}
}
