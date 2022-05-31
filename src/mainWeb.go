package main

import (
	"OSM/src/web"
	"os"
)

func main() {
	pathToFmi := os.Args[1]
	port := os.Args[2]
	web.Main(pathToFmi, port)
}
