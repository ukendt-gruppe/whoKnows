package models

import (
	"encoding/json"
	"fmt"
	"net/http"
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

const (
	openWeatherMapAPIKey = "4ebea85858f373dfd4ad3339f5a4b91b" // Replace with your actual API key
	weatherAPIURL        = "https://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&APPID=%s"
)

func FetchWeather(city string) (*WeatherData, error) {
	url := fmt.Sprintf(weatherAPIURL, city, openWeatherMapAPIKey)

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
