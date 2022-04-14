package coastlines

import (
	"fmt"
	geojson "github.com/paulmach/go.geojson"
	"github.com/qedus/osmpbf"
	"io"
	"log"
	"os"
	"runtime"
	"time"
)

func readPBF(path string) (map[int64][]float64, map[int64][]int64) {
	start := time.Now()

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d := osmpbf.NewDecoder(f)

	d.SetBufferSize(osmpbf.MaxBlobSize)

	// start decoding with several goroutines, it is faster
	err = d.Start(runtime.GOMAXPROCS(runtime.NumCPU()))
	if err != nil {
		log.Fatal(err)
	}
	nodes := make(map[int64][]float64)
	ways := make(map[int64][]int64)
	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				nodes[v.ID] = []float64{v.Lon, v.Lat}

			case *osmpbf.Way:
				// Process Way v
				value, ok := v.Tags["natural"]
				if !ok || value != "coastline" {
					continue
				}
				ways[v.ID] = v.NodeIDs

			case *osmpbf.Relation:
				continue
			default:
				log.Fatalf("unknown type %T\n", v)
			}
		}
	}
	end := time.Now()
	duration := end.Sub(start)
	fmt.Printf("Time needed to evalute pbf file: %s\n", duration)
	return nodes, ways
}

func createGeojson(nodes map[int64][]float64, ways map[int64][]int64) []byte {
	fc := geojson.NewFeatureCollection()
	for _, val := range ways {
		var lineNodes [][]float64
		for _, nodeId := range val {
			lineNodes = append(lineNodes, nodes[nodeId])
		}
		feature := geojson.NewLineStringFeature(lineNodes)
		feature.SetProperty("", 0)
		fc.AddFeature(feature)
	}
	rawJson, _ := fc.MarshalJSON()
	return rawJson
}
func Main(path string) []byte {
	return createGeojson(readPBF(path))

}
