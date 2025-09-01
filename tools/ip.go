package tools

import (
	"context"

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
