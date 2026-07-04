package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/EdgarOrtegaRamirez/geoforge/pkg/convert"
	"github.com/EdgarOrtegaRamirez/geoforge/pkg/geojson"
	"github.com/EdgarOrtegaRamirez/geoforge/pkg/geometry"
	"github.com/EdgarOrtegaRamirez/geoforge/pkg/merge"
	"github.com/EdgarOrtegaRamirez/geoforge/pkg/query"
	"github.com/EdgarOrtegaRamirez/geoforge/pkg/stats"
)

const sampleGeoJSON = `{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "geometry": {"type": "Point", "coordinates": [-74.006, 40.7128]},
      "properties": {"name": "New York", "population": 8336817}
    },
    {
      "type": "Feature",
      "geometry": {"type": "Point", "coordinates": [-118.2437, 34.0522]},
      "properties": {"name": "Los Angeles", "population": 3979576}
    },
    {
      "type": "Feature",
      "geometry": {"type": "Point", "coordinates": [-87.6298, 41.8781]},
      "properties": {"name": "Chicago", "population": 2693976}
    }
  ]
}`

func TestFullPipeline(t *testing.T) {
	// Write sample file
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "cities.geojson")
	if err := os.WriteFile(inputPath, []byte(sampleGeoJSON), 0644); err != nil {
		t.Fatal(err)
	}

	// 1. Parse
	obj, err := geojson.ParseFile(inputPath)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	fc, ok := obj.(*geojson.FeatureCollection)
	if !ok {
		t.Fatalf("Expected FeatureCollection, got %T", obj)
	}

	// 2. Stats
	s, err := stats.Analyze(fc)
	if err != nil {
		t.Fatalf("Stats error: %v", err)
	}
	if s.FeatureCount != 3 {
		t.Errorf("Expected 3 features, got %d", s.FeatureCount)
	}

	// 3. BBox
	bb, ok := geometry.CalculateBBox(fc)
	if !ok {
		t.Fatal("Expected bbox")
	}
	if bb.MinLon > bb.MaxLon {
		t.Error("MinLon should be <= MaxLon")
	}

	// 4. Filter
	filters := []query.Filter{{Key: "name", Operator: "=", Value: "New York"}}
	filtered := query.FilterFeatures(fc.Features, filters)
	if len(filtered) != 1 {
		t.Errorf("Expected 1 filtered feature, got %d", len(filtered))
	}

	// 5. Convert to WKT
	for _, f := range fc.Features {
		if f.Geometry != nil {
			wkt := convert.ToWKT(f.Geometry)
			if wkt == "" {
				t.Error("Expected non-empty WKT")
			}
		}
	}

	// 6. Merge with itself
	merged := merge.MergeFeatureCollections(fc, fc)
	if len(merged.Features) != 6 {
		t.Errorf("Expected 6 merged features, got %d", len(merged.Features))
	}

	// 7. Format stats
	output := stats.FormatStats(s)
	if len(output) == 0 {
		t.Error("Expected non-empty stats output")
	}
}

func TestParseAndFilter(t *testing.T) {
	input := `{
		"type": "FeatureCollection",
		"features": [
			{
				"type": "Feature",
				"geometry": {"type": "Polygon", "coordinates": [[[0,0],[1,0],[1,1],[0,1],[0,0]]]},
				"properties": {"type": "park", "area": 100}
			},
			{
				"type": "Feature",
				"geometry": {"type": "Point", "coordinates": [0.5, 0.5]},
				"properties": {"type": "restaurant", "area": 10}
			}
		]
	}`

	obj, err := geojson.ParseBytes([]byte(input))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	fc := obj.(*geojson.FeatureCollection)
	if len(fc.Features) != 2 {
		t.Fatalf("Expected 2 features, got %d", len(fc.Features))
	}

	// Filter by type
	filters := []query.Filter{{Key: "type", Operator: "=", Value: "park"}}
	parks := query.FilterFeatures(fc.Features, filters)
	if len(parks) != 1 {
		t.Errorf("Expected 1 park, got %d", len(parks))
	}

	// Check polygon bbox
	bb, ok := geometry.CalculateBBox(parks[0].Geometry)
	if !ok {
		t.Fatal("Expected bbox for polygon")
	}
	if bb.MinLon != 0 || bb.MaxLon != 1 {
		t.Errorf("Unexpected bbox lon: %g to %g", bb.MinLon, bb.MaxLon)
	}
}
