package helpers

import geojson "github.com/paulmach/go.geojson"

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
