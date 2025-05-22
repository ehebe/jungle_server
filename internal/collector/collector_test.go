package collector

import (
	"encoding/json"
	"testing"
)

// TestCollect verifies that Collect returns non-nil and valid CPU/Memory values
func TestCollect(t *testing.T) {
	stats := Collect()

	if stats == nil {
		t.Fatal("Expected non-nil SystemStats")
	}

	if stats.CPU < 0 || stats.CPU > 100 {
		t.Errorf("Unexpected CPU usage value: %f", stats.CPU)
	}

	if stats.Mem < 0 || stats.Mem > 100 {
		t.Errorf("Unexpected Memory usage value: %f", stats.Mem)
	}
}

// TestToJSON verifies that ToJSON outputs valid JSON
func TestToJSON(t *testing.T) {
	stats := &SystemStats{
		CPU: 12.34,
		Mem: 56.78,
	}

	jsonData := stats.ToJSON()
	if len(jsonData) == 0 {
		t.Fatal("Expected non-empty JSON output")
	}

	var parsed SystemStats
	if err := json.Unmarshal(jsonData, &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if parsed.CPU != stats.CPU || parsed.Mem != stats.Mem {
		t.Errorf("Parsed data mismatch: got %+v, expected %+v", parsed, stats)
	}
}
