package main

import (
	"log"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"LazyMCP",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithLogging(),
		server.WithRecovery(),
	)

	AddCalculatorTool(s)

	// Create HTTP transport server
	httpServer := server.NewStreamableHTTPServer(s)

	// Start the server on port 3000
	log.Printf("HTTP server starting on http://localhost:3000/mcp")
	if err := httpServer.Start(":3000"); err != nil {
		log.Fatal(err)
	}
}
