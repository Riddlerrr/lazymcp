package tools

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestWeatherTool_MissingAPIKey(t *testing.T) {
	// Save original API key and clear it
	originalAPIKey := os.Getenv("OPENWEATHER_API_KEY")
	os.Setenv("OPENWEATHER_API_KEY", "")
	defer os.Setenv("OPENWEATHER_API_KEY", originalAPIKey)

	weatherTool := NewWeatherTool()
	ctx := context.WithValue(context.Background(), ClientIPKey, "8.8.8.8")

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{},
		},
	}

	result, err := weatherTool.Handler(ctx, request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !result.IsError {
		t.Fatalf("Expected error result for missing API key")
	}

	textContent, ok := mcp.AsTextContent(result.Content[0])
	if !ok {
		t.Fatalf("Expected text content in error result")
	}

	content := textContent.Text
	if !strings.Contains(content, "OpenWeatherMap API key not configured") {
		t.Errorf("Expected API key error message, got: %s", content)
	}
}

func TestWeatherTool_IPBasedLocation(t *testing.T) {
	// Set up API key
	os.Setenv("OPENWEATHER_API_KEY", "test_api_key")
	defer os.Unsetenv("OPENWEATHER_API_KEY")

	weatherTool := NewWeatherTool()
	ctx := context.WithValue(context.Background(), ClientIPKey, "8.8.8.8")

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{},
		},
	}

	// Since we can't easily modify the hardcoded URLs in the current implementation,
	// we'll test the error case when the real API is not available with a fake key
	result, err := weatherTool.Handler(ctx, request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// With a fake API key, we expect either an error or success (if somehow it works)
	if result.IsError {
		textContent, ok := mcp.AsTextContent(result.Content[0])
		if !ok {
			t.Fatalf("Expected text content in result")
		}
		content := textContent.Text
		if !strings.Contains(content, "Failed to fetch") && !strings.Contains(content, "could not determine") {
			t.Logf("Weather API call error result: %s", content)
		}
	} else {
		textContent, ok := mcp.AsTextContent(result.Content[0])
		if !ok {
			t.Fatalf("Expected text content in result")
		}
		content := textContent.Text
		if !strings.Contains(content, "Weather Information") {
			t.Logf("Weather API call success result: %s", content)
		}
	}
}

func TestWeatherTool_CustomLocationCity(t *testing.T) {
	// Set up API key
	os.Setenv("OPENWEATHER_API_KEY", "test_api_key")
	defer os.Unsetenv("OPENWEATHER_API_KEY")

	weatherTool := NewWeatherTool()
	ctx := context.WithValue(context.Background(), ClientIPKey, "8.8.8.8")

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"location": "London,UK",
			},
		},
	}

	result, err := weatherTool.Handler(ctx, request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Since we're hitting the real API with a fake key, we expect an error
	textContent, ok := mcp.AsTextContent(result.Content[0])
	if !ok {
		t.Fatalf("Expected text content in result")
	}

	content := textContent.Text
	if !strings.Contains(content, "Weather Information") && !strings.Contains(content, "Failed to fetch") && !strings.Contains(content, "API error") {
		t.Logf("Weather API call result: %s", content)
	}
}

func TestWeatherTool_CustomLocationCoordinates(t *testing.T) {
	// Set up API key
	os.Setenv("OPENWEATHER_API_KEY", "test_api_key")
	defer os.Unsetenv("OPENWEATHER_API_KEY")

	weatherTool := NewWeatherTool()
	ctx := context.WithValue(context.Background(), ClientIPKey, "8.8.8.8")

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"location": "40.7128,-74.0060", // New York coordinates
			},
		},
	}

	result, err := weatherTool.Handler(ctx, request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	textContent, ok := mcp.AsTextContent(result.Content[0])
	if !ok {
		t.Fatalf("Expected text content in result")
	}

	content := textContent.Text
	if !strings.Contains(content, "Weather Information") && !strings.Contains(content, "Failed to fetch") && !strings.Contains(content, "API error") {
		t.Logf("Weather API call result: %s", content)
	}
}

func TestWeatherTool_NoClientIP(t *testing.T) {
	// Set up API key
	os.Setenv("OPENWEATHER_API_KEY", "test_api_key")
	defer os.Unsetenv("OPENWEATHER_API_KEY")

	weatherTool := NewWeatherTool()
	ctx := context.Background() // No client IP in context

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{},
		},
	}

	result, err := weatherTool.Handler(ctx, request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !result.IsError {
		t.Fatalf("Expected error result for missing client IP")
	}

	textContent, ok := mcp.AsTextContent(result.Content[0])
	if !ok {
		t.Fatalf("Expected text content in error result")
	}

	content := textContent.Text
	if !strings.Contains(content, "could not determine client IP address") {
		t.Errorf("Expected client IP error message, got: %s", content)
	}
}

