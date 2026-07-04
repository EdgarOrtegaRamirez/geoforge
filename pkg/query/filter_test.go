package query

import (
	"testing"

	"github.com/EdgarOrtegaRamirez/geoforge/pkg/geojson"
)

func TestFilterFeatures(t *testing.T) {
	features := []*geojson.Feature{
		{
			Type:       "Feature",
			Properties: map[string]interface{}{"name": "New York", "population": 8000000.0},
		},
		{
			Type:       "Feature",
			Properties: map[string]interface{}{"name": "Los Angeles", "population": 4000000.0},
		},
		{
			Type:       "Feature",
			Properties: map[string]interface{}{"name": "Chicago", "population": 2700000.0},
		},
	}

	// Test equality filter
	filters := []Filter{{Key: "name", Operator: "=", Value: "New York"}}
	result := FilterFeatures(features, filters)
	if len(result) != 1 {
		t.Errorf("Expected 1 result, got %d", len(result))
	}

	// Test contains filter
	filters = []Filter{{Key: "name", Operator: "contains", Value: "York"}}
	result = FilterFeatures(features, filters)
	if len(result) != 1 {
		t.Errorf("Expected 1 result, got %d", len(result))
	}

	// Test multiple filters
	filters = []Filter{
		{Key: "name", Operator: "!=", Value: "Chicago"},
		{Key: "population", Operator: ">", Value: "5000000"},
	}
	result = FilterFeatures(features, filters)
	if len(result) != 1 {
		t.Errorf("Expected 1 result, got %d", len(result))
	}
}

func TestParseFilter(t *testing.T) {
	tests := []struct {
		input    string
		wantErr  bool
		wantKey  string
		wantOp   string
		wantVal  string
	}{
		{"name=New York", false, "name", "=", "New York"},
		{"population!=0", false, "population", "!=", "0"},
		{"pop>1000", false, "pop", ">", "1000"},
		{"pop<1000", false, "pop", "<", "1000"},
		{"pop>=1000", false, "pop", ">=", "1000"},
		{"pop<=1000", false, "pop", "<=", "1000"},
		{"name:contains:York", false, "name", "contains", "York"},
		{"name:startswith:New", false, "name", "startswith", "New"},
		{"invalid", true, "", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			f, err := ParseFilter(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFilter(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if f.Key != tt.wantKey {
					t.Errorf("Key = %q, want %q", f.Key, tt.wantKey)
				}
				if f.Operator != tt.wantOp {
					t.Errorf("Operator = %q, want %q", f.Operator, tt.wantOp)
				}
				if f.Value != tt.wantVal {
					t.Errorf("Value = %q, want %q", f.Value, tt.wantVal)
				}
			}
		})
	}
}
