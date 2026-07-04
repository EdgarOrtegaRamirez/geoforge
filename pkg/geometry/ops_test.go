package geometry

import (
	"testing"

	"github.com/EdgarOrtegaRamirez/geoforge/pkg/geojson"
)

func TestCalculateBBoxPoint(t *testing.T) {
	g := &geojson.Geometry{
		Type:        geojson.PointType,
		Coordinates: []interface{}{100.0, 50.0},
	}
	bb, ok := CalculateBBox(g)
	if !ok {
		t.Fatal("Expected ok=true")
	}
	if bb.MinLon != 100.0 || bb.MaxLon != 100.0 || bb.MinLat != 50.0 || bb.MaxLat != 50.0 {
		t.Errorf("Unexpected bbox: %+v", bb)
	}
}

func TestCalculateBBoxFeatureCollection(t *testing.T) {
	fc := &geojson.FeatureCollection{
		Type: "FeatureCollection",
		Features: []*geojson.Feature{
			{
				Type: "Feature",
				Geometry: &geojson.Geometry{
					Type:        geojson.PointType,
					Coordinates: []interface{}{10.0, 20.0},
				},
			},
			{
				Type: "Feature",
				Geometry: &geojson.Geometry{
					Type:        geojson.PointType,
					Coordinates: []interface{}{30.0, 40.0},
				},
			},
		},
	}
	bb, ok := CalculateBBox(fc)
	if !ok {
		t.Fatal("Expected ok=true")
	}
	if bb.MinLon != 10.0 || bb.MaxLon != 30.0 || bb.MinLat != 20.0 || bb.MaxLat != 40.0 {
		t.Errorf("Unexpected bbox: %+v", bb)
	}
}

func TestCalculateBBoxEmpty(t *testing.T) {
	fc := &geojson.FeatureCollection{
		Type:     "FeatureCollection",
		Features: []*geojson.Feature{},
	}
	_, ok := CalculateBBox(fc)
	if ok {
		t.Error("Expected ok=false for empty collection")
	}
}

func TestCalculateCentroid(t *testing.T) {
	g := &geojson.Geometry{
		Type: geojson.LineStringType,
		Coordinates: []interface{}{
			[]interface{}{0.0, 0.0},
			[]interface{}{10.0, 10.0},
		},
	}
	c, ok := CalculateCentroid(g)
	if !ok {
		t.Fatal("Expected ok=true")
	}
	if c.Lon != 5.0 || c.Lat != 5.0 {
		t.Errorf("Expected centroid (5, 5), got (%g, %g)", c.Lon, c.Lat)
	}
}

func TestHaversineDistance(t *testing.T) {
	// New York to London approx 5570km
	d := HaversineDistance(-74.006, 40.7128, -0.1276, 51.5074)
	if d < 5500000 || d > 5600000 {
		t.Errorf("Expected ~5570km, got %g", d/1000)
	}
}

func TestDouglasPeucker(t *testing.T) {
	// Points with a detour that should be simplified
	points := []geojson.Position{
		{0, 0},
		{10, 0},
		{10, 10},
		{20, 0},
		{20, 0},
	}
	simplified := DouglasPeucker(points, 5.0)
	if len(simplified) >= len(points) {
		t.Errorf("Expected fewer points, got %d (original %d)", len(simplified), len(points))
	}
}
