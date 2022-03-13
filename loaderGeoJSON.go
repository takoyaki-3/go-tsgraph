package gotsgraph

import (
	. "github.com/takoyaki-3/go-geojson"
	gojson "github.com/takoyaki-3/go-json"
	gm "github.com/takoyaki-3/go-map/v2"
	"time"
)

func (ts *TSGraph) DumpToGeoJSON(g *gm.Graph, date time.Time, graphFileName, tripsFileName string) error {
	fc1 := FeatureCollection{
		Type: "FeatureCollection",
	}
	fc2 := FeatureCollection{
		Type: "FeatureCollection",
	}

	for _, e := range ts.Edges {
		line1Coordinates := [][]float64{}
		line2Coordinates := [][]float64{}
		line1 := Geometry{}
		line2 := Geometry{}
		fromP := ts.Points[e.FromPoint]

		from := g.GetStop(ts.Places[fromP.PlaceIndex])
		toP := ts.Points[e.ToPoint]
		to := g.GetStop(ts.Places[toP.PlaceIndex])
		if fromP.Time == toP.Time {
			continue
		}
		t := float64(date.Unix())
		line1Coordinates = append(line1Coordinates, []float64{from.Longitude, from.Latitude, float64(fromP.Time)})
		line2Coordinates = append(line2Coordinates, []float64{from.Longitude, from.Latitude, 0, float64(fromP.Time) + t})
		line2Coordinates = append(line2Coordinates, []float64{from.Longitude, from.Latitude, 0, float64(fromP.Time) + t})
		line1Coordinates = append(line1Coordinates, []float64{to.Longitude, to.Latitude, float64(toP.Time)})
		line2Coordinates = append(line2Coordinates, []float64{to.Longitude, to.Latitude, 0, float64(toP.Time) + t})
		line2Coordinates = append(line2Coordinates, []float64{to.Longitude, to.Latitude, 0, float64(toP.Time) + t})
		line1.Type = "LineString"
		line2.Type = "LineString"
		line1.Coordinates = line1Coordinates
		line2.Coordinates = line2Coordinates

		fc1.Features = append(fc1.Features, Feature{
			Type:       "Feature",
			Geometry:   line1,
			Properties: map[string]string{},
		})
		fc2.Features = append(fc2.Features, Feature{
			Type:       "Feature",
			Geometry:   line2,
			Properties: map[string]string{},
		})
	}

	// Output as a file
	if err := gojson.DumpToFile(fc1, "3dgraph.geojson"); err != nil {
		return err
	}
	if err := gojson.DumpToFile(fc2, "trips.geojson"); err != nil {
		return err
	}
	return nil
}
