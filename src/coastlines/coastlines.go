package coastlines

import (
	"fmt"
	"github.com/qedus/osmpbf"
	"io"
	"log"
	"os"
	"runtime"
	"time"
)

const path string = "C:/Users/Matthias/Downloads/antarctica-latest.osm.pbf"

func ReadPBF() {
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

	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				// Process Node v
				continue

			case *osmpbf.Way:
				// Process Way v
				value, ok := v.Tags["natural"]
				if !ok || value != "coastline" {
					continue
				}
				fmt.Println(v.ID)

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

}
