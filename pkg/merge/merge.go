// Package merge provides GeoJSON merging capabilities.
package merge

import (
	"github.com/EdgarOrtegaRamirez/geoforge/pkg/geojson"
)

// MergeFeatureCollections merges multiple FeatureCollections into one.
func MergeFeatureCollections(collections ...*geojson.FeatureCollection) *geojson.FeatureCollection {
	merged := &geojson.FeatureCollection{
		Type:     "FeatureCollection",
		Features: make([]*geojson.Feature, 0),
	}

	for _, fc := range collections {
		if fc != nil {
			merged.Features = append(merged.Features, fc.Features...)
		}
	}

	return merged
}
