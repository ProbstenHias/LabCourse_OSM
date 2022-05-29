package helpers

import (
	"OSM/src/datastructures"
	geojson "github.com/paulmach/go.geojson"
	"log"
	"os"
)

func NodesToLineString(nodes [][]float64, distance int32) []byte {
	lineNodes := make([][]float64, len(nodes))
	for i, elem := range nodes {
		lineNodes[i] = []float64{elem[1], elem[0]}
	}
	fc := geojson.NewFeatureCollection()
	feature := geojson.NewLineStringFeature(lineNodes)
	feature.SetProperty("distance", distance)
	fc.AddFeature(feature)
	rawJson, _ := fc.MarshalJSON()
	return rawJson
}

func NodeToPoint(node []float64, idx int32) []byte {
	fc := geojson.NewFeatureCollection()
	feature := geojson.NewPointFeature([]float64{node[1], node[0]})
	feature.SetProperty("idx", idx)
	fc.AddFeature(feature)
	rawJson, _ := fc.MarshalJSON()
	return rawJson
}

func NodesToPoints(nodes [][]float64) []byte {
	fc := geojson.NewFeatureCollection()
	for _, node := range nodes {
		feature := geojson.NewPointFeature([]float64{node[1], node[0]})
		feature.SetProperty("", nil)
		fc.AddFeature(feature)
	}
	rawJson, _ := fc.MarshalJSON()
	return rawJson
}

func GraphToGeoJson(graph datastructures.Graph) []byte {
	fc := geojson.NewFeatureCollection()
	for i, node := range graph.Nodes {
		pointFeature := geojson.NewPointFeature([]float64{node[1], node[0]})
		pointFeature.SetProperty("idx", i)
		fc.AddFeature(pointFeature)
		edges := graph.GetAllOutgoingEdgesOfNode(int32(i))
		for _, edgeIdx := range edges {
			toNode := graph.Nodes[graph.Edges[edgeIdx]]
			lineStringFeature := geojson.NewLineStringFeature([][]float64{
				{node[1], node[0]},
				{toNode[1], toNode[0]},
			})
			lineStringFeature.SetProperty("distance", graph.Distance[edgeIdx])
			fc.AddFeature(lineStringFeature)
		}
	}
	log.Printf("Created %d features\n", len(fc.Features))
	rawJson, _ := fc.MarshalJSON()
	return rawJson
}

func GeoJsonToFile(json []byte, pathToFile string) {
	if err := os.WriteFile(pathToFile, json, 06666); err != nil {

		log.Fatal(err)
	}
}
