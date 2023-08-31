package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
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

	start, end, downloadOnly := parseArgumentsForDateRange()

	fmt.Println("Fetching APODs...")

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
		if downloadOnly {
			downloadImage(apod.URL, sanitizeFilename(apod.Title))
		} else {
			openBrowser(apod.URL)
		}
	}
}

func parseArgumentsForDateRange() (start, end string, downloadOnly bool) {
	flag.StringVar(&start, "start", "", "start date (YYYY-MM-DD)")
	flag.StringVar(&end, "end", "", "end date (YYYY-MM-DD)")
	flag.BoolVar(&downloadOnly, "download-only", false, "Only download images without opening in browser")
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

func printPrettyFormattedAPOD(apod NasaAPOD) {
	date, err := time.Parse("2006-01-02", apod.Date)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}
	fmt.Printf("%s\n%s\n%s\n\n", apod.Title, date.Format("January 2, 2006"), apod.URL)
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

func downloadImage(url, filename string) {
	dir := "./images/"
	// Make sure the directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading image:", err)
		return
	}
	defer resp.Body.Close()

	file, err := os.Create(dir + sanitizeFilename(filename) + ".jpg")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Error saving image:", err)
	}
}

func sanitizeFilename(filename string) string {
	return strings.Map(func(r rune) rune {
		if r == ' ' || r == '_' || r == '-' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return -1
	}, filename)
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		fmt.Println("Error opening browser", url, ":", err)
	}
}
