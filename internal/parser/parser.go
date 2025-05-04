package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"YadroTest/internal/models"
)

// ParseEventsFile считывает и парсит события из файла
func ParseEventsFile(filename string) ([]models.Event, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open events file: %w", err)
	}
	defer file.Close()

	var events []models.Event
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		event, err := ParseEventLine(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum, err)
		}
		events = append(events, event)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading events file: %w", err)
	}

	return events, nil
}

// ParseEventLine парсит строку события
func ParseEventLine(line string) (models.Event, error) {
	timeStart := strings.Index(line, "[")
	timeEnd := strings.Index(line, "]")

	if timeStart != 0 || timeEnd <= 0 {
		return models.Event{}, fmt.Errorf("invalid event format: %s", line)
	}

	timeStr := line[1:timeEnd]
	eventTime, err := time.Parse("15:04:05.000", timeStr)
	if err != nil {
		return models.Event{}, fmt.Errorf("invalid time format: %v", err)
	}

	parts := strings.Fields(line[timeEnd+1:])
	if len(parts) < 2 {
		return models.Event{}, fmt.Errorf("insufficient parameters: %s", line)
	}

	eventID, err := strconv.Atoi(parts[0])
	if err != nil {
		return models.Event{}, fmt.Errorf("invalid event ID: %v", err)
	}

	competitorID, err := strconv.Atoi(parts[1])
	if err != nil {
		return models.Event{}, fmt.Errorf("invalid competitor ID: %v", err)
	}

	var extraParams string
	requires := map[int]struct{}{
		models.EventStartTimeSet:  {},
		models.EventOnFiringRange: {},
		models.EventTargetHit:     {},
		models.EventCantContinue:  {},
	}

	if _, need := requires[eventID]; need {
		if len(parts) < 3 {
			return models.Event{}, fmt.Errorf(
				"event %d requires extraParams, but none provided (line=%q)", eventID, line,
			)
		} else {
			extraParams = strings.Join(parts[2:], " ")
		}
	} else if len(parts) > 2 {
		return models.Event{}, fmt.Errorf("event %d does not requires extraParams: %v", eventID, err)
	}

	return models.Event{
		Time:         eventTime,
		EventID:      eventID,
		CompetitorID: competitorID,
		ExtraParams:  extraParams,
	}, nil
}
