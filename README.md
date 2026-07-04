# GeoForge

A comprehensive GeoJSON processing toolkit with a CLI and Go library.

## Features

- **Validate** — Validate GeoJSON files for correctness
- **Info** — Show statistics about GeoJSON data (feature counts, property types, bounding box)
- **Bounding Box** — Calculate bounding boxes for any GeoJSON object
- **Centroid** — Calculate centroids of features
- **Query** — Filter features by property conditions (equality, comparison, contains, startswith)
- **Merge** — Combine multiple GeoJSON files into one
- **Convert** — Convert between formats (GeoJSON ↔ CSV, WKT, JSON)
- **Simplify** — Simplify LineString geometries using Douglas-Peucker algorithm

## Quick Start

```bash
# Install
go install github.com/EdgarOrtegaRamirez/geoforge/cmd/geoforge@latest

# Or build from source
git clone https://github.com/EdgarOrtegaRamirez/geoforge
cd geoforge
go build -o geoforge ./cmd/geoforge/
```

## Usage

### Validate a GeoJSON file
```bash
geoforge validate cities.geojson
# VALID: *geojson.FeatureCollection
```

### Show statistics
```bash
geoforge info cities.geojson
# Features:    3
# Geometry Types:
#   Point                3
#
# Bounding Box:
#   West:  -118.2437
#   South: 34.0522
#   East:  -74.006
#   North: 41.8781
```

### Calculate bounding box
```bash
geoforge bbox cities.geojson
# West:  -118.2437
# South: 34.0522
# East:  -74.006
# North: 41.8781
```

### Filter features
```bash
geoforge query cities.geojson -f "name=New York"
geoforge query cities.geojson -f "population>1000000"
geoforge query cities.geojson -f "name:contains:York"
```

### Convert to WKT
```bash
geoforge convert cities.geojson --to wkt
# POINT(-74.006 40.7128)
# POINT(-118.2437 34.0522)
```

### Convert to CSV
```bash
geoforge convert cities.geojson --to csv
# geometry_type,longitude,latitude,name,population
# Point,-74.006,40.7128,New York,8336817
```

### Merge files
```bash
geoforge merge cities1.geojson cities2.geojson > merged.geojson
```

### Simplify LineStrings
```bash
geoforge simplify route.geojson --tolerance 0.001 > simplified.geojson
```

## Supported GeoJSON Types

- Point
- MultiPoint
- LineString
- MultiLineString
- Polygon
- MultiPolygon
- GeometryCollection

## Library API

```go
import "github.com/EdgarOrtegaRamirez/geoforge/pkg/geojson"

// Parse a GeoJSON file
obj, err := geojson.ParseFile("data.geojson")

// Parse from bytes
obj, err := geojson.ParseBytes(data)

// Type assert
fc := obj.(*geojson.FeatureCollection)

// Geometry operations
import "github.com/EdgarOrtegaRamirez/geoforge/pkg/geometry"

bbox, ok := geometry.CalculateBBox(fc)
centroid, ok := geometry.CalculateCentroid(feature.Geometry)
distance := geometry.HaversineDistance(lon1, lat1, lon2, lat2)
simplified := geometry.DouglasPeucker(points, tolerance)

// Feature filtering
import "github.com/EdgarOrtegaRamirez/geoforge/pkg/query"

filters := []query.Filter{{Key: "name", Operator: "=", Value: "test"}}
filtered := query.FilterFeatures(fc.Features, filters)

// Statistics
import "github.com/EdgarOrtegaRamirez/geoforge/pkg/stats"

s, _ := stats.Analyze(fc)
fmt.Print(stats.FormatStats(s))
```

## Architecture

```
geoforge/
├── cmd/geoforge/       # CLI entry point (Cobra)
├── pkg/
│   ├── geojson/        # GeoJSON parsing and models
│   ├── geometry/       # Geometry operations (bbox, centroid, distance, simplify)
│   ├── query/          # Feature filtering and querying
│   ├── convert/        # Format conversion (CSV, WKT, JSON)
│   ├── merge/          # Merge multiple GeoJSON files
│   └── stats/          # Statistical analysis
└── tests/              # Integration tests
```

## Filter Syntax

| Operator | Example | Description |
|----------|---------|-------------|
| `=` | `name=New York` | Exact match |
| `!=` | `name!=Chicago` | Not equal |
| `>` | `pop>1000000` | Greater than |
| `<` | `pop<5000000` | Less than |
| `>=` | `pop>=1000000` | Greater or equal |
| `<=` | `pop<=5000000` | Less or equal |
| `:contains:` | `name:contains:New` | Contains substring |
| `:startswith:` | `name:startswith:New` | Starts with |

## License

MIT
