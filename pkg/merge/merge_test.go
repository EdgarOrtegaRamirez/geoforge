package merge_test

import (
	"testing"

	"github.com/EdgarOrtegaRamirez/geoforge/pkg/geojson"
	"github.com/EdgarOrtegaRamirez/geoforge/pkg/merge"
)

func TestMergeFeatureCollections(t *testing.T) {
	fc1 := &geojson.FeatureCollection{
		Type:     "FeatureCollection",
		Features: []*geojson.Feature{{Type: "Feature", ID: "a"}, {Type: "Feature", ID: "b"}},
	}
	fc2 := &geojson.FeatureCollection{
		Type:     "FeatureCollection",
		Features: []*geojson.Feature{{Type: "Feature", ID: "c"}},
	}

	merged := merge.MergeFeatureCollections(fc1, fc2)

	if merged.Type != "FeatureCollection" {
		t.Errorf("expected FeatureCollection, got %s", merged.Type)
	}
	if len(merged.Features) != 3 {
		t.Errorf("expected 3 features, got %d", len(merged.Features))
	}
	if merged.Features[0].ID != "a" || merged.Features[2].ID != "c" {
		t.Errorf("unexpected feature order: %v", merged.Features)
	}
}

func TestMergeWithNil(t *testing.T) {
	fc1 := &geojson.FeatureCollection{
		Type:     "FeatureCollection",
		Features: []*geojson.Feature{{Type: "Feature", ID: "x"}},
	}

	merged := merge.MergeFeatureCollections(fc1, nil, nil)

	if len(merged.Features) != 1 {
		t.Errorf("expected 1 feature after merging with nils, got %d", len(merged.Features))
	}
}

func TestMergeEmpty(t *testing.T) {
	merged := merge.MergeFeatureCollections()

	if merged == nil {
		t.Fatal("expected non-nil result for empty merge")
	}
	if len(merged.Features) != 0 {
		t.Errorf("expected 0 features for empty merge, got %d", len(merged.Features))
	}
}

func TestMergeOnlyNils(t *testing.T) {
	merged := merge.MergeFeatureCollections(nil, nil)

	if merged == nil {
		t.Fatal("expected non-nil result for nil-only merge")
	}
	if len(merged.Features) != 0 {
		t.Errorf("expected 0 features, got %d", len(merged.Features))
	}
}
