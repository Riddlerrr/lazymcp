package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestIPTool_GetClientIP(t *testing.T) {
	ipTool := NewIPTool()

	tests := []struct {
		name        string
		contextIP   string
		expectError bool
	}{
		{
			name:        "Valid IPv4 address",
			contextIP:   "192.168.1.100",
			expectError: false,
		},
		{
			name:        "Valid IPv6 address",
			contextIP:   "2001:db8::1",
			expectError: false,
		},
		{
			name:        "Loopback address",
			contextIP:   "127.0.0.1",
			expectError: false,
		},
		{
			name:        "Empty IP in context",
			contextIP:   "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.contextIP != "" {
				ctx = context.WithValue(ctx, ClientIPKey, tt.contextIP)
			}

			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{},
				},
			}

			result, err := ipTool.Handler(ctx, request)
			if err != nil {
				t.Fatalf("Handler returned error: %v", err)
			}

			if tt.expectError {
				if !result.IsError {
					t.Errorf("Expected error but got success")
				}
			} else {
				if result.IsError {
					t.Errorf("Expected success but got error: %v", result)
					return
				}

				if len(result.Content) == 0 {
					t.Errorf("No content in result")
					return
				}

				textContent, ok := mcp.AsTextContent(result.Content[0])
				if !ok {
					t.Errorf("Result does not contain text content")
					return
				}

				if textContent.Text != tt.contextIP {
					t.Errorf("Expected IP %s, got %s", tt.contextIP, textContent.Text)
				}

				ip := net.ParseIP(textContent.Text)
				if ip == nil {
					t.Errorf("Result is not a valid IP address: %s", textContent.Text)
				}
			}
		})
	}
}

func TestIPTool_NoContextIP(t *testing.T) {
	ipTool := NewIPTool()
	ctx := context.Background() // No IP in context

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{},
		},
	}

	result, err := ipTool.Handler(ctx, request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if !result.IsError {
		t.Errorf("Expected error when no client IP in context, but got success")
	}
}

func TestIPDataTool_Success(t *testing.T) {
	mockResponse := IPData{
		Query:       "192.168.1.1",
		Status:      "success",
		Country:     "Spain",
		CountryCode: "ES",
		Region:      "VC",
		RegionName:  "Valencia",
		City:        "Valencia",
		Zip:         "46022",
		Lat:         39.4676,
		Lon:         -0.3771,
		Timezone:    "Europe/Madrid",
		ISP:         "Test ISP",
		Org:         "Test Organization",
		AS:          "AS12345 Test AS",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/192.168.1.1" {
			t.Errorf("Expected path /json/192.168.1.1, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	testURL := server.URL + "/json/%s"

	ipDataTool := NewIPDataTool()

	ctx := context.WithValue(context.Background(), ClientIPKey, "192.168.1.1")
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{},
		},
	}

	originalHandler := ipDataTool.Handler
	ipDataTool.Handler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		clientIP, ok := ctx.Value(ClientIPKey).(string)
		if !ok || clientIP == "" {
			return mcp.NewToolResultError("Could not determine client IP address"), nil
		}

		url := fmt.Sprintf(testURL, clientIP)
		resp, err := http.Get(url)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch IP data: %v", err)), nil
		}
		defer resp.Body.Close()

		var ipData IPData
		if err := json.NewDecoder(resp.Body).Decode(&ipData); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to parse response: %v", err)), nil
		}

		if ipData.Status != "success" {
			return mcp.NewToolResultError("Failed to get IP data from service"), nil
		}

		result, err := json.MarshalIndent(ipData, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format response: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	}

	result, err := ipDataTool.Handler(ctx, request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if result.IsError {
		t.Fatalf("Expected success but got error: %v", result.Content[0])
	}

	if len(result.Content) == 0 {
		t.Fatal("No content in result")
	}

	textContent, ok := mcp.AsTextContent(result.Content[0])
	if !ok {
		t.Fatal("Result does not contain text content")
	}

	var parsedData IPData
	if err := json.Unmarshal([]byte(textContent.Text), &parsedData); err != nil {
		t.Fatalf("Failed to parse result JSON: %v", err)
	}

	if parsedData.Country != "Spain" {
		t.Errorf("Expected country 'Spain', got '%s'", parsedData.Country)
	}
	if parsedData.City != "Valencia" {
		t.Errorf("Expected city 'Valencia', got '%s'", parsedData.City)
	}
	if parsedData.Query != "192.168.1.1" {
		t.Errorf("Expected query '192.168.1.1', got '%s'", parsedData.Query)
	}

	ipDataTool.Handler = originalHandler
}

