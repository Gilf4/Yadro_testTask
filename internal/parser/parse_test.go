package parser

import (
	"reflect"
	"testing"
	"time"
)

func TestParseEvent_ValidLine(t *testing.T) {
	line := "[10:01:09.000] 3 2"
	expectedTime, _ := time.Parse("15:04:05.000", "10:01:09.000")

	event, err := ParseEvent(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event.Type != 3 {
		t.Errorf("expected Type 3, got %d", event.Type)
	}
	if event.CompetitorID != 2 {
		t.Errorf("expected CompetitorID 2, got %d", event.CompetitorID)
	}
	if !event.Time.Equal(expectedTime) {
		t.Errorf("expected Time %v, got %v", expectedTime, event.Time)
	}
	if len(event.ExtraParams) != 0 {
		t.Errorf("expected no ExtraParams, got %v", event.ExtraParams)
	}
}

func TestParseEvent_WithExtraParams(t *testing.T) {
	line := "[10:01:00.000] 2 5 10:06:00.000"

	event, err := ParseEvent(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedParams := []string{"10:06:00.000"}
	if !reflect.DeepEqual(event.ExtraParams, expectedParams) {
		t.Errorf("expected ExtraParams %v, got %v", expectedParams, event.ExtraParams)
	}
}

func TestParseEvent_InvalidTime(t *testing.T) {
	line := "[10.01.00.000] 1 1"
	_, err := ParseEvent(line)
	if err == nil {
		t.Error("expected error for invalid time, got nil")
	}
}

func TestParseEvent_InvalidType(t *testing.T) {
	line := "[10:01:00.000] one 1"
	_, err := ParseEvent(line)
	if err == nil {
		t.Error("expected error for invalid event type, got nil")
	}
}

func TestParseEvent_InvalidEventLine(t *testing.T) {
	line := "[10:01:00.000] 2"
	_, err := ParseEvent(line)
	if err == nil {
		t.Error("expected error for invalid event line, got nil")
	}
}
