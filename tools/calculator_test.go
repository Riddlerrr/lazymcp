package tools

import (
	"context"
	"strconv"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestCalculatorTool_BasicOperations(t *testing.T) {
	calculator := NewCalculatorTool()
	ctx := context.Background()

	tests := []struct {
		name       string
		expression string
		expected   string
		hasError   bool
	}{
		{"Add positive numbers", "5 + 3", "8", false},
		{"Add negative numbers", "-5 + (-3)", "-8", false},
		{"Add mixed numbers", "5 + (-3)", "2", false},
		{"Subtract positive numbers", "10 - 4", "6", false},
		{"Subtract negative numbers", "-10 - (-4)", "-6", false},
		{"Multiply positive numbers", "6 * 7", "42", false},
		{"Multiply by zero", "5 * 0", "0", false},
		{"Multiply negative numbers", "(-3) * (-4)", "12", false},
		{"Divide positive numbers", "15 / 3", "5", false},
		{"Divide by zero", "10 / 0", "+Inf", false},
		{"Divide negative numbers", "(-12) / (-3)", "4", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{
						"expression": tt.expression,
					},
				},
			}

			result, err := calculator.Handler(ctx, request)
			if err != nil {
				t.Fatalf("Handler returned error: %v", err)
			}

			if tt.hasError {
				if !result.IsError {
					t.Errorf("Expected error but got success: %v", result)
				}
			} else {
				if result.IsError {
					t.Errorf("Expected success but got error: %v", result)
				} else {
					if len(result.Content) == 0 {
						t.Errorf("No content in result")
						return
					}
					if textContent, ok := mcp.AsTextContent(result.Content[0]); ok {
						resultFloat, err := strconv.ParseFloat(textContent.Text, 64)
						if err != nil {
							t.Errorf("Could not parse result as float: %s", textContent.Text)
						}
						expectedFloat, err := strconv.ParseFloat(tt.expected, 64)
						if err != nil {
							t.Errorf("Could not parse expected as float: %s", tt.expected)
						}
						if resultFloat != expectedFloat {
							t.Errorf("Expected %s, got %s", tt.expected, textContent.Text)
						}
					} else {
						t.Errorf("Result does not contain text content")
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
		name       string
		expression string
		expected   float64
		tolerance  float64
		hasError   bool
	}{
		{"Power operation", "pow(2, 3)", 8.0, 0.001, false},
		{"Square root", "sqrt(16)", 4.0, 0.001, false},
		{"Square root of zero", "sqrt(0)", 0.0, 0.001, false},
		{"Modulo operation", "7 % 3", 1.0, 0.001, false},
		{"Pi constant", "pi", 3.14159, 0.001, false},
		{"E constant", "e", 2.71828, 0.001, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{
						"expression": tt.expression,
					},
				},
			}

			result, err := calculator.Handler(ctx, request)
			if err != nil {
				t.Fatalf("Handler returned error: %v", err)
			}

			if tt.hasError {
				if !result.IsError {
					t.Errorf("Expected error but got success: %v", result)
				}
			} else {
				if result.IsError {
					t.Errorf("Expected success but got error: %v", result)
				} else {
					if len(result.Content) == 0 {
						t.Errorf("No content in result")
						return
					}
					if textContent, ok := mcp.AsTextContent(result.Content[0]); ok {
						resultFloat, err := strconv.ParseFloat(textContent.Text, 64)
						if err != nil {
							t.Errorf("Could not parse result as float: %s", textContent.Text)
						}
						if abs(resultFloat-tt.expected) > tt.tolerance {
							t.Errorf("Expected %f (±%f), got %f", tt.expected, tt.tolerance, resultFloat)
						}
					} else {
						t.Errorf("Result does not contain text content")
					}
				}
			}
		})
	}
}

