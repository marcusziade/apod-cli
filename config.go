package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type Config struct {
	APIKey string
}

func getOrCreateAPIKey() string {
	apiKey := os.Getenv("NASA_API_KEY")
	if apiKey != "" {
		return apiKey
	}

	filePathName := "Keys.json"
	config, err := readConfig(filePathName)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("API key not found in Keys.json")
			fmt.Println("Please sign up for an API key at https://api.nasa.gov/#signUp")
			fmt.Println("Once you have your API key, enter it below:")
			reader := bufio.NewReader(os.Stdin)
			apiKey, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("Error reading API key input: %v", err)
			}

			apiKey = strings.TrimSpace(apiKey)

			config := &Config{APIKey: apiKey}
			err = writeConfig("Keys.json", config)
			if err != nil {
				log.Fatalf("Error saving API key to file: %v", err)
			}

			fmt.Printf("API key saved to %s.\n", filePathName)
			return apiKey
		}

		log.Fatalf("Error reading Keys.json: %v", err)
	}

	apiKey = config.APIKey
	return apiKey
}

// Reads the configuration data from a file.
func readConfig(fileName string) (*Config, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		fmt.Printf("Error decoding config file: %v\n", err)
		return nil, err
	}

	return config, nil
}

// Saves the configuration data to a file.
func writeConfig(fileName string, config *Config) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(config)
	if err != nil {
		return err
	}

	return nil
}
