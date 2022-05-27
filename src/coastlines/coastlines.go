package coastlines

import (
	"io"
	"log"
	"os"
	"runtime"
	"time"

	geojson "github.com/paulmach/go.geojson"
	"github.com/qedus/osmpbf"
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
	//var happenings int64 = 0
	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			// here we just save the node id as key and the lat and long coordinates as value in the map
			case *osmpbf.Node:
				nodes[v.ID] = []float64{v.Lat, v.Lon}

			//here we process the ways
			case *osmpbf.Way:
				// check if it is a coastline
				value, ok := v.Tags["natural"]
				if !ok || value != "coastline" {
					continue
				}

				ways[v.NodeIDs[0]] = v.NodeIDs

			case *osmpbf.Relation:
				continue
			default:
				log.Fatalf("unknown type %T\n", v)
			}
		}
	}
	end := time.Now()
	duration := end.Sub(start)
	log.Printf("Time needed to read pbf file: %s\n", duration)
	return nodes, ways
}

func CreateGeojson(nodes map[int64][]float64, ways map[int64][]int64) []byte {
	fc := geojson.NewFeatureCollection()
	for _, val := range ways {
		var lineNodes [][]float64
		for _, nodeId := range val {
			lineNodes = append(lineNodes, []float64{nodes[nodeId][1], nodes[nodeId][0]})
		}
		feature := geojson.NewLineStringFeature(lineNodes)
		feature.SetProperty("", 0)
		fc.AddFeature(feature)
	}
	rawJson, _ := fc.MarshalJSON()
	return rawJson
}

// merges ways where the end node of the way is the starting node of another way
func mergeWays(ways map[int64][]int64) {
	toDelete := make(map[int64]bool)
	for key, value := range ways {
		// if this way is already merged skip it
		if _, isInToDelete := toDelete[key]; isInToDelete {
			continue
		}
		lastNode := value[len(value)-1]
		// if the way is a loop we just skip it
		if lastNode == key {
			continue
		}
		nodes, exists := ways[lastNode]
		// if there is no other way that starts with this node, skip it
		if !exists {
			continue
		}
		toDelete[lastNode] = true
		newNodeSlice := append(value, nodes...)
		ways[key] = newNodeSlice

	}
	for key := range toDelete {
		delete(ways, key)
	}

}

func Main(path string) [][][]float64 {
	nodes, ways := readPBF(path)
	//fmt.Printf("Ways before merging: %d\n", len(ways))
	//mergeWays(ways)
	//fmt.Printf("Ways after merging: %d\n", len(ways))
	//mergeWays(ways)
	//fmt.Printf("Ways after merging twice: %d\n", len(ways))
	//mergeWays(ways)
	//fmt.Printf("Ways after merging thrice %d\n", len(ways))

	oldLength := len(ways)

	startTimeMerge := time.Now()
	for {
		mergeWays(ways)
		if oldLength == len(ways) || len(ways) == 1 {
			break
		}
		oldLength = len(ways)
	}
	log.Printf("Time to merge ways: %s\n", time.Since(startTimeMerge))

	wayNodes := make([][][]float64, len(ways))
	var curr = 0
	for _, way := range ways {
		for _, node := range way {
			coordinates := nodes[node]
			wayNodes[curr] = append(wayNodes[curr], coordinates)
		}
		curr++
	}
	log.Printf("Finished processing pbf file")
	return wayNodes

}
