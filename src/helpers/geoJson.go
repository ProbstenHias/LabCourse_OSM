package helpers

import geojson "github.com/paulmach/go.geojson"

func NodesToLineString(nodes [][]float64, distance int32) []byte {
	fc := geojson.NewFeatureCollection()
	feature := geojson.NewLineStringFeature(nodes)
	feature.SetProperty("distance", distance)
	fc.AddFeature(feature)
	rawJson, _ := fc.MarshalJSON()
	return rawJson
}
