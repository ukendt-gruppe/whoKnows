// File: src/backend/cmd/weather/main.go

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type WeatherData struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Name string `json:"name"`
}

const weatherAPIURL = "https://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&APPID=%s"

func fetchWeather(city string) (*WeatherData, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("WEATHER_API_KEY environment variable is not set")
	}

	url := fmt.Sprintf(weatherAPIURL, city, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching weather data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error decoding weather data: %v", err)
	}

	return &data, nil
}

func main() {
	city := "Copenhagen"
	weather, err := fetchWeather(city)
	if err != nil {
		log.Fatalf("Error fetching weather for %s: %v", city, err)
	}
	fmt.Printf("Weather in %s: %.1f°C, %s (%s)\n", 
		weather.Name, 
		weather.Main.Temp, 
		weather.Weather[0].Main, 
		weather.Weather[0].Description)
}