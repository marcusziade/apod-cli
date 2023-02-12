package main

import (
	"encoding/json"
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

func TestReadConfig(t *testing.T) {
	t.Run("Valid Config File", func(t *testing.T) {
		// Create a temporary config file
		configFile, err := os.CreateTemp("", "config.json")
		if err != nil {
			t.Errorf("Error creating temporary config file: %v", err)
		}
		defer os.Remove(configFile.Name())

		// Write the config data to the file
		configData := &Config{APIKey: "123456"}
		err = json.NewEncoder(configFile).Encode(configData)
		if err != nil {
			t.Errorf("Error encoding config data: %v", err)
		}
		configFile.Close()

		// Call readConfig with the temporary config file
		config, err := readConfig(configFile.Name())
		if err != nil {
			t.Errorf("Error reading config file: %v", err)
		}

		// Check that the config was read correctly
		if config.APIKey != configData.APIKey {
			t.Errorf("Expected APIKey '%s', got '%s'", configData.APIKey, config.APIKey)
		}
	})

	t.Run("Invalid Config File", func(t *testing.T) {
		// Call readConfig with a non-existent config file
		_, err := readConfig("nonexistent.json")
		if err == nil {
			t.Errorf("Expected an error, but did not get one")
		}
	})
}
