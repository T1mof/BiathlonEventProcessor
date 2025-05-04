package processor

import (
	"testing"
	"time"

	"YadroTest/internal/config"
	"YadroTest/internal/models"
)

func TestDisqualification(t *testing.T) {
	cfg := &config.Config{
		Laps:        2,
		LapLen:      1000,
		PenaltyLen:  50,
		FiringLines: 1,
		DeltaTime:   30 * time.Second,
	}

	plannedStart, err := time.Parse("15:04:05.000", "10:00:00.000")
	if err != nil {
		t.Fatalf("invalid planned start time: %v", err)
	}

	actualStart, err := time.Parse("15:04:05.000", "10:00:31.000")
	if err != nil {
		t.Fatalf("invalid actual start time: %v", err)
	}

	competitor := createNewCompetitor(1, cfg.Laps)
	competitor.PlannedStartTime = plannedStart

	event := models.Event{
		Time:         actualStart,
		EventID:      models.EventStarted,
		CompetitorID: 1,
	}

	_, outEvents := processEvent(cfg, competitor, event)

	if len(outEvents) != 1 {
		t.Errorf("Expected 1 outgoing event, got %d", len(outEvents))
	}

	if competitor.Status != models.StatusNotStarted {
		t.Errorf("Expected status %q, got %q", models.StatusNotStarted, competitor.Status)
	}

	want := "[" + actualStart.Format("15:04:05.000") + "] The competitor(1) is disqualified"
	if outEvents[0] != want {
		t.Errorf("Expected outgoing event %q, got %q", want, outEvents[0])
	}
}
