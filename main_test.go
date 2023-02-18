package main

import (
	"fmt"
	"io"
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

func TestPrintPrettyFormattedAPOD(t *testing.T) {
	apod := NasaAPOD{Date: "2022-02-12", Title: "Test Title", URL: "https://testurl.com"}

	// Redirect stdout to buffer to capture output
	old := os.Stdout
	defer func() { os.Stdout = old }()
	r, w, _ := os.Pipe()
	os.Stdout = w

	printPrettyFormattedAPOD(apod)

	// Read from buffer to check output
	w.Close()
	out, _ := io.ReadAll(r)

	expected := "Test Title\nFebruary 12, 2022\nhttps://testurl.com\n\n"
	if string(out) != expected {
		t.Errorf("Unexpected output. Expected: %q, got: %q", expected, string(out))
	}
}

func TestPrintPrettyFormattedAPODError(t *testing.T) {
	apod := NasaAPOD{Date: "invalid date", Title: "Test Title", URL: "https://testurl.com"}

	err := printPrettyFormattedAPOD(apod)
	if err == nil {
		t.Errorf("Expected function to return an error, but it didn't")
	}
}

func TestConstructURL(t *testing.T) {
	apiKey := "test-api-key"
	start := "2023-02-06"
	end := "2023-02-13"

	expectedURL := fmt.Sprintf(
		"%s?api_key=%s&start_date=%s&end_date=%s",
		apiURL, apiKey, start, end)
	resultURL := constructURL(apiKey, start, end)

	if resultURL != expectedURL {
		t.Errorf("URL does not match expected value. Got %s, expected %s", resultURL, expectedURL)
	}

	// Test case for empty start and end dates
	expectedURL = fmt.Sprintf(
		"%s?api_key=%s&start_date=%s&end_date=%s",
		apiURL, apiKey,
		time.Now().AddDate(0, 0, -7).Format("2006-01-02"),
		time.Now().Format("2006-01-02"))
	resultURL = constructURL(apiKey, "", "")

	if resultURL != expectedURL {
		t.Errorf("URL does not match expected value. Got %s, expected %s", resultURL, expectedURL)
	}
}