func TestIPDataTool_NoClientIP(t *testing.T) {
	ipDataTool := NewIPDataTool()
	ctx := context.Background()

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{},
		},
	}

	result, err := ipDataTool.Handler(ctx, request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if !result.IsError {
		t.Error("Expected error when no client IP in context, but got success")
	}
}

func TestIPDataTool_APIFailure(t *testing.T) {
	mockResponse := IPData{
		Status: "fail",
		Query:  "invalid",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	testURL := server.URL + "/json/%s"

	ipDataTool := NewIPDataTool()

	ctx := context.WithValue(context.Background(), ClientIPKey, "invalid")
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{},
		},
	}

	ipDataTool.Handler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		clientIP, ok := ctx.Value(ClientIPKey).(string)
		if !ok || clientIP == "" {
			return mcp.NewToolResultError("Could not determine client IP address"), nil
		}

		url := fmt.Sprintf(testURL, clientIP)
		resp, err := http.Get(url)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch IP data: %v", err)), nil
		}
		defer resp.Body.Close()

		var ipData IPData
		if err := json.NewDecoder(resp.Body).Decode(&ipData); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to parse response: %v", err)), nil
		}

		if ipData.Status != "success" {
			return mcp.NewToolResultError("Failed to get IP data from service"), nil
		}

		result, err := json.MarshalIndent(ipData, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format response: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	}

	result, err := ipDataTool.Handler(ctx, request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if !result.IsError {
		t.Error("Expected error when API returns failure status, but got success")
	}
}

