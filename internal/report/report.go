package report

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"YadroTest/internal/config"
	"YadroTest/internal/models"
)

// GenerateReport создает итоговый отчет о результатах соревнования
func GenerateReport(cfg *config.Config, competitors map[int]*models.Competitor) []string {
	var result []string
	var sortedCompetitors []*models.Competitor

	for _, comp := range competitors {
		if comp.Registered {
			sortedCompetitors = append(sortedCompetitors, comp)
		}
	}

	sortCompetitors(sortedCompetitors)

	for _, comp := range sortedCompetitors {
		line := formatCompetitorResult(comp, cfg)
		result = append(result, line)
	}

	return result
}

// sortCompetitors сортирует участников по времени/статусу
func sortCompetitors(competitors []*models.Competitor) {
	sort.Slice(competitors, func(i, j int) bool {
		ci, cj := competitors[i], competitors[j]

		if ci.Status == models.StatusFinished && cj.Status != models.StatusFinished {
			return true
		}
		if ci.Status != models.StatusFinished && cj.Status == models.StatusFinished {
			return false
		}

		if ci.Status == models.StatusFinished && cj.Status == models.StatusFinished {
			timeI := getTotalTime(ci)
			timeJ := getTotalTime(cj)
			return timeI < timeJ
		}
		return ci.ID < cj.ID
	})
}

// getTotalTime рассчитывает общее время участника
func getTotalTime(comp *models.Competitor) time.Duration {
	if comp.Status != models.StatusFinished {
		return 0
	}

	return comp.FinishTime.Sub(comp.ActualStartTime)
}

// formatCompetitorResult форматирует результат участника
func formatCompetitorResult(comp *models.Competitor, cfg *config.Config) string {
	var sb strings.Builder

	switch comp.Status {
	case models.StatusFinished:
		totalTime := getTotalTime(comp)
		sb.WriteString(fmt.Sprintf("[%s] %d", formatDuration(totalTime), comp.ID))
	case models.StatusNotStarted:
		sb.WriteString(fmt.Sprintf("[NotStarted] %d", comp.ID))
	default:
		sb.WriteString(fmt.Sprintf("[NotFinished] %d", comp.ID))
	}

	sb.WriteString(" [")
	for i := 0; i < len(comp.LapEndTimes); i++ {
		if i > 0 {
			sb.WriteString(", ")
		}

		if !comp.LapEndTimes[i].IsZero() && !comp.LapStartTimes[i].IsZero() {
			lapTime := comp.LapEndTimes[i].Sub(comp.LapStartTimes[i])
			lapSpeed := float64(cfg.LapLen) / lapTime.Seconds()
			lapSpeed = math.Floor(lapSpeed*1000) / 1000
			sb.WriteString(fmt.Sprintf("{%s, %.3f}", formatDuration(lapTime), lapSpeed))
		} else {
			sb.WriteString("{,}")
		}
	}
	sb.WriteString("]")

	if comp.TotalPenaltyTime > 0 {
		penaltySpeed := float64(len(comp.PenaltyTimes)*cfg.PenaltyLen) / comp.TotalPenaltyTime.Seconds()
		sb.WriteString(fmt.Sprintf(" {%s, %.3f}", formatDuration(comp.TotalPenaltyTime), penaltySpeed))
	} else {
		sb.WriteString(" {,}")
	}

	sb.WriteString(fmt.Sprintf(" %d/%d", comp.Hits, comp.Shots))

	return sb.String()
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	ms := int(d.Milliseconds()) % 1000

	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}
