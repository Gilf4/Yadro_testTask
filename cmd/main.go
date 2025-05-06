package main

import (
	"fmt"
	"github.com/Gilf4/testTask/internal/config"
	"github.com/Gilf4/testTask/internal/eventsProcessor"
	"github.com/Gilf4/testTask/internal/parser"
	"github.com/Gilf4/testTask/internal/report"
	"log"
)

const (
	configPath = "config/config.json"
	dataPath   = "testData/events"
)

func main() {
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	events, err := parser.ParseFile(dataPath)
	if err != nil {
		fmt.Printf("parse file error: %v\n", err)
		return
	}

	competition, err := eventsProcessor.NewCompetition(cfg)
	if err != nil {
		fmt.Printf("init competition error: %v\n", err)
		return
	}

	if err := competition.ProcessEvents(events); err != nil {
		fmt.Printf("process event error: %v\n", err)
		return
	}

	fmt.Printf("\n==Race report==\n")
	if err := report.CreateFinalReport(competition.CompetitorsList, competition.Config); err != nil {
		fmt.Printf("create final report error: %v\n", err)
	}

	fmt.Printf("\n==Outgoing events==\n")
	for _, v := range competition.OutgoingEvents {
		fmt.Print(v)
	}
}
