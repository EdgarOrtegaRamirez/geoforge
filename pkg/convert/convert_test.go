package convert

import (
	"bytes"
	"strings"
	"testing"

	"github.com/EdgarOrtegaRamirez/geoforge/pkg/geojson"
)

func TestToCSV(t *testing.T) {
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

	var buf bytes.Buffer
	err := ToCSV(fc, &buf, true)
	if err != nil {
		t.Fatalf("ToCSV error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "geometry_type") {
		t.Error("Expected header to contain geometry_type")
	}
	if !strings.Contains(output, "Test") {
		t.Error("Expected output to contain 'Test'")
	}
}

func TestToWKTPoint(t *testing.T) {
	g := &geojson.Geometry{
		Type:        geojson.PointType,
		Coordinates: []interface{}{100.0, 50.0},
	}
	wkt := ToWKT(g)
	if wkt != "POINT(100 50)" {
		t.Errorf("Expected 'POINT(100 50)', got %q", wkt)
	}
}

func TestToWKTLineString(t *testing.T) {
	g := &geojson.Geometry{
		Type: geojson.LineStringType,
		Coordinates: []interface{}{
			[]interface{}{0.0, 0.0},
			[]interface{}{1.0, 1.0},
		},
	}
	wkt := ToWKT(g)
	if !strings.HasPrefix(wkt, "LINESTRING(") {
		t.Errorf("Expected LINESTRING(...), got %q", wkt)
	}
}

func TestToWKTPolygon(t *testing.T) {
	g := &geojson.Geometry{
		Type: geojson.PolygonType,
		Coordinates: []interface{}{
			[]interface{}{
				[]interface{}{0.0, 0.0},
				[]interface{}{1.0, 0.0},
				[]interface{}{1.0, 1.0},
				[]interface{}{0.0, 0.0},
			},
		},
	}
	wkt := ToWKT(g)
	if !strings.HasPrefix(wkt, "POLYGON(") {
		t.Errorf("Expected POLYGON(...), got %q", wkt)
	}
}

func TestToJSON(t *testing.T) {
	obj := map[string]interface{}{"type": "test"}
	json, err := ToJSON(obj)
	if err != nil {
		t.Fatalf("ToJSON error: %v", err)
	}
	if !strings.Contains(json, "test") {
		t.Error("Expected JSON to contain 'test'")
	}
}
