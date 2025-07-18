package tools

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestCalculatorTool_BasicOperations(t *testing.T) {
	calculator := NewCalculatorTool()
	ctx := context.Background()

	tests := []struct {
		name      string
		operation string
		x         float64
		y         float64
		expected  string
		hasError  bool
	}{
		{"Add positive numbers", "add", 5, 3, "8.00", false},
		{"Add negative numbers", "add", -5, -3, "-8.00", false},
		{"Add mixed numbers", "add", 5, -3, "2.00", false},
		{"Subtract positive numbers", "subtract", 10, 4, "6.00", false},
		{"Subtract negative numbers", "subtract", -10, -4, "-6.00", false},
		{"Multiply positive numbers", "multiply", 6, 7, "42.00", false},
		{"Multiply by zero", "multiply", 5, 0, "0.00", false},
		{"Multiply negative numbers", "multiply", -3, -4, "12.00", false},
		{"Divide positive numbers", "divide", 15, 3, "5.00", false},
		{"Divide by zero", "divide", 10, 0, "", true},
		{"Divide negative numbers", "divide", -12, -3, "4.00", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{
						"operation": tt.operation,
						"x":         tt.x,
						"y":         tt.y,
					},
				},
			}

			result, err := calculator.Handler(ctx, request)
			if err != nil {
				t.Fatalf("Handler returned error: %v", err)
			}

			if tt.hasError {
				if !result.IsError {
					t.Errorf("Expected error but got success")
				}
			} else {
				if result.IsError {
					t.Errorf("Expected success but got error: %v", result.Content)
				}
				if len(result.Content) > 0 {
					if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
						if textContent.Text != tt.expected {
							t.Errorf("Expected %s, got %s", tt.expected, textContent.Text)
						}
					}
				}
			}
		})
	}
}

func TestCalculatorTool_AdvancedOperations(t *testing.T) {
	calculator := NewCalculatorTool()
	ctx := context.Background()

	tests := []struct {
		name      string
		operation string
		x         float64
		y         float64
		expected  string
		hasError  bool
	}{
		{"Power positive base and exponent", "power", 2, 3, "8.00", false},
		{"Power negative base even exponent", "power", -2, 2, "4.00", false},
		{"Power negative base odd exponent", "power", -2, 3, "-8.00", false},
		{"Power zero base", "power", 0, 5, "0.00", false},
		{"Power base to zero", "power", 5, 0, "1.00", false},
		{"Power fractional exponent", "power", 4, 0.5, "2.00", false},
		{"Modulo positive numbers", "modulo", 10, 3, "1.00", false},
		{"Modulo negative dividend", "modulo", -10, 3, "-1.00", false},
		{"Modulo by zero", "modulo", 10, 0, "", true},
		{"Modulo exact division", "modulo", 15, 5, "0.00", false},
		{"Modulo fractional numbers", "modulo", 7.5, 2.5, "0.00", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{
						"operation": tt.operation,
						"x":         tt.x,
						"y":         tt.y,
					},
				},
			}

			result, err := calculator.Handler(ctx, request)
			if err != nil {
				t.Fatalf("Handler returned error: %v", err)
			}

			if tt.hasError {
				if !result.IsError {
					t.Errorf("Expected error but got success")
				}
			} else {
				if result.IsError {
					t.Errorf("Expected success but got error: %v", result.Content)
				}
				if len(result.Content) > 0 {
					if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
						if textContent.Text != tt.expected {
							t.Errorf("Expected %s, got %s", tt.expected, textContent.Text)
						}
					}
				}
			}
		})
	}
}

func TestCalculatorTool_SqrtOperation(t *testing.T) {
	calculator := NewCalculatorTool()
	ctx := context.Background()

	tests := []struct {
		name     string
		x        float64
		expected string
		hasError bool
	}{
		{"Square root of positive number", 9, "3.00", false},
		{"Square root of zero", 0, "0.00", false},
		{"Square root of negative number", -4, "", true},
		{"Square root of fractional number", 2.25, "1.50", false},
		{"Square root of large number", 100, "10.00", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{
						"operation": "sqrt",
						"x":         tt.x,
					},
				},
			}

			result, err := calculator.Handler(ctx, request)
			if err != nil {
				t.Fatalf("Handler returned error: %v", err)
			}

			if tt.hasError {
				if !result.IsError {
					t.Errorf("Expected error but got success")
				}
			} else {
				if result.IsError {
					t.Errorf("Expected success but got error: %v", result.Content)
				}
				if len(result.Content) > 0 {
					if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
						if textContent.Text != tt.expected {
							t.Errorf("Expected %s, got %s", tt.expected, textContent.Text)
						}
					}
				}
			}
		})
	}
}

func TestCalculatorTool_InvalidOperations(t *testing.T) {
	calculator := NewCalculatorTool()
	ctx := context.Background()

	tests := []struct {
		name      string
		operation string
		x         float64
		y         float64
	}{
		{"Invalid operation", "invalid", 5, 3},
		{"Empty operation", "", 5, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{
						"operation": tt.operation,
						"x":         tt.x,
						"y":         tt.y,
					},
				},
			}

			result, err := calculator.Handler(ctx, request)
			if err != nil {
				t.Fatalf("Handler returned error: %v", err)
			}

			if !result.IsError {
				t.Errorf("Expected error for invalid operation %s", tt.operation)
			}
		})
	}
}

