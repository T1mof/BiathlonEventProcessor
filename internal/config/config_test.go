package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	tempFile, err := os.CreateTemp("", "config.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	configJSON := `{
        "laps": 2,
        "lapLen": 3651,
        "penaltyLen": 50,
        "firingLines": 1,
        "start": "09:30:00",
        "startDelta": "00:00:30"
    }`

	if _, err := tempFile.Write([]byte(configJSON)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	cfg, err := Load(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Laps != 2 {
		t.Errorf("Expected Laps=2, got %d", cfg.Laps)
	}
	if cfg.LapLen != 3651 {
		t.Errorf("Expected LapLen=3651, got %d", cfg.LapLen)
	}
	if cfg.PenaltyLen != 50 {
		t.Errorf("Expected PenaltyLen=50, got %d", cfg.PenaltyLen)
	}
	if cfg.FiringLines != 1 {
		t.Errorf("Expected FiringLines=1, got %d", cfg.FiringLines)
	}

	expectedStart, _ := time.Parse("15:04:05", "09:30:00")
	if !cfg.StartTime.Equal(expectedStart) {
		t.Errorf("Expected StartTime=%v, got %v", expectedStart, cfg.StartTime)
	}

	expectedDelta := 30 * time.Second
	if cfg.DeltaTime != expectedDelta {
		t.Errorf("Expected DeltaTime=%v, got %v", expectedDelta, cfg.DeltaTime)
	}
}
