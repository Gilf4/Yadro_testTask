package eventsProcessor

import (
	"fmt"
	"github.com/Gilf4/testTask/internal/utils"
	"strings"
	"time"

	"github.com/Gilf4/testTask/internal/config"
	"github.com/Gilf4/testTask/internal/models"
)

type Competition struct {
	Config          config.RaceConfig
	CompetitorsList map[models.CompetitorID]*models.Competitor
	StartTime       time.Time
	StartDelta      time.Duration
	OutgoingEvents  []string
}

func NewCompetition(cfg *config.RaceConfig) (*Competition, error) {
	start, err := time.Parse("15:04:05.000", cfg.Start)
	if err != nil {
		return nil, fmt.Errorf("invalid start time: %w", err)
	}

	startDelta, err := utils.ToDuration(cfg.StartDelta)
	if err != nil {
		return nil, fmt.Errorf("invalid start delta: %w", err)
	}

	return &Competition{
		Config:          *cfg,
		CompetitorsList: make(map[models.CompetitorID]*models.Competitor),
		StartTime:       start,
		StartDelta:      startDelta,
		OutgoingEvents:  make([]string, 0),
	}, nil
}

func (c *Competition) ProcessEvents(events []*models.Event) error {
	for _, event := range events {
		if err := c.processEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (c *Competition) processEvent(event *models.Event) error {
	competitor := c.getCompetitor(event.CompetitorID)
	timeStr := event.Time.Format("15:04:05.000")

	switch event.Type {
	case 1: // Competitor registration
		fmt.Printf("[%s] The competitor(%d) registered\n", timeStr, event.CompetitorID)

	case 2: // Setting scheduled the start time
		startTime, err := time.Parse("15:04:05.000", event.ExtraParams[0])
		if err != nil {
			return fmt.Errorf("invalid start time format: %v", err)
		}
		competitor.ScheduledStart = startTime

		fmt.Printf("[%s] The start time for the competitor(%d) was set by a draw to %s\n",
			timeStr, event.CompetitorID, event.ExtraParams[0])

	case 3: // Competitor on the start line
		fmt.Printf("[%s] The competitor(%d) is on the start line\n", timeStr, event.CompetitorID)

	case 4: //Competitor starting a race
		// Checking for late start
		maxStartTime := competitor.ScheduledStart.Add(c.StartDelta)
		if event.Time.After(maxStartTime) {
			competitor.Disqualified = true
			competitor.Comment = "NotStarted"
			outgoingEvent := fmt.Sprintf("[%s] 32 %d\n", timeStr, event.CompetitorID)
			c.OutgoingEvents = append(c.OutgoingEvents, outgoingEvent)
			return nil
		}

		competitor.ActualStart = event.Time
		competitor.CurrentLap = 1
		competitor.LapStartTimes = append(competitor.LapStartTimes, event.Time)

		fmt.Printf("[%s] The competitor(%d) has started\n", timeStr, event.CompetitorID)

	case 5: // At the firing range
		fmt.Printf("[%s] The competitor(%d) is on the firing range(%s)\n",
			timeStr, event.CompetitorID, event.ExtraParams[0])

	case 6: // Hitting the target
		competitor.ShotsHit++
		fmt.Printf("[%s] The target(%s) has been hit by competitor(%d)\n",
			timeStr, event.ExtraParams[0], event.CompetitorID)

	case 7: // Left the firing range
		fmt.Printf("[%s] The competitor(%d) left the firing range\n", timeStr, event.CompetitorID)

	case 8: // Entering the penalty lap
		competitor.PenaltyEnterTime = event.Time
		fmt.Printf("[%s] The competitor(%d) entered the penalty laps\n", timeStr, event.CompetitorID)

	case 9: // Left the penalty lap
		competitor.PenaltyTime += event.Time.Sub(competitor.PenaltyEnterTime)
		fmt.Printf("[%s] The competitor(%d) left the penalty laps\n", timeStr, event.CompetitorID)

	case 10: // End the main lap
		lapTime := event.Time.Sub(competitor.LapStartTimes[competitor.CurrentLap-1])
		competitor.TotalTime += lapTime
		competitor.LapTimes = append(competitor.LapTimes, lapTime)

		competitor.CurrentLap++
		competitor.LapStartTimes = append(competitor.LapStartTimes, event.Time)
		fmt.Printf("[%s] The competitor(%d) ended the main lap\n", timeStr, event.CompetitorID)

		if competitor.CurrentLap > c.Config.Laps {
			outgoingEvent := fmt.Sprintf("[%s] 33 %d\n", timeStr, competitor.ID)
			c.OutgoingEvents = append(c.OutgoingEvents, outgoingEvent)
		}

	case 11: // Competitor can`t continue
		competitor.Disqualified = true
		competitor.Comment = "NotFinished"

		fmt.Printf("[%s] The competitor(%d) can`t continue: %s\n",
			timeStr, event.CompetitorID, strings.Join(event.ExtraParams, " "))
		outgoingEvent := fmt.Sprintf("[%s] 32 %d\n", timeStr, event.CompetitorID)
		c.OutgoingEvents = append(c.OutgoingEvents, outgoingEvent)

	default:
		return fmt.Errorf("unknown event type: %d", event.Type)
	}

	return nil
}

func (c *Competition) getCompetitor(id models.CompetitorID) *models.Competitor {
	_, ok := c.CompetitorsList[id]
	if !ok {
		c.CompetitorsList[id] = &models.Competitor{
			ID: id,
		}
	}

	return c.CompetitorsList[id]
}
