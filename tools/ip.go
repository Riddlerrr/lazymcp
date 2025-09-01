package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

type contextKey string

const ClientIPKey contextKey = "client-ip"

type IPTool struct {
	Tool    mcp.Tool
	Handler func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

func NewIPTool() *IPTool {
	return &IPTool{
		Tool:    ipTool(),
		Handler: ipToolHandler,
	}
}

func ipTool() mcp.Tool {
	return mcp.NewTool("get_ip",
		mcp.WithDescription("Get the IP address of the client making the request"),
	)
}

func ipToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract client IP from context
	clientIP, ok := ctx.Value(ClientIPKey).(string)
	if !ok || clientIP == "" {
		return mcp.NewToolResultError("Could not determine client IP address"), nil
	}

	return mcp.NewToolResultText(clientIP), nil
}

type IPDataTool struct {
	Tool    mcp.Tool
	Handler func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

type IPData struct {
	Query       string  `json:"query"`
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	AS          string  `json:"as"`
}

func NewIPDataTool() *IPDataTool {
	return &IPDataTool{
		Tool:    ipDataTool(),
		Handler: ipDataToolHandler,
	}
}

func ipDataTool() mcp.Tool {
	return mcp.NewTool("get_ip_data",
		mcp.WithDescription("Get detailed information about the client's IP address including geolocation data"),
		mcp.WithString("ip",
			mcp.Description("IP address to lookup (optional, uses client IP if not provided)"),
		),
	)
}

// FetchIPData fetches IP geolocation data for a given IP address or client IP from context
func FetchIPData(ctx context.Context, targetIP string) (*IPData, error) {
	// If no target IP provided, get client IP from context
	if targetIP == "" {
		clientIP, ok := ctx.Value(ClientIPKey).(string)
		if !ok || clientIP == "" {
			return nil, fmt.Errorf("could not determine client IP address and no IP parameter provided")
		}
		targetIP = clientIP
	}

	url := fmt.Sprintf("http://ip-api.com/json/%s", targetIP)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch IP data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var ipData IPData
	if err := json.Unmarshal(body, &ipData); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if ipData.Status != "success" {
		return nil, fmt.Errorf("failed to get IP data from service")
	}

	return &ipData, nil
}

func ipDataToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var targetIP string

	// Check if IP parameter is provided
	if args, ok := request.Params.Arguments.(map[string]any); ok {
		if ipParam, exists := args["ip"]; exists && ipParam != nil {
			if ipStr, ok := ipParam.(string); ok && ipStr != "" {
				targetIP = ipStr
			}
		}
	}

	// Fetch IP data using shared function
	ipData, err := FetchIPData(ctx, targetIP)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result := FormatIPDataAsMarkdown(*ipData)
	return mcp.NewToolResultText(result), nil
}

func FormatIPDataAsMarkdown(data IPData) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("# IP Address Information: %s\n\n", data.Query))

	builder.WriteString("## Location\n")
	builder.WriteString(fmt.Sprintf("- **Country:** %s (%s)\n", data.Country, data.CountryCode))
	builder.WriteString(fmt.Sprintf("- **Region:** %s (%s)\n", data.RegionName, data.Region))
	builder.WriteString(fmt.Sprintf("- **City:** %s\n", data.City))
	if data.Zip != "" {
		builder.WriteString(fmt.Sprintf("- **ZIP Code:** %s\n", data.Zip))
	}
	builder.WriteString(fmt.Sprintf("- **Coordinates:** %.4f, %.4f\n", data.Lat, data.Lon))
	builder.WriteString(fmt.Sprintf("- **Timezone:** %s\n\n", data.Timezone))

	builder.WriteString("## Network Information\n")
	builder.WriteString(fmt.Sprintf("- **ISP:** %s\n", data.ISP))
	builder.WriteString(fmt.Sprintf("- **Organization:** %s\n", data.Org))
	builder.WriteString(fmt.Sprintf("- **AS:** %s\n", data.AS))

	return builder.String()
}
