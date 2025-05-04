package models

import "time"

type Event struct {
	Time         time.Time
	EventID      int
	CompetitorID int
	ExtraParams  string
}

type Competitor struct {
	ID               int
	Registered       bool
	PlannedStartTime time.Time
	ActualStartTime  time.Time
	FinishTime       time.Time
	Status           string

	CurrentLap    int
	LapStartTimes []time.Time
	LapEndTimes   []time.Time

	PenaltyStartTime time.Time
	PenaltyEndTime   time.Time
	PenaltyTimes     []time.Duration
	TotalPenaltyTime time.Duration

	Shots        int
	Hits         int
	CurrentRange int

	NotFinishReason string
}

type LapResult struct {
	Time  time.Duration
	Speed float64
}

type PenaltyResult struct {
	Time  time.Duration
	Speed float64
}

const (
	EventRegistered      = 1
	EventStartTimeSet    = 2
	EventOnStartLine     = 3
	EventStarted         = 4
	EventOnFiringRange   = 5
	EventTargetHit       = 6
	EventLeftFiringRange = 7
	EventEnteredPenalty  = 8
	EventLeftPenalty     = 9
	EventEndedLap        = 10
	EventCantContinue    = 11

	EventDisqualified = 32
	EventFinished     = 33
)

const (
	StatusStarted     = "Started"
	StatusNotStarted  = "NotStarted"
	StatusNotFinished = "NotFinished"
	StatusFinished    = "Finished"
)
