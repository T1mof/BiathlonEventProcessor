package report

import (
	"strings"
	"testing"
	"time"

	"YadroTest/internal/config"
	"YadroTest/internal/models"
)

func TestFormatCompetitorResult(t *testing.T) {
	cfg := &config.Config{
		Laps:       2,
		LapLen:     1000,
		PenaltyLen: 50,
	}

	t.Run("Finished competitor", func(t *testing.T) {
		comp := &models.Competitor{
			ID:     1,
			Status: models.StatusFinished,
			Hits:   4,
			Shots:  5,
		}

		startTime, _ := time.Parse("15:04:05.000", "10:00:00.000")
		finishTime, _ := time.Parse("15:04:05.000", "10:20:00.000")
		comp.ActualStartTime = startTime
		comp.FinishTime = finishTime

		comp.LapStartTimes = make([]time.Time, 2)
		comp.LapEndTimes = make([]time.Time, 2)

		lap1Start, _ := time.Parse("15:04:05.000", "10:00:00.000")
		lap1End, _ := time.Parse("15:04:05.000", "10:10:00.000")
		comp.LapStartTimes[0] = lap1Start
		comp.LapEndTimes[0] = lap1End

		lap2Start, _ := time.Parse("15:04:05.000", "10:10:00.000")
		lap2End, _ := time.Parse("15:04:05.000", "10:20:00.000")
		comp.LapStartTimes[1] = lap2Start
		comp.LapEndTimes[1] = lap2End

		penaltyTime := 1 * time.Minute
		comp.TotalPenaltyTime = penaltyTime
		comp.PenaltyTimes = []time.Duration{penaltyTime}

		result := formatCompetitorResult(comp, cfg)
		expected := "[00:20:00.000] 1 [{00:10:00.000, 1.666}, {00:10:00.000, 1.666}] {00:01:00.000, 0.833} 4/5"

		if result != expected {
			t.Errorf("Expected: %s\nGot: %s", expected, result)
		}
	})

	t.Run("NotStarted competitor", func(t *testing.T) {
		comp := &models.Competitor{
			ID:     2,
			Status: models.StatusNotStarted,
		}
		result := formatCompetitorResult(comp, cfg)
		if !strings.HasPrefix(result, "[NotStarted]") {
			t.Errorf("Expected NotStarted status, got: %s", result)
		}
	})

	t.Run("NotFinished competitor", func(t *testing.T) {
		comp := &models.Competitor{
			ID:              3,
			Status:          models.StatusNotFinished,
			NotFinishReason: "Lost in the forest",
		}
		result := formatCompetitorResult(comp, cfg)
		if !strings.HasPrefix(result, "[NotFinished]") {
			t.Errorf("Expected NotFinished status, got: %s", result)
		}
	})
}
