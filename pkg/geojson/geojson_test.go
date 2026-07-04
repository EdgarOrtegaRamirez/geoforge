package geojson

import (
	"strings"
	"testing"
)

func TestParsePoint(t *testing.T) {
	input := `{"type": "Point", "coordinates": [100.0, 0.0]}`
	obj, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g, ok := obj.(*Geometry)
	if !ok {
		t.Fatalf("Expected *Geometry, got %T", obj)
	}
	if g.Type != PointType {
		t.Errorf("Expected PointType, got %s", g.Type)
	}
}

func TestParseFeatureCollection(t *testing.T) {
	input := `{
		"type": "FeatureCollection",
		"features": [
			{
				"type": "Feature",
				"geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
				{"type": "Feature", "geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
				"properties": {"name": "Test Point"}
			}
		]
	}`
	// Invalid JSON - fix it
	input = `{
		"type": "FeatureCollection",
		"features": [
			{
				"type": "Feature",
				"geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
				"properties": {"name": "Test Point"}
			}
		]
	}`
	obj, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	fc, ok := obj.(*FeatureCollection)
	if !ok {
		t.Fatalf("Expected *FeatureCollection, got %T", obj)
	}
	if len(fc.Features) != 1 {
		t.Errorf("Expected 1 feature, got %d", len(fc.Features))
	}
	if fc.Features[0].Properties["name"] != "Test Point" {
		t.Errorf("Expected property name='Test Point', got %v", fc.Features[0].Properties["name"])
	}
}

func TestParseFeature(t *testing.T) {
	input := `{
		"type": "Feature",
		"geometry": {"type": "LineString", "coordinates": [[0,0], [1,1], [2,2]]},
		"properties": {"id": 42}
	}`
	obj, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	f, ok := obj.(*Feature)
	if !ok {
		t.Fatalf("Expected *Feature, got %T", obj)
	}
	if f.Geometry.Type != LineStringType {
		t.Errorf("Expected LineStringType, got %s", f.Geometry.Type)
	}
}

func TestParseInvalidJSON(t *testing.T) {
	input := `not json`
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestParseMissingType(t *testing.T) {
	input := `{"coordinates": [1, 2]}`
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Error("Expected error for missing type")
	}
}

func TestParseUnknownType(t *testing.T) {
	input := `{"type": "Unknown"}`
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Error("Expected error for unknown type")
	}
}

func TestParseBytes(t *testing.T) {
	input := []byte(`{"type": "Point", "coordinates": [100.0, 0.0]}`)
	obj, err := ParseBytes(input)
	if err != nil {
		t.Fatalf("ParseBytes error: %v", err)
	}
	g, ok := obj.(*Geometry)
	if !ok {
		t.Fatalf("Expected *Geometry, got %T", obj)
	}
	if g.Type != PointType {
		t.Errorf("Expected PointType, got %s", g.Type)
	}
}

func TestParsePolygon(t *testing.T) {
	input := `{
		"type": "Polygon",
		"coordinates": [[[0,0], [1,0], [1,1], [0,1], [0,0]]]
	}`
	obj, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g, ok := obj.(*Geometry)
	if !ok {
		t.Fatalf("Expected *Geometry, got %T", obj)
	}
	if g.Type != PolygonType {
		t.Errorf("Expected PolygonType, got %s", g.Type)
	}
}

func TestParseMultiPoint(t *testing.T) {
	input := `{
		"type": "MultiPoint",
		"coordinates": [[0,0], [1,1]]
	}`
	obj, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g, ok := obj.(*Geometry)
	if !ok {
		t.Fatalf("Expected *Geometry, got %T", obj)
	}
	if g.Type != MultiPointType {
		t.Errorf("Expected MultiPointType, got %s", g.Type)
	}
}
