package main

import (
	"context"
	"log"
	"net/http"
	"strings"

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

	ipTool := tools.NewIPTool()
	mcpServer.AddTool(ipTool.Tool, ipTool.Handler)
	
	return mcpServer
}

func getClientIP(r *http.Request) string {
	// Try X-Forwarded-For header first (for proxies/load balancers)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Try X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Try CF-Connecting-IP for Cloudflare
	if cfIP := r.Header.Get("CF-Connecting-IP"); cfIP != "" {
		return cfIP
	}

	// Fall back to RemoteAddr
	// RemoteAddr might be in format "IP:port", so we need to extract just the IP
	if idx := strings.LastIndex(r.RemoteAddr, ":"); idx != -1 {
		return r.RemoteAddr[:idx]
	}
	return r.RemoteAddr
}

func main() {
	s := NewMCPServer()

	// Create HTTP transport server with custom context function
	httpServer := server.NewStreamableHTTPServer(s,
		server.WithHTTPContextFunc(func(ctx context.Context, r *http.Request) context.Context {
			clientIP := getClientIP(r)
			return context.WithValue(ctx, tools.ClientIPKey, clientIP)
		}),
	)

	// Start the server on port 3000
	log.Printf("HTTP server starting on http://localhost:3000/mcp")
	if err := httpServer.Start(":3000"); err != nil {
		log.Fatal(err)
	}
}
