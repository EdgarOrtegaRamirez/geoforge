// Package geometry provides geometric operations on GeoJSON geometries.
package geometry

import (
	"math"

	"github.com/EdgarOrtegaRamirez/geoforge/pkg/geojson"
)

// BBox represents a bounding box [west, south, east, north].
type BBox struct {
	MinLon, MinLat, MaxLon, MaxLat float64
}

// Centroid represents the geometric center.
type Centroid struct {
	Lon, Lat float64
}

// CalculateBBox calculates the bounding box for any GeoJSON object.
func CalculateBBox(obj interface{}) (BBox, bool) {
	switch v := obj.(type) {
	case *geojson.FeatureCollection:
		return bboxFeatureCollection(v)
	case *geojson.Feature:
		return bboxFeature(v)
	case *geojson.Geometry:
		return bboxGeometry(v)
	}
	return BBox{}, false
}

func bboxFeatureCollection(fc *geojson.FeatureCollection) (BBox, bool) {
	if len(fc.Features) == 0 {
		return BBox{}, false
	}
	bb := BBox{MinLon: math.MaxFloat64, MinLat: math.MaxFloat64, MaxLon: -math.MaxFloat64, MaxLat: -math.MaxFloat64}
	for _, f := range fc.Features {
		if fb, ok := bboxFeature(f); ok {
			bb.MinLon = math.Min(bb.MinLon, fb.MinLon)
			bb.MinLat = math.Min(bb.MinLat, fb.MinLat)
			bb.MaxLon = math.Max(bb.MaxLon, fb.MaxLon)
			bb.MaxLat = math.Max(bb.MaxLat, fb.MaxLat)
		}
	}
	return bb, true
}

func bboxFeature(f *geojson.Feature) (BBox, bool) {
	if f.Geometry == nil {
		return BBox{}, false
	}
	return bboxGeometry(f.Geometry)
}

func bboxGeometry(g *geojson.Geometry) (BBox, bool) {
	if g == nil {
		return BBox{}, false
	}

	bb := BBox{MinLon: math.MaxFloat64, MinLat: math.MaxFloat64, MaxLon: -math.MaxFloat64, MaxLat: -math.MaxFloat64}

	switch g.Type {
	case geojson.PointType:
		if coords, ok := g.Coordinates.([]interface{}); ok && len(coords) >= 2 {
			lon := coords[0].(float64)
			lat := coords[1].(float64)
			return BBox{MinLon: lon, MinLat: lat, MaxLon: lon, MaxLat: lat}, true
		}
	case geojson.MultiPointType, geojson.LineStringType:
		if coords, ok := g.Coordinates.([]interface{}); ok {
			for _, c := range coords {
				if point, ok := c.([]interface{}); ok && len(point) >= 2 {
					lon := point[0].(float64)
					lat := point[1].(float64)
					bb.MinLon = math.Min(bb.MinLon, lon)
					bb.MinLat = math.Min(bb.MinLat, lat)
					bb.MaxLon = math.Max(bb.MaxLon, lon)
					bb.MaxLat = math.Max(bb.MaxLat, lat)
				}
			}
			return bb, true
		}
	case geojson.PolygonType:
		if rings, ok := g.Coordinates.([]interface{}); ok {
			for _, ring := range rings {
				if points, ok := ring.([]interface{}); ok {
					for _, p := range points {
						if point, ok := p.([]interface{}); ok && len(point) >= 2 {
							lon := point[0].(float64)
							lat := point[1].(float64)
							bb.MinLon = math.Min(bb.MinLon, lon)
							bb.MinLat = math.Min(bb.MinLat, lat)
							bb.MaxLon = math.Max(bb.MaxLon, lon)
							bb.MaxLat = math.Max(bb.MaxLat, lat)
						}
					}
				}
			}
			return bb, true
		}
	case geojson.MultiPolygonType:
		if polygons, ok := g.Coordinates.([]interface{}); ok {
			for _, poly := range polygons {
				if rings, ok := poly.([]interface{}); ok {
					for _, ring := range rings {
						if points, ok := ring.([]interface{}); ok {
							for _, p := range points {
								if point, ok := p.([]interface{}); ok && len(point) >= 2 {
									lon := point[0].(float64)
									lat := point[1].(float64)
									bb.MinLon = math.Min(bb.MinLon, lon)
									bb.MinLat = math.Min(bb.MinLat, lat)
									bb.MaxLon = math.Max(bb.MaxLon, lon)
									bb.MaxLat = math.Max(bb.MaxLat, lat)
								}
							}
						}
					}
				}
			}
			return bb, true
		}
	case geojson.GeometryCollectionType:
		for _, geomRaw := range g.Geometries {
			if sub, ok := geomRaw.(*geojson.Geometry); ok {
				if subBb, ok := bboxGeometry(sub); ok {
					bb.MinLon = math.Min(bb.MinLon, subBb.MinLon)
					bb.MinLat = math.Min(bb.MinLat, subBb.MinLat)
					bb.MaxLon = math.Max(bb.MaxLon, subBb.MaxLon)
					bb.MaxLat = math.Max(bb.MaxLat, subBb.MaxLat)
				}
			}
		}
		return bb, true
	}

	return BBox{}, false
}

