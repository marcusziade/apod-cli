package main

import (
	"os"
	"testing"
	"time"
)

func TestParseArgumentsForDateRange(t *testing.T) {
	t.Run("Start and end flags provided", func(t *testing.T) {
		// Both start and end flags provided, so the function should return the start and end dates
		os.Args = []string{"cmd", "-start", "2023-02-01", "-end", "2023-02-28"}
		start, end := parseArgumentsForDateRange()
		if start != "2023-02-01" {
			t.Errorf("Expected start date to be '2023-02-01', but got %s", start)
		}
		if end != "2023-02-28" {
			t.Errorf("Expected end date to be '2023-02-28', but got %s", end)
		}
	})
}

func TestGetAPODsForDateRange_Success(t *testing.T) {
	start := "2023-02-05"
	end := "2023-02-11"

	apods, err := getAPODsForDateRange(start, end)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(apods) == 0 {
		t.Error("No APODs returned")
	}

	for _, apod := range apods {
		if date, err := time.Parse("2006-01-02", apod.Date); err != nil || date.Before(time.Date(2023, 2, 5, 0, 0, 0, 0, time.UTC)) || date.After(time.Date(2023, 2, 11, 0, 0, 0, 0, time.UTC)) {
			t.Errorf("APOD date out of range: %v", apod.Date)
		}
	}
}

func TestGetAPODsForDateRange_Fail(t *testing.T) {
	start := time.Now().AddDate(0, 0, 1).Format("2006-01-02") // Future date
	end := time.Now().AddDate(0, 0, 2).Format("2006-01-02")

	apods, err := getAPODsForDateRange(start, end)

	if err == nil {
		t.Error("Expected error but got none")
	}

	if apods != nil {
		t.Error("Expected nil APODs but got non-nil")
	}
}
