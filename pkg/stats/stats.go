// Package stats provides statistical analysis of GeoJSON data.
package stats

import (
	"fmt"
	"sort"
	"strings"

	"github.com/EdgarOrtegaRamirez/geoforge/pkg/geojson"
	"github.com/EdgarOrtegaRamirez/geoforge/pkg/geometry"
)

// Stats represents statistical information about a GeoJSON dataset.
type Stats struct {
	FeatureCount   int                       `json:"feature_count"`
	GeometryTypes  map[string]int            `json:"geometry_types"`
	PropertyKeys   []string                  `json:"property_keys"`
	PropertyTypes  map[string]map[string]int `json:"property_types"`
	BBox           *geometry.BBox            `json:"bbox,omitempty"`
	Centroid       *geometry.Centroid        `json:"centroid,omitempty"`
	SampleFeatures map[string]interface{}    `json:"sample_features,omitempty"`
}

// Analyze computes statistics for a GeoJSON object.
func Analyze(obj interface{}) (*Stats, error) {
	stats := &Stats{
		GeometryTypes: make(map[string]int),
		PropertyTypes: make(map[string]map[string]int),
	}

	switch v := obj.(type) {
	case *geojson.FeatureCollection:
		stats.FeatureCount = len(v.Features)
		stats.PropertyKeys = make([]string, 0)
		keySet := make(map[string]bool)

		for _, f := range v.Features {
			if f.Geometry != nil {
				stats.GeometryTypes[string(f.Geometry.Type)]++
			}
			for k, val := range f.Properties {
				if !keySet[k] {
					stats.PropertyKeys = append(stats.PropertyKeys, k)
					keySet[k] = true
				}
				if stats.PropertyTypes[k] == nil {
					stats.PropertyTypes[k] = make(map[string]int)
				}
				stats.PropertyTypes[k][valType(val)]++
			}
		}

		sort.Strings(stats.PropertyKeys)

		if bb, ok := geometry.CalculateBBox(v); ok {
			stats.BBox = &bb
		}

		// Calculate centroid from bounding box center
		if stats.BBox != nil {
			c := geometry.Centroid{
				Lon: (stats.BBox.MinLon + stats.BBox.MaxLon) / 2,
				Lat: (stats.BBox.MinLat + stats.BBox.MaxLat) / 2,
			}
			stats.Centroid = &c
		}

	case *geojson.Feature:
		stats.FeatureCount = 1
		if v.Geometry != nil {
			stats.GeometryTypes[string(v.Geometry.Type)]++
		}
		if v.Properties != nil {
			stats.PropertyKeys = make([]string, 0, len(v.Properties))
			for k, val := range v.Properties {
				stats.PropertyKeys = append(stats.PropertyKeys, k)
				if stats.PropertyTypes[k] == nil {
					stats.PropertyTypes[k] = make(map[string]int)
				}
				stats.PropertyTypes[k][valType(val)]++
			}
			sort.Strings(stats.PropertyKeys)
		}

	case *geojson.Geometry:
		stats.GeometryTypes[string(v.Type)]++
	}

	return stats, nil
}

// FormatStats returns a human-readable string representation of stats.
func FormatStats(s *Stats) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Features:    %d\n", s.FeatureCount))
	sb.WriteString("Geometry Types:\n")
	for gtype, count := range s.GeometryTypes {
		sb.WriteString(fmt.Sprintf("  %-20s %d\n", gtype, count))
	}

	if s.BBox != nil {
		sb.WriteString(fmt.Sprintf("\nBounding Box:\n"))
		sb.WriteString(fmt.Sprintf("  West:  %g\n", s.BBox.MinLon))
		sb.WriteString(fmt.Sprintf("  South: %g\n", s.BBox.MinLat))
		sb.WriteString(fmt.Sprintf("  East:  %g\n", s.BBox.MaxLon))
		sb.WriteString(fmt.Sprintf("  North: %g\n", s.BBox.MaxLat))
	}

	if s.Centroid != nil {
		sb.WriteString(fmt.Sprintf("\nCentroid: %g, %g\n", s.Centroid.Lon, s.Centroid.Lat))
	}

	if len(s.PropertyKeys) > 0 {
		sb.WriteString(fmt.Sprintf("\nProperties (%d):\n", len(s.PropertyKeys)))
		for _, k := range s.PropertyKeys {
			types := s.PropertyTypes[k]
			typeStrs := make([]string, 0, len(types))
			for t, c := range types {
				typeStrs = append(typeStrs, fmt.Sprintf("%s:%d", t, c))
			}
			sb.WriteString(fmt.Sprintf("  %-25s %s\n", k, strings.Join(typeStrs, ", ")))
		}
	}

	return sb.String()
}

func valType(v interface{}) string {
	switch v.(type) {
	case string:
		return "string"
	case float64:
		return "number"
	case bool:
		return "boolean"
	case nil:
		return "null"
	case map[string]interface{}:
		return "object"
	case []interface{}:
		return "array"
	default:
		return "unknown"
	}
}
