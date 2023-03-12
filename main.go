package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"
)

type NasaAPOD struct {
	Date  string `json:"date"`
	Title string `json:"title"`
	URL   string `json:"hdurl"`
}

const apiURL = "https://api.nasa.gov/planetary/apod"

func main() {
	apiKey := getOrCreateAPIKey()
	start, end := parseArgumentsForDateRange()

	loadMsg := "Fetching APODs...\n\n"
	fmt.Println(loadMsg)

	apods, err := getAPODsForDateRange(apiKey, start, end)
	if err != nil {
		fmt.Println("Error retrieving data:", err)
		return
	}

	if len(apods) == 0 {
		fmt.Println("No APODs found in the given date range")
		return
	}

	for _, apod := range apods {
		printPrettyFormattedAPOD(apod)
	}
}

func parseArgumentsForDateRange() (start, end string) {
	flag.StringVar(&start, "start", "", "start date (YYYY-MM-DD)")
	flag.StringVar(&end, "end", "", "end date (YYYY-MM-DD)")
	flag.Parse()
	return
}

/*
This function is used to retrieve Astronomy Picture of the Day (APOD) data for a given date range from the NASA API.
It takes in three parameters: apiKey, start, and end, which represent the API key and the date range to retrieve.
If either start or end is empty, the function will retrieve the APODs for the last week.
*/
func getAPODsForDateRange(apiKey, start, end string) ([]NasaAPOD, error) {
	url := constructURL(apiKey, start, end)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apods []NasaAPOD
	err = json.NewDecoder(resp.Body).Decode(&apods)
	if err != nil {
		return nil, err
	}

	return apods, nil
}

func printPrettyFormattedAPOD(apod NasaAPOD) error {
	date, err := time.Parse("2006-01-02", apod.Date)
	if err != nil {
		return fmt.Errorf("error parsing date: %v", err)
	}
	fmt.Printf("%s\n%s\n%s\n\n", apod.Title, date.Format("January 2, 2006"), apod.URL)
	return nil
}

func constructURL(apiKey, start, end string) string {
	if start == "" || end == "" {
		endDate := time.Now()
		startDate := endDate.AddDate(0, 0, -7)
		return fmt.Sprintf(
			"%s?api_key=%s&start_date=%s&end_date=%s",
			apiURL, apiKey,
			startDate.Format("2006-01-02"),
			endDate.Format("2006-01-02"))
	}
	return fmt.Sprintf(
		"%s?api_key=%s&start_date=%s&end_date=%s",
		apiURL, apiKey,
		start, end)
}