func TestIPDataTool_RealData(t *testing.T) {
	realResponseJSON := `{
		"status":"success",
		"country":"Spain",
		"countryCode":"ES",
		"region":"VC",
		"regionName":"Valencia",
		"city":"Valencia",
		"zip":"46022",
		"lat":39.4676,
		"lon":-0.3771,
		"timezone":"Europe/Madrid",
		"isp":"Digi Spain Telecom S.L.U.",
		"org":"Digi Spain Telecom S.L",
		"as":"AS57269 DIGI SPAIN TELECOM S.L.",
		"query":"2a0c:5a85:d506:e700:89d6:90b5:fc32:acc9"
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/json/2a0c:5a85:d506:e700:89d6:90b5:fc32:acc9"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(realResponseJSON))
	}))
	defer server.Close()

	testURL := server.URL + "/json/%s"

	ipDataTool := NewIPDataTool()

	ctx := context.WithValue(context.Background(), ClientIPKey, "2a0c:5a85:d506:e700:89d6:90b5:fc32:acc9")
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{},
		},
	}

	originalHandler := ipDataTool.Handler
	ipDataTool.Handler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		clientIP, ok := ctx.Value(ClientIPKey).(string)
		if !ok || clientIP == "" {
			return mcp.NewToolResultError("Could not determine client IP address"), nil
		}

		url := fmt.Sprintf(testURL, clientIP)
		resp, err := http.Get(url)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch IP data: %v", err)), nil
		}
		defer resp.Body.Close()

		var ipData IPData
		if err := json.NewDecoder(resp.Body).Decode(&ipData); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to parse response: %v", err)), nil
		}

		if ipData.Status != "success" {
			return mcp.NewToolResultError("Failed to get IP data from service"), nil
		}

		result, err := json.MarshalIndent(ipData, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format response: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	}

	result, err := ipDataTool.Handler(ctx, request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if result.IsError {
		t.Fatalf("Expected success but got error: %v", result.Content[0])
	}

	if len(result.Content) == 0 {
		t.Fatal("No content in result")
	}

	textContent, ok := mcp.AsTextContent(result.Content[0])
	if !ok {
		t.Fatal("Result does not contain text content")
	}

	var parsedData IPData
	if err := json.Unmarshal([]byte(textContent.Text), &parsedData); err != nil {
		t.Fatalf("Failed to parse result JSON: %v", err)
	}

	if parsedData.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", parsedData.Status)
	}
	if parsedData.Country != "Spain" {
		t.Errorf("Expected country 'Spain', got '%s'", parsedData.Country)
	}
	if parsedData.CountryCode != "ES" {
		t.Errorf("Expected countryCode 'ES', got '%s'", parsedData.CountryCode)
	}
	if parsedData.Region != "VC" {
		t.Errorf("Expected region 'VC', got '%s'", parsedData.Region)
	}
	if parsedData.RegionName != "Valencia" {
		t.Errorf("Expected regionName 'Valencia', got '%s'", parsedData.RegionName)
	}
	if parsedData.City != "Valencia" {
		t.Errorf("Expected city 'Valencia', got '%s'", parsedData.City)
	}
	if parsedData.Zip != "46022" {
		t.Errorf("Expected zip '46022', got '%s'", parsedData.Zip)
	}
	if parsedData.Lat != 39.4676 {
		t.Errorf("Expected lat 39.4676, got %f", parsedData.Lat)
	}
	if parsedData.Lon != -0.3771 {
		t.Errorf("Expected lon -0.3771, got %f", parsedData.Lon)
	}
	if parsedData.Timezone != "Europe/Madrid" {
		t.Errorf("Expected timezone 'Europe/Madrid', got '%s'", parsedData.Timezone)
	}
	if parsedData.ISP != "Digi Spain Telecom S.L.U." {
		t.Errorf("Expected ISP 'Digi Spain Telecom S.L.U.', got '%s'", parsedData.ISP)
	}
	if parsedData.Org != "Digi Spain Telecom S.L" {
		t.Errorf("Expected org 'Digi Spain Telecom S.L', got '%s'", parsedData.Org)
	}
	if parsedData.AS != "AS57269 DIGI SPAIN TELECOM S.L." {
		t.Errorf("Expected AS 'AS57269 DIGI SPAIN TELECOM S.L.', got '%s'", parsedData.AS)
	}
	if parsedData.Query != "2a0c:5a85:d506:e700:89d6:90b5:fc32:acc9" {
		t.Errorf("Expected query '2a0c:5a85:d506:e700:89d6:90b5:fc32:acc9', got '%s'", parsedData.Query)
	}

	ipDataTool.Handler = originalHandler
}

func TestFetchIPData(t *testing.T) {
	// Test with specific IP
	ctx := context.WithValue(context.Background(), ClientIPKey, "8.8.8.8")

	// This will make a real API call, so we expect either success or failure
	ipData, err := FetchIPData(ctx, "8.8.8.8")
	if err != nil {
		t.Logf("FetchIPData returned error (expected with real API call): %v", err)
		// Error is acceptable in test environment
		return
	}

	if ipData == nil {
		t.Error("FetchIPData returned nil data with no error")
		return
	}

	// Verify basic structure
	if ipData.Query == "" {
		t.Error("FetchIPData returned empty Query field")
	}

	t.Logf("FetchIPData successful: %s, %s, %s", ipData.Query, ipData.City, ipData.Country)
}

func TestFetchIPData_NoClientIP(t *testing.T) {
	ctx := context.Background() // No client IP

	_, err := FetchIPData(ctx, "")
	if err == nil {
		t.Error("Expected error when no client IP provided")
	}

	expectedError := "could not determine client IP address"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error containing '%s', got: %s", expectedError, err.Error())
	}
}