// CalculateCentroid calculates the centroid for a geometry.
func CalculateCentroid(g *geojson.Geometry) (Centroid, bool) {
	if g == nil {
		return Centroid{}, false
	}

	switch g.Type {
	case geojson.PointType:
		if coords, ok := g.Coordinates.([]interface{}); ok && len(coords) >= 2 {
			return Centroid{Lon: coords[0].(float64), Lat: coords[1].(float64)}, true
		}
	case geojson.MultiPointType, geojson.LineStringType:
		if coords, ok := g.Coordinates.([]interface{}); ok && len(coords) > 0 {
			sumLon, sumLat := 0.0, 0.0
			count := 0
			for _, c := range coords {
				if point, ok := c.([]interface{}); ok && len(point) >= 2 {
					sumLon += point[0].(float64)
					sumLat += point[1].(float64)
					count++
				}
			}
			if count > 0 {
				return Centroid{Lon: sumLon / float64(count), Lat: sumLat / float64(count)}, true
			}
		}
	case geojson.PolygonType:
		return polygonCentroid(g.Coordinates)
	case geojson.MultiPolygonType:
		if polygons, ok := g.Coordinates.([]interface{}); ok && len(polygons) > 0 {
			return polygonCentroid(polygons[0])
		}
	}
	return Centroid{}, false
}

func polygonCentroid(coords interface{}) (Centroid, bool) {
	if rings, ok := coords.([]interface{}); ok && len(rings) > 0 {
		if ring, ok := rings[0].([]interface{}); ok && len(ring) > 0 {
			sumLon, sumLat := 0.0, 0.0
			for _, p := range ring {
				if point, ok := p.([]interface{}); ok && len(point) >= 2 {
					sumLon += point[0].(float64)
					sumLat += point[1].(float64)
				}
			}
			n := float64(len(ring))
			return Centroid{Lon: sumLon / n, Lat: sumLat / n}, true
		}
	}
	return Centroid{}, false
}

// HaversineDistance calculates the distance in meters between two points.
func HaversineDistance(lon1, lat1, lon2, lat2 float64) float64 {
	const R = 6371000 // Earth radius in meters
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

// ApproximateArea calculates the approximate area in square meters using the Shoelace formula.
func ApproximateArea(coords []geojson.Position) float64 {
	if len(coords) < 3 {
		return 0
	}

	area := 0.0
	n := len(coords)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		// Convert to radians
		lat1 := coords[i][1] * math.Pi / 180
		lat2 := coords[j][1] * math.Pi / 180
		lon1 := coords[i][0] * math.Pi / 180
		lon2 := coords[j][0] * math.Pi / 180
		area += (lon2 - lon1) * (2 + math.Sin(lat1) + math.Sin(lat2))
	}
	area = math.Abs(area * 6378137 * 6378137 / 2) // Earth radius squared
	return area
}

// DouglasPeucker simplifies a line using the Douglas-Peucker algorithm.
func DouglasPeucker(points []geojson.Position, tolerance float64) []geojson.Position {
	if len(points) <= 2 {
		return points
	}

	// Find the point with the maximum distance
	maxDist := 0.0
	maxIdx := 0
	end := len(points) - 1

	for i := 1; i < end; i++ {
		dist := perpendicularDistance(points[i], points[0], points[end])
		if dist > maxDist {
			maxDist = dist
			maxIdx = i
		}
	}

	// If max distance is greater than tolerance, recursively simplify
	if maxDist > tolerance {
		left := DouglasPeucker(points[:maxIdx+1], tolerance)
		right := DouglasPeucker(points[maxIdx:], tolerance)
		return append(left[:len(left)-1], right...)
	}

	return []geojson.Position{points[0], points[end]}
}

func perpendicularDistance(point, lineStart, lineEnd geojson.Position) float64 {
	lon1, lat1 := lineStart[0], lineStart[1]
	lon2, lat2 := lineEnd[0], lineEnd[1]
	lon0, lat0 := point[0], point[1]

	// Calculate distance using Haversine-like approach
	d1 := HaversineDistance(lon1, lat1, lon0, lat0)
	d2 := HaversineDistance(lon2, lat2, lon0, lat0)
	dLine := HaversineDistance(lon1, lat1, lon2, lat2)

	if dLine == 0 {
		return d1
	}

	// Using Heron's formula for triangle area
	s := (d1 + d2 + dLine) / 2
	area := math.Sqrt(s * (s - d1) * (s - d2) * (s - dLine))
	return 2 * area / dLine
}
