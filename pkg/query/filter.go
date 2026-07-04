// Package query provides feature querying capabilities.
package query

import (
	"fmt"
	"strings"

	"github.com/EdgarOrtegaRamirez/geoforge/pkg/geojson"
)

// Filter represents a property filter condition.
type Filter struct {
	Key      string
	Operator string // =, !=, >, <, >=, <=, contains, startswith
	Value    string
}

// FilterFeatures filters features based on property conditions.
func FilterFeatures(features []*geojson.Feature, filters []Filter) []*geojson.Feature {
	var result []*geojson.Feature
	for _, f := range features {
		if matchesFilters(f, filters) {
			result = append(result, f)
		}
	}
	return result
}

func matchesFilters(f *geojson.Feature, filters []Filter) bool {
	for _, filter := range filters {
		if !matchesFilter(f, filter) {
			return false
		}
	}
	return true
}

func matchesFilter(f *geojson.Feature, filter Filter) bool {
	val, ok := f.Properties[filter.Key]
	if !ok {
		return false
	}

	valStr := fmt.Sprintf("%v", val)

	switch filter.Operator {
	case "=":
		return valStr == filter.Value
	case "!=":
		return valStr != filter.Value
	case ">":
		return valStr > filter.Value
	case "<":
		return valStr < filter.Value
	case ">=":
		return valStr >= filter.Value
	case "<=":
		return valStr <= filter.Value
	case "contains":
		return strings.Contains(valStr, filter.Value)
	case "startswith":
		return strings.HasPrefix(valStr, filter.Value)
	default:
		return false
	}
}

// ParseFilter parses a filter string like "name=New York" into a Filter.
func ParseFilter(s string) (Filter, error) {
	operators := []string{"!=", ">=", "<=", "=", ">", "<"}
	for _, op := range operators {
		if idx := strings.Index(s, op); idx > 0 {
			return Filter{
				Key:      s[:idx],
				Operator: op,
				Value:    s[idx+len(op):],
			}, nil
		}
	}

	// Check for contains/startswith
	if idx := strings.Index(s, ":contains:"); idx > 0 {
		return Filter{
			Key:      s[:idx],
			Operator: "contains",
			Value:    s[idx+10:],
		}, nil
	}
	if idx := strings.Index(s, ":startswith:"); idx > 0 {
		return Filter{
			Key:      s[:idx],
			Operator: "startswith",
			Value:    s[idx+len(":startswith:"):],
		}, nil
	}

	return Filter{}, fmt.Errorf("invalid filter expression: %s", s)
}
