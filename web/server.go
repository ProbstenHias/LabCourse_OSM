package main

import (
	"fmt"
	geojson "github.com/paulmach/go.geojson"
	"log"
	"net/http"
	"strconv"
)

func main() {
	fileServer := http.FileServer(http.Dir("web/static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/route", routeHandler)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
func routeHandler(w http.ResponseWriter, r *http.Request) {
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

	start := []float64{startLng, startLat}
	end := []float64{endLng, endLat}
	lineNodes := [][]float64{start, end}
	fc := geojson.NewFeatureCollection()
	feature := geojson.NewLineStringFeature(lineNodes)
	fc.AddFeature(feature)
	rawJson, _ := fc.MarshalJSON()

	w.WriteHeader(200)
	w.Write(rawJson)
}
