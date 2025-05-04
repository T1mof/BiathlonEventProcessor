package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config содержит настройки соревнования
type Config struct {
	Laps        int    `json:"laps"`
	LapLen      int    `json:"lapLen"`
	PenaltyLen  int    `json:"penaltyLen"`
	FiringLines int    `json:"firingLines"`
	Start       string `json:"start"`
	StartDelta  string `json:"startDelta"`

	StartTime time.Time
	DeltaTime time.Duration
}

// Load загружает конфигурацию из JSON-файла
func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	startTime, err := time.Parse("15:04:05", cfg.Start)
	if err != nil {
		return nil, fmt.Errorf("invalid start time format: %w", err)
	}
	cfg.StartTime = startTime

	deltaTime, err := parseTimeDelta(cfg.StartDelta)
	if err != nil {
		return nil, fmt.Errorf("invalid start delta format: %w", err)
	}
	cfg.DeltaTime = deltaTime

	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// parseTimeDelta парсит строку формата "HH:MM:SS" в Duration
func parseTimeDelta(timeStr string) (time.Duration, error) {
	t, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return 0, err
	}

	hours := time.Duration(t.Hour()) * time.Hour
	minutes := time.Duration(t.Minute()) * time.Minute
	seconds := time.Duration(t.Second()) * time.Second

	return hours + minutes + seconds, nil
}

// validateConfig проверяет корректность значений конфигурации
func validateConfig(cfg *Config) error {
	if cfg.Laps <= 0 {
		return fmt.Errorf("laps must be positive")
	}
	if cfg.LapLen <= 0 {
		return fmt.Errorf("lapLen must be positive")
	}
	if cfg.PenaltyLen <= 0 {
		return fmt.Errorf("penaltyLen must be positive")
	}
	if cfg.FiringLines <= 0 {
		return fmt.Errorf("firingLines must be positive")
	}
	return nil
}
