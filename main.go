package main

import (
	"context"
	"log"

	"github.com/Riddlerrr/lazymcp/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func NewMCPServer() *server.MCPServer {
	hooks := &server.Hooks{}

	hooks.AddOnSuccess(func(ctx context.Context, id any, method mcp.MCPMethod, message any, result any) {
		log.Printf("Success: %s, %v, %v, %v\n", method, id, message, result)
	})

	hooks.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error) {
		log.Printf("Error: %s, %v, %v, %v\n", method, id, message, err)
	})

	hooks.AddOnRegisterSession(func(ctx context.Context, session server.ClientSession) {
		log.Printf("onRegisterSession: %s\n", session.SessionID())
	})

	mcpServer := server.NewMCPServer(
		"LazyMCP",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
		server.WithLogging(),
		server.WithHooks(hooks),
	)

	calculator := tools.NewCalculatorTool()
	mcpServer.AddTool(calculator.Tool, calculator.Handler)
	
	return mcpServer
}

func main() {
	s := NewMCPServer()

	// Create HTTP transport server
	httpServer := server.NewStreamableHTTPServer(s)

	// Start the server on port 3000
	log.Printf("HTTP server starting on http://localhost:3000/mcp")
	if err := httpServer.Start(":3000"); err != nil {
		log.Fatal(err)
	}
}
