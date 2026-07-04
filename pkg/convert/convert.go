// Package convert provides format conversion capabilities.
package convert

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/EdgarOrtegaRamirez/geoforge/pkg/geojson"
)

// ToCSV converts a GeoJSON FeatureCollection to CSV format.
func ToCSV(fc *geojson.FeatureCollection, w io.Writer, includeGeometry bool) error {
	if len(fc.Features) == 0 {
		return nil
	}

	// Collect all property keys
	keySet := make(map[string]bool)
	for _, f := range fc.Features {
		for k := range f.Properties {
			keySet[k] = true
		}
	}

	// Build header
	var header []string
	if includeGeometry {
		header = append(header, "geometry_type", "longitude", "latitude")
	}
	for k := range keySet {
		header = append(header, k)
	}

	writer := csv.NewWriter(w)
	defer writer.Flush()

	if err := writer.Write(header); err != nil {
		return fmt.Errorf("writing header: %w", err)
	}

	// Write rows
	for _, f := range fc.Features {
		var row []string
		if includeGeometry && f.Geometry != nil {
			row = append(row, string(f.Geometry.Type))
			lon, lat := extractPoint(f.Geometry)
			row = append(row, fmt.Sprintf("%g", lon))
			row = append(row, fmt.Sprintf("%g", lat))
		}
		for _, k := range header {
			if includeGeometry && (k == "geometry_type" || k == "longitude" || k == "latitude") {
				continue
			}
			val := f.Properties[k]
			row = append(row, fmt.Sprintf("%v", val))
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("writing row: %w", err)
		}
	}

	return nil
}

// FromCSV converts a CSV file to a GeoJSON FeatureCollection.
func FromCSV(r io.Reader, lonCol, latCol string) (*geojson.FeatureCollection, error) {
	reader := csv.NewReader(r)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("reading CSV: %w", err)
	}

	if len(records) < 2 {
		return &geojson.FeatureCollection{Type: "FeatureCollection"}, nil
	}

	header := records[0]
	lonIdx, latIdx := -1, -1
	for i, h := range header {
		switch strings.ToLower(h) {
		case strings.ToLower(lonCol):
			lonIdx = i
		case strings.ToLower(latCol):
			latIdx = i
		}
	}

	fc := &geojson.FeatureCollection{Type: "FeatureCollection", Features: make([]*geojson.Feature, 0, len(records)-1)}

	for _, row := range records[1:] {
		f := &geojson.Feature{
			Type:       "Feature",
			Properties: make(map[string]interface{}),
		}

		for i, val := range row {
			if i == lonIdx || i == latIdx {
				continue
			}
			if i < len(header) {
				f.Properties[header[i]] = val
			}
		}

		if lonIdx >= 0 && latIdx >= 0 && lonIdx < len(row) && latIdx < len(row) {
			var lon, lat float64
			fmt.Sscanf(row[lonIdx], "%f", &lon)
			fmt.Sscanf(row[latIdx], "%f", &lat)
			f.Geometry = &geojson.Geometry{
				Type:        geojson.PointType,
				Coordinates: []interface{}{lon, lat},
			}
		}

		fc.Features = append(fc.Features, f)
	}

	return fc, nil
}

// ToWKT converts a GeoJSON geometry to WKT (Well-Known Text) format.
func ToWKT(g *geojson.Geometry) string {
	if g == nil {
		return ""
	}

	switch g.Type {
	case geojson.PointType:
		if coords, ok := g.Coordinates.([]interface{}); ok && len(coords) >= 2 {
			return fmt.Sprintf("POINT(%g %g)", coords[0].(float64), coords[1].(float64))
		}
	case geojson.MultiPointType:
		return multiPointToWKT(g.Coordinates)
	case geojson.LineStringType:
		return lineStringToWKT(g.Coordinates)
	case geojson.MultiLineStringType:
		return multiLineStringToWKT(g.Coordinates)
	case geojson.PolygonType:
		return polygonToWKT(g.Coordinates)
	case geojson.MultiPolygonType:
		return multiPolygonToWKT(g.Coordinates)
	}
	return ""
}

func multiPointToWKT(coords interface{}) string {
	if points, ok := coords.([]interface{}); ok {
		var parts []string
		for _, p := range points {
			if point, ok := p.([]interface{}); ok && len(point) >= 2 {
				parts = append(parts, fmt.Sprintf("%g %g", point[0].(float64), point[1].(float64)))
			}
		}
		return fmt.Sprintf("MULTIPOINT(%s)", strings.Join(parts, ", "))
	}
	return ""
}

func lineStringToWKT(coords interface{}) string {
	if points, ok := coords.([]interface{}); ok {
		var parts []string
		for _, p := range points {
			if point, ok := p.([]interface{}); ok && len(point) >= 2 {
				parts = append(parts, fmt.Sprintf("%g %g", point[0].(float64), point[1].(float64)))
			}
		}
		return fmt.Sprintf("LINESTRING(%s)", strings.Join(parts, ", "))
	}
	return ""
}

func multiLineStringToWKT(coords interface{}) string {
	if lines, ok := coords.([]interface{}); ok {
		var parts []string
		for _, line := range lines {
			parts = append(parts, lineStringToWKT(line))
		}
		return fmt.Sprintf("MULTILINESTRING(%s)", strings.Join(parts, ", "))
	}
	return ""
}

func polygonToWKT(coords interface{}) string {
	if rings, ok := coords.([]interface{}); ok {
		var parts []string
		for _, ring := range rings {
			parts = append(parts, ringToWKT(ring))
		}
		return fmt.Sprintf("POLYGON(%s)", strings.Join(parts, ", "))
	}
	return ""
}

func ringToWKT(ring interface{}) string {
	if points, ok := ring.([]interface{}); ok {
		var parts []string
		for _, p := range points {
			if point, ok := p.([]interface{}); ok && len(point) >= 2 {
				parts = append(parts, fmt.Sprintf("%g %g", point[0].(float64), point[1].(float64)))
			}
		}
		return fmt.Sprintf("(%s)", strings.Join(parts, ", "))
	}
	return ""
}

func multiPolygonToWKT(coords interface{}) string {
	if polygons, ok := coords.([]interface{}); ok {
		var parts []string
		for _, poly := range polygons {
			parts = append(parts, polygonToWKT(poly))
		}
		return fmt.Sprintf("MULTIPOLYGON(%s)", strings.Join(parts, ", "))
	}
	return ""
}

func extractPoint(g *geojson.Geometry) (lon, lat float64) {
	if g == nil {
		return 0, 0
	}
	switch g.Type {
	case geojson.PointType:
		if coords, ok := g.Coordinates.([]interface{}); ok && len(coords) >= 2 {
			return coords[0].(float64), coords[1].(float64)
		}
	case geojson.MultiPointType:
		if coords, ok := g.Coordinates.([]interface{}); ok && len(coords) > 0 {
			if point, ok := coords[0].([]interface{}); ok && len(point) >= 2 {
				return point[0].(float64), point[1].(float64)
			}
		}
	}
	return 0, 0
}

// ToJSON marshals a GeoJSON object to a formatted JSON string.
func ToJSON(obj interface{}) (string, error) {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling JSON: %w", err)
	}
	return string(data), nil
}

// WriteJSON writes a GeoJSON object to a file.
func WriteJSON(obj interface{}, path string) error {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}
