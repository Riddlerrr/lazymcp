package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/mark3labs/mcp-go/mcp"
)

type WeatherTool struct {
	Tool    mcp.Tool
	Handler func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

type WeatherData struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int64 `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int64  `json:"sunrise"`
		Sunset  int64  `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

func NewWeatherTool() *WeatherTool {
	return &WeatherTool{
		Tool:    weatherTool(),
		Handler: weatherToolHandler,
	}
}

func weatherTool() mcp.Tool {
	return mcp.NewTool("get_weather",
		mcp.WithDescription("Get current weather for a location. Uses client's IP location by default, or accepts a custom location parameter (city name or 'lat,lon' coordinates)"),
		mcp.WithString("location",
			mcp.Description("Location to get weather for (optional). Can be city name (e.g., 'London' or 'New York,US') or coordinates (e.g., '40.7128,-74.0060'). Uses client IP location if not provided."),
		),
	)
}

func weatherToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Check if API key is configured
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		return mcp.NewToolResultError("OpenWeatherMap API key not configured. Please set OPENWEATHER_API_KEY environment variable. Get your free API key at https://openweathermap.org/api"), nil
	}

	var weatherURL string
	var locationName string
	var units string

	// Check if location parameter is provided
	if args, ok := request.Params.Arguments.(map[string]any); ok {
		if locationParam, exists := args["location"]; exists && locationParam != nil {
			if locationStr, ok := locationParam.(string); ok && locationStr != "" {
				locationName = locationStr
				units = determineUnitsFromLocation(locationStr)
				weatherURL = buildWeatherURLFromLocation(locationStr, apiKey, units)
			}
		}
	}

	// If no location parameter, use IP-based location
	if weatherURL == "" {
		var err error
		weatherURL, locationName, units, err = getWeatherURLFromIP(ctx, apiKey)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get location from IP: %v", err)), nil
		}
	}

	// Fetch weather data
	resp, err := http.Get(weatherURL)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch weather data: %v", err)), nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to read weather response: %v", err)), nil
	}

	if resp.StatusCode != 200 {
		return mcp.NewToolResultError(fmt.Sprintf("Weather API error (status %d): %s", resp.StatusCode, string(body))), nil
	}

	var weatherData WeatherData
	if err := json.Unmarshal(body, &weatherData); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to parse weather response: %v", err)), nil
	}

	result := FormatWeatherAsMarkdown(weatherData, locationName, units)
	return mcp.NewToolResultText(result), nil
}

func buildWeatherURLFromLocation(location, apiKey, units string) string {
	// Check if location is coordinates (lat,lon format)
	if strings.Contains(location, ",") && len(strings.Split(location, ",")) == 2 {
		coords := strings.Split(location, ",")
		lat := strings.TrimSpace(coords[0])
		lon := strings.TrimSpace(coords[1])

		// Validate that both are numbers
		if _, err1 := strconv.ParseFloat(lat, 64); err1 == nil {
			if _, err2 := strconv.ParseFloat(lon, 64); err2 == nil {
				return fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&appid=%s&units=%s", lat, lon, apiKey, units)
			}
		}
	}

	// Treat as city name
	return fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=%s", location, apiKey, units)
}

func getWeatherURLFromIP(ctx context.Context, apiKey string) (string, string, string, error) {
	// Use shared IP data fetching logic
	ipData, err := FetchIPData(ctx, "")
	if err != nil {
		return "", "", "", err
	}

	// Determine units based on country
	units := "metric"
	if ipData.CountryCode == "US" {
		units = "imperial"
	}

	weatherURL := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%.4f&lon=%.4f&appid=%s&units=%s", ipData.Lat, ipData.Lon, apiKey, units)
	locationName := fmt.Sprintf("%s, %s", ipData.City, ipData.Country)

	return weatherURL, locationName, units, nil
}

