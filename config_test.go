package main

import (
	"encoding/json"
	"os"
	"testing"
)

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
	tempFile, err := os.CreateTemp("", "Keys.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	config := &Config{APIKey: "test-api-key"}

	err = writeConfig(tempFile.Name(), config)
	if err != nil {
		t.Errorf("writeConfig() returned error: %v", err)
	}

	data, err := os.ReadFile(tempFile.Name())
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
	tempFile, err := os.CreateTemp("", "Keys.json")
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
