// File: src/backend/tests/integration/weather_test.go

package integration

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/ukendt-gruppe/whoKnows/src/backend/internal/models"
)

func TestWeatherAPIIntegration(t *testing.T) {
	// Determine the environment
	isCI := os.Getenv("CI") == "true"

	// Load .env file for local development
	if !isCI {
		err := godotenv.Load("../../.env")
		if err != nil {
			t.Logf("No .env file found. This is expected in CI environment, but should exist for local development.")
		} else {
			t.Logf("Loaded .env file for local development")
		}
	} else {
		t.Log("Running in CI environment")
	}

	// Check if WEATHER_API_KEY is set
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		t.Fatal("WEATHER_API_KEY is not set")
	}

	// Attempt to fetch weather data
	city := "Copenhagen"
	weather, err := models.FetchWeather(city)
	if err != nil {
		t.Fatalf("Error fetching weather for %s: %v", city, err)
	}

	// Verify that we got a response
	if weather == nil {
		t.Fatal("Received nil weather data")
	}

	// Check if we got data for the correct city
	if weather.Name != city {
		t.Errorf("Expected city name %s, got %s", city, weather.Name)
	}

	// Check if temperature is present (not checking for 0 as it can be a valid temperature)
	if weather.Main.Temp == 0 {
		t.Logf("Warning: Temperature is 0°C. This is possible but unusual.")
	}

	// Check if we have weather conditions
	if len(weather.Weather) == 0 {
		t.Error("Weather conditions array is empty")
	} else if weather.Weather[0].Main == "" || weather.Weather[0].Description == "" {
		t.Error("Weather condition or description is empty")
	}

	// Log the retrieved weather information
	t.Logf("Successfully retrieved weather for %s: %.1f°C, %s (%s)", 
		weather.Name, 
		weather.Main.Temp, 
		weather.Weather[0].Main, 
		weather.Weather[0].Description)

	// Log environment-specific message
	if isCI {
		t.Log("API key successfully used in CI environment")
	} else {
		t.Log("API key successfully used from .env file in local environment")
	}
}