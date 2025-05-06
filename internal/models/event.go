package models

import (
	"time"
)

type Event struct {
	Time         time.Time
	Type         int
	CompetitorID CompetitorID
	ExtraParams  []string
}
