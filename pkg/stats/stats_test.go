package stats

import (
	"testing"

	"github.com/EdgarOrtegaRamirez/geoforge/pkg/geojson"
)

func TestAnalyzeFeatureCollection(t *testing.T) {
	fc := &geojson.FeatureCollection{
		Type: "FeatureCollection",
		Features: []*geojson.Feature{
			{
				Type: "Feature",
				Geometry: &geojson.Geometry{
					Type:        geojson.PointType,
					Coordinates: []interface{}{100.0, 50.0},
				},
				Properties: map[string]interface{}{
					"name": "Test",
					"pop":  12345.0,
				},
			},
		},
	}

	s, err := Analyze(fc)
	if err != nil {
		t.Fatalf("Analyze error: %v", err)
	}

	if s.FeatureCount != 1 {
		t.Errorf("Expected 1 feature, got %d", s.FeatureCount)
	}

	if s.GeometryTypes["Point"] != 1 {
		t.Errorf("Expected 1 Point geometry, got %d", s.GeometryTypes["Point"])
	}

	if len(s.PropertyKeys) != 2 {
		t.Errorf("Expected 2 property keys, got %d", len(s.PropertyKeys))
	}

	if s.BBox == nil {
		t.Error("Expected bounding box to be calculated")
	}
}

func TestAnalyzeEmpty(t *testing.T) {
	fc := &geojson.FeatureCollection{
		Type:     "FeatureCollection",
		Features: []*geojson.Feature{},
	}

	s, err := Analyze(fc)
	if err != nil {
		t.Fatalf("Analyze error: %v", err)
	}

	if s.FeatureCount != 0 {
		t.Errorf("Expected 0 features, got %d", s.FeatureCount)
	}
}

func TestAnalyzeGeometry(t *testing.T) {
	g := &geojson.Geometry{
		Type:        geojson.PointType,
		Coordinates: []interface{}{100.0, 50.0},
	}

	s, err := Analyze(g)
	if err != nil {
		t.Fatalf("Analyze error: %v", err)
	}

	if s.GeometryTypes["Point"] != 1 {
		t.Errorf("Expected 1 Point geometry, got %d", s.GeometryTypes["Point"])
	}
}

func TestFormatStats(t *testing.T) {
	s := &Stats{
		FeatureCount:  5,
		GeometryTypes: map[string]int{"Point": 3, "Polygon": 2},
		PropertyKeys:  []string{"name", "pop"},
		PropertyTypes: map[string]map[string]int{
			"name": {"string": 5},
			"pop":  {"number": 5},
		},
	}

	output := FormatStats(s)
	if len(output) == 0 {
		t.Error("Expected non-empty output")
	}
}
