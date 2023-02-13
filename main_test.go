package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func TestWriteConfig(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "Keys.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	config := &Config{APIKey: "test-api-key"}

	err = writeConfig(tempFile.Name(), config)
	if err != nil {
		t.Errorf("writeConfig() returned error: %v", err)
	}

	data, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	var configFromFile Config
	err = json.Unmarshal(data, &configFromFile)
	if err != nil {
		t.Errorf("Error unmarshaling JSON: %v", err)
	}

	if configFromFile.APIKey != config.APIKey {
		t.Errorf("API key does not match expected value. Got %s, expected %s", configFromFile.APIKey, config.APIKey)
	}
}

func TestWriteConfigFail(t *testing.T) {
	// Create temporary file and make it read-only
	tempFile, err := ioutil.TempFile("", "Keys.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	err = tempFile.Chmod(0400)
	if err != nil {
		t.Fatalf("Error setting file permissions: %v", err)
	}

	// Attempt to write config to read-only file
	config := &Config{APIKey: "test-api-key"}
	err = writeConfig(tempFile.Name(), config)

	if err == nil {
		t.Error("Expected error, but got nil")
	}
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
	out, _ := ioutil.ReadAll(r)

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
