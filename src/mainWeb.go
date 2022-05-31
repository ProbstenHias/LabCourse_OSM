package main

import (
	"OSM/src/web"
	"os"
)

func main() {
	pathToFmi := os.Args[1]
	web.Main(pathToFmi)
}
