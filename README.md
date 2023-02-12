# APOD CLI
The APOD CLI is a command-line interface that allows you to fetch and view Astronomy Picture of the Day (APOD) data from NASA's API.

## Features
- Fetch APOD data for a specified date range
- Show APOD title, date, and URL in the terminal
- Optional command-line flags to specify the date range

## Installation
### Download and run (Linux, macOS, Windows) 
Download the executable from [releases](https://github.com/marcusziade/apod-cli/releases)

### Homebrew (macOS)
```
brew install apod-cli
```
### Build and run
1. Make sure you have Go installed on your computer.
2. Clone this repository to your local machine.
3. Run go build in the project directory to build the binary.
4. Optionally, you can move the built binary to a directory in your system PATH to make it available globally.

## Usage
### Basic usage:
```
apod-cli
```
This will fetch APOD data for the last week and display it in the terminal.

### Specifying a date range:
```
apod-cli -start=2022-01-01 -end=2022-01-07
```
This will fetch APOD data for the specified date range (inclusive) and display it in the terminal.

## Contributing
If you find any issues or have a feature request, feel free to create an issue on the GitHub repository. Pull requests are also welcome!

## License
This project is licensed under the MIT License. See the LICENSE file for details.