func TestCalculatorTool_TrigonometricOperations(t *testing.T) {
	calculator := NewCalculatorTool()
	ctx := context.Background()

	tests := []struct {
		name       string
		expression string
		expected   float64
		tolerance  float64
		hasError   bool
	}{
		{"Sin of 0", "sin(0)", 0.0, 0.001, false},
		{"Sin of pi/2", "sin(pi/2)", 1.0, 0.001, false},
		{"Cos of 0", "cos(0)", 1.0, 0.001, false},
		{"Cos of pi", "cos(pi)", -1.0, 0.001, false},
		{"Tan of 0", "tan(0)", 0.0, 0.001, false},
		{"Tan of pi/4", "tan(pi/4)", 1.0, 0.001, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{
						"expression": tt.expression,
					},
				},
			}

			result, err := calculator.Handler(ctx, request)
			if err != nil {
				t.Fatalf("Handler returned error: %v", err)
			}

			if tt.hasError {
				if !result.IsError {
					t.Errorf("Expected error but got success: %v", result)
				}
			} else {
				if result.IsError {
					t.Errorf("Expected success but got error: %v", result)
				} else {
					if len(result.Content) == 0 {
						t.Errorf("No content in result")
						return
					}
					if textContent, ok := mcp.AsTextContent(result.Content[0]); ok {
						resultFloat, err := strconv.ParseFloat(textContent.Text, 64)
						if err != nil {
							t.Errorf("Could not parse result as float: %s", textContent.Text)
						}
						if abs(resultFloat-tt.expected) > tt.tolerance {
							t.Errorf("Expected %f (±%f), got %f", tt.expected, tt.tolerance, resultFloat)
						}
					} else {
						t.Errorf("Result does not contain text content")
					}
				}
			}
		})
	}
}

func TestCalculatorTool_ComplexExpressions(t *testing.T) {
	calculator := NewCalculatorTool()
	ctx := context.Background()

	tests := []struct {
		name       string
		expression string
		expected   float64
		tolerance  float64
		hasError   bool
	}{
		{"Order of operations", "2 + 3 * 4", 14.0, 0.001, false},
		{"Parentheses", "(2 + 3) * 4", 20.0, 0.001, false},
		{"Mixed operations", "sqrt(16) + pow(2, 3)", 12.0, 0.001, false},
		{"Nested functions", "sin(asin(0.5))", 0.5, 0.001, false},
		{"Expression with constants", "pi * 2", 6.28318, 0.001, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{
						"expression": tt.expression,
					},
				},
			}

			result, err := calculator.Handler(ctx, request)
			if err != nil {
				t.Fatalf("Handler returned error: %v", err)
			}

			if tt.hasError {
				if !result.IsError {
					t.Errorf("Expected error but got success: %v", result)
				}
			} else {
				if result.IsError {
					t.Errorf("Expected success but got error: %v", result)
				} else {
					if len(result.Content) == 0 {
						t.Errorf("No content in result")
						return
					}
					if textContent, ok := mcp.AsTextContent(result.Content[0]); ok {
						resultFloat, err := strconv.ParseFloat(textContent.Text, 64)
						if err != nil {
							t.Errorf("Could not parse result as float: %s", textContent.Text)
						}
						if abs(resultFloat-tt.expected) > tt.tolerance {
							t.Errorf("Expected %f (±%f), got %f", tt.expected, tt.tolerance, resultFloat)
						}
					} else {
						t.Errorf("Result does not contain text content")
					}
				}
			}
		})
	}
}

func TestCalculatorTool_ErrorCases(t *testing.T) {
	calculator := NewCalculatorTool()
	ctx := context.Background()

	tests := []struct {
		name       string
		expression string
		hasError   bool
	}{
		{"Invalid expression", "2 +", true},
		{"Unknown function", "foo(2)", true},
		{"Missing closing parenthesis", "2 * (3 + 4", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{
						"expression": tt.expression,
					},
				},
			}

			result, err := calculator.Handler(ctx, request)
			if err != nil {
				t.Fatalf("Handler returned error: %v", err)
			}

			if tt.hasError {
				if !result.IsError {
					t.Errorf("Expected error but got success: %v", result)
				}
			} else {
				if result.IsError {
					t.Errorf("Expected success but got error: %v", result)
				}
			}
		})
	}
}

// Helper function for absolute value
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}