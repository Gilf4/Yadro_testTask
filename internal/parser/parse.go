package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Gilf4/testTask/internal/models"
)

const layout = "15:04:05.000"

func ParseFile(path string) ([]*models.Event, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var events []*models.Event
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		event, err := ParseEvent(line)
		if err != nil {
			return nil, fmt.Errorf("parse error: %v", err)
		}

		events = append(events, event)
	}

	return events, nil
}

func ParseEvent(line string) (*models.Event, error) {
	splitLine := strings.Split(line, " ")
	if len(splitLine) < 3 {
		return nil, fmt.Errorf("invalid event line: %s", line)
	}

	timeStr := strings.Trim(splitLine[0], "[]")
	eventTime, err := time.Parse(layout, timeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid time format: %v", err)
	}

	eventType, err := strconv.Atoi(splitLine[1])
	if err != nil {
		return nil, fmt.Errorf("invalid event ID: %v", err)
	}

	competitorID, err := strconv.Atoi(splitLine[2])
	if err != nil {
		return nil, fmt.Errorf("invalid competitor ID: %v", err)
	}

	extraParams := splitLine[3:]

	return &models.Event{
		Time:         eventTime,
		Type:         eventType,
		CompetitorID: models.CompetitorID(competitorID),
		ExtraParams:  extraParams,
	}, nil
}
