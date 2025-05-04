package processor

import (
	"fmt"
	"strconv"
	"time"

	"YadroTest/internal/config"
	"YadroTest/internal/models"
)

// ProcessEvents обрабатывает все события и возвращает состояние участников и лог
func ProcessEvents(cfg *config.Config, events []models.Event) (map[int]*models.Competitor, []string) {
	competitors := make(map[int]*models.Competitor)
	logLines := []string{}

	for _, event := range events {
		competitor, exists := competitors[event.CompetitorID]
		if !exists {
			competitor = createNewCompetitor(event.CompetitorID, cfg.Laps)
			competitors[event.CompetitorID] = competitor
		}

		log, outEvents := processEvent(cfg, competitor, event)
		logLines = append(logLines, log)

		logLines = append(logLines, outEvents...)
	}

	return competitors, logLines
}

// createNewCompetitor создает нового участника с начальными значениями
func createNewCompetitor(id int, laps int) *models.Competitor {
	return &models.Competitor{
		ID:            id,
		Status:        "",
		LapStartTimes: make([]time.Time, laps),
		LapEndTimes:   make([]time.Time, laps),
		PenaltyTimes:  []time.Duration{},
	}
}

// processEvent обрабатывает одно событие и возвращает лог и исходящие события
func processEvent(cfg *config.Config, competitor *models.Competitor, event models.Event) (string, []string) {
	var outEvents []string

	logLine := fmt.Sprintf("[%s]", event.Time.Format("15:04:05.000"))

	switch event.EventID {
	case models.EventRegistered:
		competitor.Registered = true
		logLine += fmt.Sprintf(" The competitor(%d) registered", event.CompetitorID)

	case models.EventStartTimeSet:
		startTime, _ := time.Parse("15:04:05.000", event.ExtraParams)
		competitor.PlannedStartTime = startTime
		logLine += fmt.Sprintf(" The start time for the competitor(%d) was set by a draw to %s",
			event.CompetitorID, event.ExtraParams)

	case models.EventOnStartLine:
		logLine += fmt.Sprintf(" The competitor(%d) is on the start line", event.CompetitorID)

	case models.EventStarted:
		competitor.ActualStartTime = event.Time
		competitor.Status = models.StatusStarted
		competitor.CurrentLap = 1
		competitor.LapStartTimes[0] = competitor.PlannedStartTime
		logLine += fmt.Sprintf(" The competitor(%d) has started", event.CompetitorID)

		if event.Time.Sub(competitor.PlannedStartTime) > cfg.DeltaTime {
			competitor.Status = models.StatusNotStarted
			disqEvent := fmt.Sprintf("[%s] The competitor(%d) is disqualified",
				event.Time.Format("15:04:05.000"), event.CompetitorID)
			outEvents = append(outEvents, disqEvent)
		}

	case models.EventOnFiringRange:
		firingRange, _ := strconv.Atoi(event.ExtraParams)
		competitor.CurrentRange = firingRange
		competitor.Shots = cfg.FiringLines * 5
		logLine += fmt.Sprintf(" The competitor(%d) is on the firing range(%s)",
			event.CompetitorID, event.ExtraParams)

	case models.EventTargetHit:
		competitor.Hits++
		logLine += fmt.Sprintf(" The target(%s) has been hit by competitor(%d)",
			event.ExtraParams, event.CompetitorID)

	case models.EventLeftFiringRange:
		logLine += fmt.Sprintf(" The competitor(%d) left the firing range", event.CompetitorID)

	case models.EventEnteredPenalty:
		competitor.PenaltyStartTime = event.Time
		logLine += fmt.Sprintf(" The competitor(%d) entered the penalty laps", event.CompetitorID)

	case models.EventLeftPenalty:
		competitor.PenaltyEndTime = event.Time
		penaltyTime := competitor.PenaltyEndTime.Sub(competitor.PenaltyStartTime)
		competitor.PenaltyTimes = append(competitor.PenaltyTimes, penaltyTime)
		competitor.TotalPenaltyTime += penaltyTime
		logLine += fmt.Sprintf(" The competitor(%d) left the penalty laps", event.CompetitorID)

	case models.EventEndedLap:
		lapIndex := competitor.CurrentLap - 1
		competitor.LapEndTimes[lapIndex] = event.Time
		logLine += fmt.Sprintf(" The competitor(%d) ended the main lap", event.CompetitorID)

		if competitor.CurrentLap >= cfg.Laps {
			competitor.Status = models.StatusFinished
			competitor.FinishTime = event.Time
			finishEvent := fmt.Sprintf("[%s] The competitor(%d) has finished",
				event.Time.Format("15:04:05.000"), event.CompetitorID)
			outEvents = append(outEvents, finishEvent)
		} else {
			competitor.CurrentLap++
			competitor.LapStartTimes[competitor.CurrentLap-1] = event.Time
		}

	case models.EventCantContinue:
		competitor.Status = models.StatusNotFinished
		competitor.NotFinishReason = event.ExtraParams
		logLine += fmt.Sprintf(" The competitor(%d) can`t continue: %s",
			event.CompetitorID, event.ExtraParams)
	}

	return logLine, outEvents
}
