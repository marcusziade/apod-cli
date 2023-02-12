package main

import (
	"encoding/json"
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
	apods := getAPODsForLastWeek()

	for _, apod := range apods {
		printAPOD(apod)
	}
}

func getAPODsForLastWeek() []NasaAPOD {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -7)

	url := fmt.Sprintf(
		"https://api.nasa.gov/planetary/apod?api_key=efxF8oNY5ZPU48KF3waxgvnQnmITHxLknZpZz6Q8&start_date=%s&end_date=%s",
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"))

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error retrieving data:", err)
		return nil
	}
	defer resp.Body.Close()

	var apods []NasaAPOD
	err = json.NewDecoder(resp.Body).Decode(&apods)
	if err != nil {
		fmt.Println("Error decoding JSON data:", err)
		return nil
	}

	return apods
}

func printAPOD(apod NasaAPOD) {
	date, err := time.Parse("2006-01-02", apod.Date)
	if err != nil {
		fmt.Printf("Error parsing date: %v", err)
		return
	}
	fmt.Printf("%s\n%s\n%s\n\n", apod.Title, date.Format("January 2, 2006"), apod.URL)
}
