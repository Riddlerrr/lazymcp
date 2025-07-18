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
		mcp.WithDescription("Perform basic, advanced arithmetic, and trigonometric operations"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform (add, subtract, multiply, divide, power, sqrt, modulo, sin, cos, tan, asin, acos, atan)"),
			mcp.Enum("add", "subtract", "multiply", "divide", "power", "sqrt", "modulo", "sin", "cos", "tan", "asin", "acos", "atan"),
		),
		mcp.WithNumber("x",
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber("y",
			mcp.Description("Second number (not required for sqrt, sin, cos, tan, asin, acos, atan operations)"),
		),
		mcp.WithString("angle_unit",
			mcp.Description("Unit for trigonometric operations (degrees or radians, defaults to radians)"),
			mcp.Enum("degrees", "radians"),
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
	case "sin":
		angle := x
		angleUnit := request.GetString("angle_unit", "radians")
		if angleUnit == "degrees" {
			angle = x * math.Pi / 180
		}
		result = math.Sin(angle)
	case "cos":
		angle := x
		angleUnit := request.GetString("angle_unit", "radians")
		if angleUnit == "degrees" {
			angle = x * math.Pi / 180
		}
		result = math.Cos(angle)
	case "tan":
		angle := x
		angleUnit := request.GetString("angle_unit", "radians")
		if angleUnit == "degrees" {
			angle = x * math.Pi / 180
		}
		result = math.Tan(angle)
	case "asin":
		if x < -1 || x > 1 {
			return mcp.NewToolResultError("Input for asin must be between -1 and 1"), nil
		}
		result = math.Asin(x)
		angleUnit := request.GetString("angle_unit", "radians")
		if angleUnit == "degrees" {
			result = result * 180 / math.Pi
		}
	case "acos":
		if x < -1 || x > 1 {
			return mcp.NewToolResultError("Input for acos must be between -1 and 1"), nil
		}
		result = math.Acos(x)
		angleUnit := request.GetString("angle_unit", "radians")
		if angleUnit == "degrees" {
			result = result * 180 / math.Pi
		}
	case "atan":
		result = math.Atan(x)
		angleUnit := request.GetString("angle_unit", "radians")
		if angleUnit == "degrees" {
			result = result * 180 / math.Pi
		}
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Unknown operation: %s", op)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("%.2f", result)), nil
}
