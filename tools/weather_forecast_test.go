package tools

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// loadFixture loads test data from a JSON fixture file
func loadFixture(t *testing.T, filename string) []byte {
	t.Helper()

	fixturesDir := filepath.Join("..", "test", "fixtures")
	filePath := filepath.Join(fixturesDir, filename)

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to load fixture %s: %v", filename, err)
	}

	return data
}

// loadForecastFixture loads and unmarshals forecast data from a fixture file
func loadForecastFixture(t *testing.T, filename string) ForecastData {
	t.Helper()

	data := loadFixture(t, filename)

	var forecastData ForecastData
	err := json.Unmarshal(data, &forecastData)
	if err != nil {
		t.Fatalf("Failed to unmarshal fixture %s: %v", filename, err)
	}

	return forecastData
}

func TestFormatForecastAsMarkdown(t *testing.T) {
	// Load real API response data from fixture
	forecastData := loadForecastFixture(t, "valencia_forecast_response.json")

	t.Run("MetricUnits", func(t *testing.T) {
		result := FormatForecastAsMarkdown(forecastData, "", "metric")

		// Check header
		if !strings.Contains(result, "# Weather Forecast: Valencia") {
			t.Error("Missing forecast header")
		}

		// Check location info
		if !strings.Contains(result, "Valencia, ES") {
			t.Error("Missing city and country info")
		}

		if !strings.Contains(result, "39.4739, -0.3797") {
			t.Error("Missing coordinates")
		}

		// Check Next 24 Hours section
		if !strings.Contains(result, "## Next 24 Hours") {
			t.Error("Missing Next 24 Hours section")
		}

		// Check temperature in Celsius
		if !strings.Contains(result, "°C") {
			t.Error("Missing Celsius temperature units")
		}

		// Check 5-Day Forecast section
		if !strings.Contains(result, "## 5-Day Forecast") {
			t.Error("Missing 5-Day Forecast section")
		}

		// Check precipitation probability formatting
		if !strings.Contains(result, "20% chance rain") {
			t.Error("Missing precipitation probability formatting")
		}

		// Check weather conditions are properly formatted
		if !strings.Contains(result, "Clear Sky") {
			t.Error("Weather condition not properly title-cased")
		}

		if !strings.Contains(result, "Scattered Clouds") {
			t.Error("Weather condition not properly title-cased")
		}

		if !strings.Contains(result, "Light Rain") {
			t.Error("Weather condition not properly title-cased")
		}
	})

	t.Run("ImperialUnits", func(t *testing.T) {
		result := FormatForecastAsMarkdown(forecastData, "Valencia, Spain", "imperial")

		// Check requested location is shown
		if !strings.Contains(result, "*Requested location: Valencia, Spain*") {
			t.Error("Missing requested location info")
		}

		// Check temperature in Fahrenheit
		if !strings.Contains(result, "°F") {
			t.Error("Missing Fahrenheit temperature units")
		}

		// Should not contain Celsius
		if strings.Contains(result, "°C") {
			t.Error("Should not contain Celsius when using imperial units")
		}
	})

	t.Run("StructureValidation", func(t *testing.T) {
		result := FormatForecastAsMarkdown(forecastData, "", "metric")

		// Verify markdown structure
		lines := strings.Split(result, "\n")
		foundSections := make(map[string]bool)

		for _, line := range lines {
			if strings.HasPrefix(line, "# ") {
				foundSections["header"] = true
			}
			if strings.HasPrefix(line, "## Next 24 Hours") {
				foundSections["next24"] = true
			}
			if strings.HasPrefix(line, "## 5-Day Forecast") {
				foundSections["5day"] = true
			}
			if strings.HasPrefix(line, "## Location") {
				foundSections["location"] = true
			}
		}

		expectedSections := []string{"header", "next24", "5day", "location"}
		for _, section := range expectedSections {
			if !foundSections[section] {
				t.Errorf("Missing section: %s", section)
			}
		}
	})
}

func TestForecastDataUnmarshal(t *testing.T) {
	// Test that our structs correctly unmarshal the API response
	data := loadForecastFixture(t, "minimal_forecast_response.json")

	// Validate core fields
	if data.Cod != "200" {
		t.Errorf("Expected cod '200', got '%s'", data.Cod)
	}

	if len(data.List) != 1 {
		t.Errorf("Expected 1 forecast item, got %d", len(data.List))
	}

	item := data.List[0]
	if item.Dt != 1756803600 {
		t.Errorf("Expected dt 1756803600, got %d", item.Dt)
	}

	if item.Main.Temp != 296.71 {
		t.Errorf("Expected temp 296.71, got %f", item.Main.Temp)
	}

	if len(item.Weather) == 0 || item.Weather[0].Main != "Clear" {
		t.Error("Expected weather main 'Clear'")
	}

	if item.Pop != 0.2 {
		t.Errorf("Expected pop 0.2, got %f", item.Pop)
	}

	// Test optional rain field
	if item.Rain == nil || item.Rain.ThreeH != 0.24 {
		t.Error("Expected rain 3h value of 0.24")
	}

	// Validate city info
	if data.City.Name != "Valencia" {
		t.Errorf("Expected city 'Valencia', got '%s'", data.City.Name)
	}

	if data.City.Coord.Lat != 39.4739 {
		t.Errorf("Expected lat 39.4739, got %f", data.City.Coord.Lat)
	}
}

func TestBuildForecastURLFromLocation(t *testing.T) {
	tests := []struct {
		name        string
		location    string
		apiKey      string
		units       string
		expectedURL string
	}{
		{
			name:        "City name",
			location:    "Valencia",
			apiKey:      "test_key",
			units:       "metric",
			expectedURL: "https://api.openweathermap.org/data/2.5/forecast?q=Valencia&appid=test_key&units=metric",
		},
		{
			name:        "City with country",
			location:    "Valencia,ES",
			apiKey:      "test_key",
			units:       "metric",
			expectedURL: "https://api.openweathermap.org/data/2.5/forecast?q=Valencia,ES&appid=test_key&units=metric",
		},
		{
			name:        "Coordinates",
			location:    "39.4676,-0.3771",
			apiKey:      "test_key",
			units:       "metric",
			expectedURL: "https://api.openweathermap.org/data/2.5/forecast?lat=39.4676&lon=-0.3771&appid=test_key&units=metric",
		},
		{
			name:        "Imperial units",
			location:    "New York,US",
			apiKey:      "test_key",
			units:       "imperial",
			expectedURL: "https://api.openweathermap.org/data/2.5/forecast?q=New York,US&appid=test_key&units=imperial",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildForecastURLFromLocation(tt.location, tt.apiKey, tt.units)
			if result != tt.expectedURL {
				t.Errorf("Expected URL: %s\nGot: %s", tt.expectedURL, result)
			}
		})
	}
}

func TestRealAPIExample(t *testing.T) {
	// Test with Valencia coordinates and test API key
	expectedURL := "https://api.openweathermap.org/data/2.5/forecast?lat=39.4676&lon=-0.3771&appid=test_api_key&units=metric"

	result := buildForecastURLFromLocation("39.4676,-0.3771", "test_api_key", "metric")

	if result != expectedURL {
		t.Errorf("Expected URL: %s\nGot: %s", expectedURL, result)
	}
}