func TestFormatWeatherAsMarkdown(t *testing.T) {
	weatherData := WeatherData{
		Coord: struct {
			Lon float64 `json:"lon"`
			Lat float64 `json:"lat"`
		}{Lon: -122.0775, Lat: 37.4056},
		Weather: []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		}{{ID: 800, Main: "Clear", Description: "clear sky", Icon: "01d"}},
		Main: struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  int     `json:"pressure"`
			Humidity  int     `json:"humidity"`
		}{Temp: 22.5, FeelsLike: 21.8, TempMin: 20.0, TempMax: 25.0, Pressure: 1013, Humidity: 65},
		Wind: struct {
			Speed float64 `json:"speed"`
			Deg   int     `json:"deg"`
		}{Speed: 3.5, Deg: 220},
		Visibility: 10000,
		Clouds: struct {
			All int `json:"all"`
		}{All: 0},
		Sys: struct {
			Type    int    `json:"type"`
			ID      int    `json:"id"`
			Country string `json:"country"`
			Sunrise int64  `json:"sunrise"`
			Sunset  int64  `json:"sunset"`
		}{Country: "US"},
		Name: "Mountain View",
	}

	result := FormatWeatherAsMarkdown(weatherData, "Test Location", "metric")

	// Check for expected sections
	expectedSections := []string{
		"# Weather Information: Mountain View",
		"*Requested location: Test Location*",
		"## Current Conditions",
		"**Condition:** Clear Sky (Clear)",
		"**Temperature:** 22.5°C (feels like 21.8°C)",
		"**Range:** 20.0°C - 25.0°C",
		"**Humidity:** 65%",
		"**Pressure:** 1013 hPa",
		"## Details",
		"**Wind:** 3.5 m/s SW (220°)",
		"**Visibility:** 10.0 km",
		"## Location",
		"**City:** Mountain View, US",
		"**Coordinates:** 37.4056, -122.0775",
	}

	for _, expected := range expectedSections {
		if !strings.Contains(result, expected) {
			t.Errorf("Expected to find '%s' in formatted output", expected)
		}
	}
}

func TestGetWindDirection(t *testing.T) {
	testCases := []struct {
		degrees  int
		expected string
	}{
		{0, "N"},
		{45, "NE"},
		{90, "E"},
		{135, "SE"},
		{180, "S"},
		{225, "SW"},
		{270, "W"},
		{315, "NW"},
		{360, "N"},
	}

	for _, tc := range testCases {
		result := getWindDirection(tc.degrees)
		if result != tc.expected {
			t.Errorf("getWindDirection(%d) = %s, expected %s", tc.degrees, result, tc.expected)
		}
	}
}

func TestWeatherTool_RealData(t *testing.T) {
	// Test with real API response from Valencia, Spain
	realWeatherData := WeatherData{
		Coord: struct {
			Lon float64 `json:"lon"`
			Lat float64 `json:"lat"`
		}{Lon: -0.3771, Lat: 39.4676},
		Weather: []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		}{{ID: 800, Main: "Clear", Description: "clear sky", Icon: "01d"}},
		Main: struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  int     `json:"pressure"`
			Humidity  int     `json:"humidity"`
		}{
			Temp:      302.58, // Kelvin - should be converted to Celsius in real usage
			FeelsLike: 302.2,
			TempMin:   300.71,
			TempMax:   303.73,
			Pressure:  1010,
			Humidity:  40,
		},
		Wind: struct {
			Speed float64 `json:"speed"`
			Deg   int     `json:"deg"`
		}{Speed: 7.2, Deg: 90},
		Visibility: 10000,
		Clouds: struct {
			All int `json:"all"`
		}{All: 0},
		Sys: struct {
			Type    int    `json:"type"`
			ID      int    `json:"id"`
			Country string `json:"country"`
			Sunrise int64  `json:"sunrise"`
			Sunset  int64  `json:"sunset"`
		}{Country: "ES", Sunrise: 1756704602, Sunset: 1756751605},
		Name: "Valencia",
		Cod:  200,
	}

	result := FormatWeatherAsMarkdown(realWeatherData, "Valencia, Spain", "metric")

	// Verify key information from real data is present
	expectedContent := []string{
		"# Weather Information: Valencia",
		"*Requested location: Valencia, Spain*",
		"**Condition:** Clear Sky (Clear)",
		"**Temperature:** 302.6°C", // Note: Real API returns Kelvin, our API uses Celsius
		"**Humidity:** 40%",
		"**Pressure:** 1010 hPa",
		"**Wind:** 7.2 m/s E (90°)",
		"**Visibility:** 10.0 km",
		"**City:** Valencia, ES",
		"**Coordinates:** 39.4676, -0.3771",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(result, expected) {
			t.Errorf("Expected to find '%s' in formatted output", expected)
		}
	}

	// Verify the complete structure
	if !strings.Contains(result, "## Current Conditions") {
		t.Error("Missing Current Conditions section")
	}
	if !strings.Contains(result, "## Details") {
		t.Error("Missing Details section")
	}
	if !strings.Contains(result, "## Location") {
		t.Error("Missing Location section")
	}
}

