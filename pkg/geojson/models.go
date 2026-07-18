// Package geojson provides GeoJSON data models and parsing.
package geojson

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// GeometryType represents the type of a GeoJSON geometry.
type GeometryType string

const (
	PointType              GeometryType = "Point"
	MultiPointType         GeometryType = "MultiPoint"
	LineStringType         GeometryType = "LineString"
	MultiLineStringType    GeometryType = "MultiLineString"
	PolygonType            GeometryType = "Polygon"
	MultiPolygonType       GeometryType = "MultiPolygon"
	GeometryCollectionType GeometryType = "GeometryCollection"
)

// Position represents a geographic coordinate [longitude, latitude, altitude?].
type Position = []float64

// Geometry represents a GeoJSON geometry object.
type Geometry struct {
	Type        GeometryType  `json:"type"`
	Coordinates interface{}   `json:"coordinates"`
	Geometries  []interface{} `json:"geometries,omitempty"` // For GeometryCollection
}

// Feature represents a GeoJSON feature.
type Feature struct {
	Type       string                 `json:"type"`
	Geometry   *Geometry              `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
	ID         interface{}            `json:"id,omitempty"`
	BBox       []float64              `json:"bbox,omitempty"`
}

// FeatureCollection represents a GeoJSON FeatureCollection.
type FeatureCollection struct {
	Type     string     `json:"type"`
	Features []*Feature `json:"features"`
	BBox     []float64  `json:"bbox,omitempty"`
}

// ParseFile reads and parses a GeoJSON file.
func ParseFile(path string) (interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()
	return Parse(f)
}

// Parse reads and parses GeoJSON from a reader.
func Parse(r io.Reader) (interface{}, error) {
	var raw map[string]interface{}
	if err := json.NewDecoder(r).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decoding JSON: %w", err)
	}

	t, ok := raw["type"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'type' field")
	}

	// Re-decode into proper type
	data, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling: %w", err)
	}

	switch t {
	case "FeatureCollection":
		var fc FeatureCollection
		if err := json.Unmarshal(data, &fc); err != nil {
			return nil, fmt.Errorf("parsing FeatureCollection: %w", err)
		}
		return &fc, nil
	case "Feature":
		var f Feature
		if err := json.Unmarshal(data, &f); err != nil {
			return nil, fmt.Errorf("parsing Feature: %w", err)
		}
		return &f, nil
	case "Point", "MultiPoint", "LineString", "MultiLineString", "Polygon", "MultiPolygon", "GeometryCollection":
		var g Geometry
		if err := json.Unmarshal(data, &g); err != nil {
			return nil, fmt.Errorf("parsing geometry: %w", err)
		}
		return &g, nil
	default:
		return nil, fmt.Errorf("unknown GeoJSON type: %s", t)
	}
}

// ParseBytes parses GeoJSON from bytes.
func ParseBytes(data []byte) (interface{}, error) {
	return Parse(bytesReader(data))
}

type bytesReaderWrapper struct {
	data []byte
	pos  int
}

func bytesReader(data []byte) io.Reader {
	return &bytesReaderWrapper{data: data}
}

func (r *bytesReaderWrapper) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n := copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}
