package tools

import (
	"context"
	"net"
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
					Arguments: map[string]interface{}{},
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
			Arguments: map[string]interface{}{},
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