package tools

import (
	"context"
	"fmt"
	"math"

	"github.com/mark3labs/mcp-go/mcp"
)

type CalculatorTool struct {
	Tool    mcp.Tool
	Handler func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

func NewCalculatorTool() *CalculatorTool {
	return &CalculatorTool{
		Tool:    calculatorTool(),
		Handler: calculatorToolHandler,
	}
}

func calculatorTool() mcp.Tool {
	return mcp.NewTool("calculate",
		mcp.WithDescription("Perform basic and some advanced arithmetic operations"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform (add, subtract, multiply, divide, power, sqrt, modulo)"),
			mcp.Enum("add", "subtract", "multiply", "divide", "power", "sqrt", "modulo"),
		),
		mcp.WithNumber("x",
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber("y",
			mcp.Description("Second number (not required for sqrt operation)"),
		),
	)
}

func calculatorToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	op, err := request.RequireString("operation")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	x, err := request.RequireFloat("x")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var result float64
	switch op {
	case "add":
		y, err := request.RequireFloat("y")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result = x + y
	case "subtract":
		y, err := request.RequireFloat("y")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result = x - y
	case "multiply":
		y, err := request.RequireFloat("y")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result = x * y
	case "divide":
		y, err := request.RequireFloat("y")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if y == 0 {
			return mcp.NewToolResultError("Cannot divide by zero"), nil
		}
		result = x / y
	case "power":
		y, err := request.RequireFloat("y")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result = math.Pow(x, y)
	case "sqrt":
		if x < 0 {
			return mcp.NewToolResultError("Cannot calculate square root of negative number"), nil
		}
		result = math.Sqrt(x)
	case "modulo":
		y, err := request.RequireFloat("y")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if y == 0 {
			return mcp.NewToolResultError("Cannot perform modulo with zero"), nil
		}
		result = math.Mod(x, y)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Unknown operation: %s", op)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("%.2f", result)), nil
}
