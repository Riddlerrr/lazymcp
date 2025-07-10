package main

import (
	"fmt"
	"log"
	"os"
	"time"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

type TimeArgs struct {
	Format   string `json:"format,omitempty" jsonschema:"description=Time format: 'unix', 'iso8601', 'rfc3339', 'kitchen', or custom Go time format,default=iso8601"`
	Timezone string `json:"timezone,omitempty" jsonschema:"description=Timezone (e.g., 'UTC', 'America/New_York', 'Europe/London'),default=UTC"`
}

func main() {
	// Enable debug logging to stderr
	log.SetOutput(os.Stderr)
	log.Println("Starting lazymcp server...")

	// Create a new server using stdio transport
	transport := stdio.NewStdioServerTransport()
	server := mcp_golang.NewServer(transport)

	log.Println("Server created, registering tools...")

	// Register the get_current_time tool
	err := server.RegisterTool("get_current_time", "Get the current date and time in various formats", func(args TimeArgs) (*mcp_golang.ToolResponse, error) {
		log.Printf("Tool called with args: %+v", args)
		// Use defaults if not provided
		format := args.Format
		if format == "" {
			format = "iso8601"
		}

		timezone := args.Timezone
		if timezone == "" {
			timezone = "UTC"
		}

		// Load timezone
		loc, err := time.LoadLocation(timezone)
		if err != nil {
			return nil, fmt.Errorf("invalid timezone: %v", err)
		}

		// Get current time in the specified timezone
		now := time.Now().In(loc)

		// Format time based on the requested format
		var formatted string
		switch format {
		case "unix":
			formatted = fmt.Sprintf("%d", now.Unix())
		case "iso8601", "rfc3339":
			formatted = now.Format(time.RFC3339)
		case "kitchen":
			formatted = now.Format(time.Kitchen)
		default:
			// Try to use custom format
			formatted = now.Format(format)
		}

		// Build detailed response
		responseText := fmt.Sprintf(`Current time in %s: %s

Details:
- Year: %d
- Month: %s
- Day: %d
- Hour: %d
- Minute: %d
- Second: %d
- Weekday: %s
- Timezone: %s
- UTC Offset: %s
- Unix Timestamp: %d`,
			timezone, formatted,
			now.Year(),
			now.Month().String(),
			now.Day(),
			now.Hour(),
			now.Minute(),
			now.Second(),
			now.Weekday().String(),
			loc.String(),
			now.Format("-07:00"),
			now.Unix())

		// Create the response
		response := mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(responseText),
		)

		return response, nil
	})

	if err != nil {
		log.Fatalf("Failed to register tool: %v", err)
	}

	log.Println("Tool registered successfully, starting server...")

	// Start the server - this blocks and handles incoming messages
	if err := server.Serve(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
