package coastlines

import (
	"fmt"
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
				// use the first node in the way as the key for the map
				// since a node might be the first node for more than one way we need a special encoding to not overwrite an entry
				// we concatenate an index to the node id, the index describes how many ways with the same node as the starting node came before this way
				// so for the first occurance we use <NodeID>+0 for the second <NodeID>+1, ...
				//var index int64 = 0
				//
				//for true {
				//	_, exists := ways[v.NodeIDs[0]*10+index]
				//	if exists {
				//		index++
				//		happenings++
				//	}
				//	break
				//}
				//ways[v.NodeIDs[0]*10+index] = v.NodeIDs
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
	fmt.Printf("Time needed to evalute pbf file: %s\n", duration)
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

	for {
		mergeWays(ways)
		if oldLength == len(ways) || len(ways) == 1 {
			break
		}
		oldLength = len(ways)
	}

	wayNodes := make([][][]float64, len(ways))
	for i, way := range ways {
		for _, node := range way {
			coordinates := nodes[node]
			wayNodes[i] = append(wayNodes[i], coordinates)
		}
	}
	return wayNodes

}
