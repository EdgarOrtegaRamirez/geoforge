# AGENTS.md

## Project Overview

GeoForge is a comprehensive GeoJSON processing toolkit with a CLI and Go library for parsing, analyzing, querying, converting, merging, and simplifying GeoJSON data.

## Building

```bash
go build -o geoforge ./cmd/geoforge/
```

## Testing

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test ./... -v
```

Run specific package tests:

```bash
go test ./pkg/geojson/...
go test ./pkg/geometry/...
go test ./pkg/query/...
go test ./tests/...
```

## Architecture

- `cmd/geoforge/`: CLI entry point using Cobra with 8 commands
- `pkg/geojson/`: GeoJSON parsing, models, and validation
- `pkg/geometry/`: Geometry operations (bbox, centroid, distance, simplification)
- `pkg/query/`: Feature filtering and querying with multiple operators
- `pkg/convert/`: Format conversion (CSV, WKT, JSON)
- `pkg/merge/`: Merge multiple GeoJSON files
- `pkg/stats/`: Statistical analysis of GeoJSON data

## Key Data Structures

- `geojson.FeatureCollection`: Collection of GeoJSON features
- `geojson.Feature`: Individual feature with geometry and properties
- `geojson.Geometry`: Geometry object (Point, LineString, Polygon, etc.)
- `geometry.BBox`: Bounding box (West, South, East, North)
- `geometry.Centroid`: Geographic centroid (Lon, Lat)
- `query.Filter`: Property filter condition

## Key Algorithms

- **Haversine Distance**: Great-circle distance between two points on Earth
- **Douglas-Peucker**: Line simplification algorithm
- **Shoelace Formula**: Polygon area calculation
- **Perpendicular Distance**: For line simplification

## Dependencies

- `github.com/spf13/cobra`: CLI framework
- Standard library only for core geometry operations

## Code Style

- Follow Go conventions
- Use meaningful variable names
- Add comments for exported functions
- Handle errors explicitly
- Use `interface{}` for flexible GeoJSON coordinate handling
