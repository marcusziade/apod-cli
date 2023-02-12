package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gen2brain/beeep"
)

type NasaAPOD struct {
	Date  string `json:"date"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

func main() {
	ticker := time.Tick(24 * time.Hour)

	for range ticker {
		if time.Now().Hour() == 18 {
			apod, err := getLatestAPOD()
			if err != nil {
				fmt.Println("Error retrieving APOD:", err)
				continue
			}

			fmt.Printf("%s (%s)\n%s\n\n", apod.Date, apod.Title, apod.URL)

			err = displayNotification(apod)
			if err != nil {
				fmt.Println("Error displaying notification:", err)
				continue
			}
		}
	}
}

// This function is used to retrieve the latest Astronomy Picture of the Day (APOD) from the NASA API.
func getLatestAPOD() (NasaAPOD, error) {
	url := "https://api.nasa.gov/planetary/apod?api_key=efxF8oNY5ZPU48KF3waxgvnQnmITHxLknZpZz6Q8"

	resp, err := http.Get(url)
	if err != nil {
		return NasaAPOD{}, err
	}
	defer resp.Body.Close()

	var apod NasaAPOD
	err = json.NewDecoder(resp.Body).Decode(&apod)
	if err != nil {
		return NasaAPOD{}, err
	}

	return apod, nil
}

// This function is used to display a desktop notification for a given Astronomy Picture of the Day (APOD). The function downloads the image for the APOD, then saves the image data to a temporary file and creates a desktop notification using the beeep.Notify function from the beeep package, displaing the APODs title and the APOD image.
func displayNotification(apod NasaAPOD) error {
	resp, err := http.Get(apod.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	tmpfile, err := os.CreateTemp("", "apod*.jpg")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write(imgData)
	if err != nil {
		return err
	}

	return beeep.Notify("Picture of the day", apod.Title, tmpfile.Name())
}
