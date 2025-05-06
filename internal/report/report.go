package report

import (
	"fmt"
	"sort"
	"time"

	"github.com/Gilf4/testTask/internal/config"
	"github.com/Gilf4/testTask/internal/models"
)

func CreateFinalReport(competitors map[models.CompetitorID]*models.Competitor, cfg config.RaceConfig) error {
	// Sort by total time
	sorted := sortCompetitors(competitors)

	for _, c := range sorted {
		line := buildReportLine(c, cfg)
		fmt.Println(line)
	}

	return nil
}

func sortCompetitors(comp map[models.CompetitorID]*models.Competitor) []*models.Competitor {
	var list []*models.Competitor
	for _, c := range comp {
		list = append(list, c)
	}

	sort.Slice(list, func(i, j int) bool {
		if list[i].Disqualified != list[j].Disqualified {
			return !list[i].Disqualified
		}

		ti := list[i].TotalTime
		tj := list[j].TotalTime
		return ti < tj
	})

	return list
}

func buildReportLine(c *models.Competitor, cfg config.RaceConfig) string {
	var report string
	if c.Disqualified {
		report = fmt.Sprintf("[%s]", c.Comment)
	} else {
		diff := c.ActualStart.Sub(c.ScheduledStart)
		report = fmt.Sprintf("[%s]", formatDuration(diff))
	}

	// Laps info
	laps := "["
	for i := 0; i < cfg.Laps; i++ {
		if i < len(c.LapTimes) {
			t := c.LapTimes[i]
			speed := float64(cfg.LapLen) / t.Seconds()
			laps += fmt.Sprintf("{%s, %.4f}", formatDuration(t), speed)
		} else {
			laps += "{,}"
		}
		if i != cfg.Laps-1 {
			laps += ", "
		}
	}
	laps += "]"

	// Penalty time
	avgPenaltySpeed := 0.0
	if c.PenaltyTime > 0 {
		avgPenaltySpeed = float64(cfg.PenaltyLen*((cfg.Laps*5)-c.ShotsHit)) / c.PenaltyTime.Seconds()
	}
	penalty := fmt.Sprintf("{%s, %.3f}", formatDuration(c.PenaltyTime), avgPenaltySpeed)
	totalShots := cfg.Laps * 5

	return fmt.Sprintf("%s %d %s %s %d/%d", report, c.ID, laps, penalty, c.ShotsHit, totalShots)
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	millis := int(d.Milliseconds()) % 1000

	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, millis)
}
