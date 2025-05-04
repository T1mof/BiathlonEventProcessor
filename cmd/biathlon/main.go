package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"

	"YadroTest/internal/config"
	"YadroTest/internal/models"
	"YadroTest/internal/parser"
	"YadroTest/internal/processor"
	"YadroTest/internal/report"
)

func main() {
	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot determine working directory: %v", err)
	}

	configFile := filepath.Join(baseDir, "config.json")
	eventsDir := filepath.Join(baseDir, "events")

	cfg, err := config.Load(configFile)
	if err != nil {
		log.Fatalf("failed to load config %q: %v", configFile, err)
	}

	var allEvents []models.Event
	err = filepath.WalkDir(eventsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			evs, err := parser.ParseEventsFile(path)
			if err != nil {
				return fmt.Errorf("parse events file %q: %w", path, err)
			}
			allEvents = append(allEvents, evs...)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("cannot read events from %q: %v", eventsDir, err)
	}

	sort.Slice(allEvents, func(i, j int) bool {
		return allEvents[i].Time.Before(allEvents[j].Time)
	})

	competitors, logLines := processor.ProcessEvents(cfg, allEvents)

	for _, line := range logLines {
		fmt.Println(line)
	}

	reportLines := report.GenerateReport(cfg, competitors)
	for _, line := range reportLines {
		fmt.Println(line)
	}
}