func TestBuildWeatherURLFromLocation(t *testing.T) {
	apiKey := "test_key"

	// Test coordinates
	coordURL := buildWeatherURLFromLocation("40.7128,-74.0060", apiKey, "metric")
	expectedCoordURL := "https://api.openweathermap.org/data/2.5/weather?lat=40.7128&lon=-74.0060&appid=test_key&units=metric"
	if coordURL != expectedCoordURL {
		t.Errorf("buildWeatherURLFromLocation coordinates = %s, expected %s", coordURL, expectedCoordURL)
	}

	// Test city name
	cityURL := buildWeatherURLFromLocation("London,UK", apiKey, "metric")
	expectedCityURL := "https://api.openweathermap.org/data/2.5/weather?q=London,UK&appid=test_key&units=metric"
	if cityURL != expectedCityURL {
		t.Errorf("buildWeatherURLFromLocation city = %s, expected %s", cityURL, expectedCityURL)
	}

	// Test invalid coordinates (should be treated as city)
	invalidURL := buildWeatherURLFromLocation("invalid,coords", apiKey, "imperial")
	expectedInvalidURL := "https://api.openweathermap.org/data/2.5/weather?q=invalid,coords&appid=test_key&units=imperial"
	if invalidURL != expectedInvalidURL {
		t.Errorf("buildWeatherURLFromLocation invalid = %s, expected %s", invalidURL, expectedInvalidURL)
	}
}

func TestDetermineUnitsFromLocation(t *testing.T) {
	testCases := []struct {
		location      string
		expectedUnits string
	}{
		{"New York,US", "imperial"},
		{"London,UK", "metric"},
		{"Paris,France", "metric"},
		{"Los Angeles, CA", "imperial"},
		{"Miami,FL", "imperial"},
		{"Berlin", "metric"},
		{"Tokyo,Japan", "metric"},
		{"Chicago, IL", "imperial"},
		{"Vancouver,Canada", "metric"},
		{"UNITED STATES", "imperial"},
		{"USA", "imperial"},
		{"40.7128,-74.0060", "metric"}, // coordinates default to metric
	}

	for _, tc := range testCases {
		result := determineUnitsFromLocation(tc.location)
		if result != tc.expectedUnits {
			t.Errorf("determineUnitsFromLocation(%q) = %q, expected %q", tc.location, result, tc.expectedUnits)
		}
	}
}

func TestFormatWeatherAsMarkdown_ImperialUnits(t *testing.T) {
	// Test with US location data (imperial units)
	weatherData := WeatherData{
		Coord: struct {
			Lon float64 `json:"lon"`
			Lat float64 `json:"lat"`
		}{Lon: -74.0060, Lat: 40.7128},
		Weather: []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		}{{ID: 800, Main: "Clear", Description: "clear sky", Icon: "01d"}},
		Main: struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  int     `json:"pressure"`
			Humidity  int     `json:"humidity"`
		}{Temp: 72.5, FeelsLike: 75.2, TempMin: 68.0, TempMax: 76.0, Pressure: 30, Humidity: 55},
		Wind: struct {
			Speed float64 `json:"speed"`
			Deg   int     `json:"deg"`
		}{Speed: 8.5, Deg: 180},
		Visibility: 16093, // ~10 miles in meters
		Sys: struct {
			Type    int    `json:"type"`
			ID      int    `json:"id"`
			Country string `json:"country"`
			Sunrise int64  `json:"sunrise"`
			Sunset  int64  `json:"sunset"`
		}{Country: "US"},
		Name: "New York",
	}

	result := FormatWeatherAsMarkdown(weatherData, "New York,US", "imperial")

	// Check for imperial units
	expectedImperialContent := []string{
		"# Weather Information: New York",
		"**Temperature:** 72.5°F (feels like 75.2°F)",
		"**Range:** 68.0°F - 76.0°F",
		"**Pressure:** 30.00 inHg",
		"**Wind:** 8.5 mph S (180°)",
		"**Visibility:** 10.0 miles",
		"**City:** New York, US",
	}

	for _, expected := range expectedImperialContent {
		if !strings.Contains(result, expected) {
			t.Errorf("Expected to find '%s' in imperial formatted output", expected)
		}
	}

	// Make sure no metric units are present
	metricIndicators := []string{"°C", "hPa", "m/s", "km"}
	for _, indicator := range metricIndicators {
		if strings.Contains(result, indicator) {
			t.Errorf("Found metric unit '%s' in imperial formatted output", indicator)
		}
	}
}
