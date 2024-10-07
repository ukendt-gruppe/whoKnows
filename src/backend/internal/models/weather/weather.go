// File: src/backend/internal/weather/weather.go

package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func FetchWeather(city string) (*WeatherData, error) {
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

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data WeatherData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("error decoding weather data: %v", err)
	}

	return &data, nil
}