func TestCalculatorTool_MissingParameters(t *testing.T) {
	calculator := NewCalculatorTool()
	ctx := context.Background()

	tests := []struct {
		name      string
		arguments map[string]interface{}
	}{
		{"Missing operation", map[string]interface{}{"x": 5, "y": 3}},
		{"Missing x parameter", map[string]interface{}{"operation": "add", "y": 3}},
		{"Missing y parameter for add", map[string]interface{}{"operation": "add", "x": 5}},
		{"Missing y parameter for power", map[string]interface{}{"operation": "power", "x": 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: tt.arguments,
				},
			}

			result, err := calculator.Handler(ctx, request)
			if err != nil {
				t.Fatalf("Handler returned error: %v", err)
			}

			if !result.IsError {
				t.Errorf("Expected error for missing parameters")
			}
		})
	}
}

func TestCalculatorTool_TrigonometricOperations(t *testing.T) {
	calculator := NewCalculatorTool()
	ctx := context.Background()

	tests := []struct {
		name      string
		operation string
		x         float64
		angleUnit string
		expected  string
		hasError  bool
	}{
		{"Sin 0 radians", "sin", 0, "radians", "0.00", false},
		{"Sin 90 degrees", "sin", 90, "degrees", "1.00", false},
		{"Sin π/2 radians", "sin", 1.5708, "radians", "1.00", false},
		{"Cos 0 radians", "cos", 0, "radians", "1.00", false},
		{"Cos 90 degrees", "cos", 90, "degrees", "0.00", false},
		{"Cos π radians", "cos", 3.1416, "radians", "-1.00", false},
		{"Tan 0 radians", "tan", 0, "radians", "0.00", false},
		{"Tan 45 degrees", "tan", 45, "degrees", "1.00", false},
		{"Tan π/4 radians", "tan", 0.7854, "radians", "1.00", false},
		{"Sin default unit (radians)", "sin", 0, "", "0.00", false},
		{"Cos default unit (radians)", "cos", 0, "", "1.00", false},
		{"Tan default unit (radians)", "tan", 0, "", "0.00", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := map[string]interface{}{
				"operation": tt.operation,
				"x":         tt.x,
			}
			if tt.angleUnit != "" {
				args["angle_unit"] = tt.angleUnit
			}

			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: args,
				},
			}

			result, err := calculator.Handler(ctx, request)
			if err != nil {
				t.Fatalf("Handler returned error: %v", err)
			}

			if tt.hasError {
				if !result.IsError {
					t.Errorf("Expected error but got success")
				}
			} else {
				if result.IsError {
					t.Errorf("Expected success but got error: %v", result.Content)
				}
				if len(result.Content) > 0 {
					if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
						if textContent.Text != tt.expected {
							t.Errorf("Expected %s, got %s", tt.expected, textContent.Text)
						}
					}
				}
			}
		})
	}
}

func TestCalculatorTool_InverseTrigonometricOperations(t *testing.T) {
	calculator := NewCalculatorTool()
	ctx := context.Background()

	tests := []struct {
		name      string
		operation string
		x         float64
		angleUnit string
		expected  string
		hasError  bool
	}{
		{"Asin 0 radians", "asin", 0, "radians", "0.00", false},
		{"Asin 1 degrees", "asin", 1, "degrees", "90.00", false},
		{"Asin 0.5 radians", "asin", 0.5, "radians", "0.52", false},
		{"Asin invalid input", "asin", 2, "radians", "", true},
		{"Asin invalid input negative", "asin", -2, "radians", "", true},
		{"Acos 1 radians", "acos", 1, "radians", "0.00", false},
		{"Acos 0 degrees", "acos", 0, "degrees", "90.00", false},
		{"Acos -1 radians", "acos", -1, "radians", "3.14", false},
		{"Acos invalid input", "acos", 2, "radians", "", true},
		{"Acos invalid input negative", "acos", -2, "radians", "", true},
		{"Atan 0 radians", "atan", 0, "radians", "0.00", false},
		{"Atan 1 degrees", "atan", 1, "degrees", "45.00", false},
		{"Atan -1 radians", "atan", -1, "radians", "-0.79", false},
		{"Atan default unit (radians)", "atan", 0, "", "0.00", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := map[string]interface{}{
				"operation": tt.operation,
				"x":         tt.x,
			}
			if tt.angleUnit != "" {
				args["angle_unit"] = tt.angleUnit
			}

			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: args,
				},
			}

			result, err := calculator.Handler(ctx, request)
			if err != nil {
				t.Fatalf("Handler returned error: %v", err)
			}

			if tt.hasError {
				if !result.IsError {
					t.Errorf("Expected error but got success")
				}
			} else {
				if result.IsError {
					t.Errorf("Expected success but got error: %v", result.Content)
				}
				if len(result.Content) > 0 {
					if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
						if textContent.Text != tt.expected {
							t.Errorf("Expected %s, got %s", tt.expected, textContent.Text)
						}
					}
				}
			}
		})
	}
}
