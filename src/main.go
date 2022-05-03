package main

import (
	"OSM/src/coastlines"
	"log"
	"os"
)

func main() {
	path := os.Getenv("OSM_ANTARCTICAPBF")
	json := coastlines.Main(path)
	if err := os.WriteFile("../out/out.json", json, 06666); err != nil {

		log.Fatal(err)
	}
}
