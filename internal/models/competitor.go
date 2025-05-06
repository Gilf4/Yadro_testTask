package models

import "time"

type CompetitorID int
type Competitor struct {
	ID               CompetitorID
	ScheduledStart   time.Time
	ActualStart      time.Time
	LapStartTimes    []time.Time
	Disqualified     bool
	Comment          string
	CurrentLap       int
	LapTimes         []time.Duration
	TotalTime        time.Duration
	PenaltyEnterTime time.Time
	PenaltyTime      time.Duration
	ShotsHit         int
}
