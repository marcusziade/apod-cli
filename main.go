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
	URL   string `json:"url"`
}

func main() {
	var start, end string
	flag.StringVar(&start, "start", "", "start date (YYYY-MM-DD)")
	flag.StringVar(&end, "end", "", "end date (YYYY-MM-DD)")
	flag.Parse()

	loadMsg := "Fetching APODs...\n\n"
	fmt.Println(loadMsg)

	apods, err := getAPODsForDateRange(start, end)
	if err != nil {
		fmt.Println("Error retrieving data:", err)
		return
	}

	if len(apods) == 0 {
		fmt.Println("No APODs found in the given date range")
		return
	}

	for _, apod := range apods {
		printAPOD(apod)
	}
}

func getAPODsForDateRange(start, end string) ([]NasaAPOD, error) {
	var url string
	if start == "" || end == "" {
		endDate := time.Now()
		startDate := endDate.AddDate(0, 0, -7)
		url = fmt.Sprintf(
			"https://api.nasa.gov/planetary/apod?api_key=efxF8oNY5ZPU48KF3waxgvnQnmITHxLknZpZz6Q8&start_date=%s&end_date=%s",
			startDate.Format("2006-01-02"),
			endDate.Format("2006-01-02"))
	} else {
		url = fmt.Sprintf(
			"https://api.nasa.gov/planetary/apod?api_key=efxF8oNY5ZPU48KF3waxgvnQnmITHxLknZpZz6Q8&start_date=%s&end_date=%s",
			start,
			end)
	}

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

func printAPOD(apod NasaAPOD) {
	date, err := time.Parse("2006-01-02", apod.Date)
	if err != nil {
		fmt.Printf("Error parsing date: %v", err)
		return
	}
	fmt.Printf("%s\n%s\n%s\n\n", apod.Title, date.Format("January 2, 2006"), apod.URL)
}