func determineUnitsFromLocation(location string) string {
	locationUpper := strings.ToUpper(location)

	// Check for Canada first (should be metric, not imperial)
	if strings.Contains(locationUpper, "CANADA") {
		return "metric"
	}

	// Check if location contains US indicators
	usIndicators := []string{",US", ", US", "USA", "UNITED STATES"}

	for _, indicator := range usIndicators {
		if strings.Contains(locationUpper, indicator) {
			return "imperial"
		}
	}

	// Check for US state codes (common ones) - need space or comma before state code
	usStates := []string{", CA", ",CA", " CA", ", NY", ",NY", " NY", ", TX", ",TX", " TX",
		", FL", ",FL", " FL", ", IL", ",IL", " IL", ", PA", ",PA", " PA", ", OH", ",OH", " OH",
		", GA", ",GA", " GA", ", NC", ",NC", " NC", ", MI", ",MI", " MI"}
	for _, state := range usStates {
		if strings.Contains(locationUpper, state) {
			return "imperial"
		}
	}

	return "metric"
}

func FormatWeatherAsMarkdown(data WeatherData, originalLocation string, units string) string {
	var builder strings.Builder
	isImperial := (units == "imperial")

	builder.WriteString(fmt.Sprintf("# Weather Information: %s\n\n", data.Name))

	if originalLocation != "" && originalLocation != data.Name {
		builder.WriteString(fmt.Sprintf("*Requested location: %s*\n\n", originalLocation))
	}

	// Current conditions
	builder.WriteString("## Current Conditions\n")
	if len(data.Weather) > 0 {
		weather := data.Weather[0]
		builder.WriteString(fmt.Sprintf("- **Condition:** %s (%s)\n", toTitle(weather.Description), weather.Main))
	}

	// Temperature formatting based on units
	if isImperial {
		builder.WriteString(fmt.Sprintf("- **Temperature:** %.1f°F (feels like %.1f°F)\n", data.Main.Temp, data.Main.FeelsLike))
		if data.Main.TempMin != data.Main.TempMax {
			builder.WriteString(fmt.Sprintf("- **Range:** %.1f°F - %.1f°F\n", data.Main.TempMin, data.Main.TempMax))
		}
	} else {
		builder.WriteString(fmt.Sprintf("- **Temperature:** %.1f°C (feels like %.1f°C)\n", data.Main.Temp, data.Main.FeelsLike))
		if data.Main.TempMin != data.Main.TempMax {
			builder.WriteString(fmt.Sprintf("- **Range:** %.1f°C - %.1f°C\n", data.Main.TempMin, data.Main.TempMax))
		}
	}

	builder.WriteString(fmt.Sprintf("- **Humidity:** %d%%\n", data.Main.Humidity))
	if isImperial {
		builder.WriteString(fmt.Sprintf("- **Pressure:** %.2f inHg\n", float64(data.Main.Pressure)))
	} else {
		builder.WriteString(fmt.Sprintf("- **Pressure:** %d hPa\n", data.Main.Pressure))
	}

	// Wind and visibility
	builder.WriteString("\n## Details\n")
	if data.Wind.Speed > 0 {
		windDirection := getWindDirection(data.Wind.Deg)
		if isImperial {
			builder.WriteString(fmt.Sprintf("- **Wind:** %.1f mph %s (%d°)\n", data.Wind.Speed, windDirection, data.Wind.Deg))
		} else {
			builder.WriteString(fmt.Sprintf("- **Wind:** %.1f m/s %s (%d°)\n", data.Wind.Speed, windDirection, data.Wind.Deg))
		}
	}
	if data.Visibility > 0 {
		if isImperial {
			visibilityMiles := float64(data.Visibility) / 1609.34 // meters to miles
			builder.WriteString(fmt.Sprintf("- **Visibility:** %.1f miles\n", visibilityMiles))
		} else {
			visibilityKm := float64(data.Visibility) / 1000
			builder.WriteString(fmt.Sprintf("- **Visibility:** %.1f km\n", visibilityKm))
		}
	}
	if data.Clouds.All > 0 {
		builder.WriteString(fmt.Sprintf("- **Cloudiness:** %d%%\n", data.Clouds.All))
	}

	// Location info
	builder.WriteString("\n## Location\n")
	builder.WriteString(fmt.Sprintf("- **City:** %s, %s\n", data.Name, data.Sys.Country))
	builder.WriteString(fmt.Sprintf("- **Coordinates:** %.4f, %.4f\n", data.Coord.Lat, data.Coord.Lon))

	return builder.String()
}

func getWindDirection(degrees int) string {
	directions := []string{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"}
	index := int((float64(degrees)+11.25)/22.5) % 16
	return directions[index]
}

func toTitle(s string) string {
	if s == "" {
		return s
	}

	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			runes := []rune(word)
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}
	return strings.Join(words, " ")
}
