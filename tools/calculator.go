package tools

import (
	"context"
	"fmt"

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
		mcp.WithDescription("Perform basic arithmetic operations"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform (add, subtract, multiply, divide)"),
			mcp.Enum("add", "subtract", "multiply", "divide"),
		),
		mcp.WithNumber("x",
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber("y",
			mcp.Required(),
			mcp.Description("Second number"),
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

	y, err := request.RequireFloat("y")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var result float64
	switch op {
	case "add":
		result = x + y
	case "subtract":
		result = x - y
	case "multiply":
		result = x * y
	case "divide":
		if y == 0 {
			return mcp.NewToolResultError("Cannot divide by zero"), nil
		}
		result = x / y
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Unknown operation: %s", op)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("%.2f", result)), nil
}
