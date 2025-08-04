package tools

import (
	"context"
	"fmt"
	"math"

	"github.com/expr-lang/expr"
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
		mcp.WithDescription("Evaluate mathematical expressions using natural syntax (e.g., '2 + 3 * 4', 'sin(pi/4)', 'sqrt(16)')"),
		mcp.WithString("expression",
			mcp.Required(),
			mcp.Description("Mathematical expression to evaluate. Supports +, -, *, /, ^, sqrt(), sin(), cos(), tan(), asin(), acos(), atan(), log(), ln(), abs(), ceil(), floor(), round(), pi, e"),
		),
	)
}

func calculatorToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	expression, err := request.RequireString("expression")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	env := map[string]interface{}{
		"pi":    math.Pi,
		"e":     math.E,
		"sqrt":  math.Sqrt,
		"sin":   math.Sin,
		"cos":   math.Cos,
		"tan":   math.Tan,
		"asin":  math.Asin,
		"acos":  math.Acos,
		"atan":  math.Atan,
		"log":   math.Log10,
		"ln":    math.Log,
		"abs":   math.Abs,
		"ceil":  math.Ceil,
		"floor": math.Floor,
		"round": math.Round,
		"pow":   math.Pow,
	}

	program, err := expr.Compile(expression, expr.Env(env))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Expression compilation error: %v", err)), nil
	}

	result, err := expr.Run(program, env)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Expression evaluation error: %v", err)), nil
	}

	switch v := result.(type) {
	case float64:
		return mcp.NewToolResultText(fmt.Sprintf("%.6g", v)), nil
	case int64:
		return mcp.NewToolResultText(fmt.Sprintf("%d", v)), nil
	case int:
		return mcp.NewToolResultText(fmt.Sprintf("%d", v)), nil
	default:
		return mcp.NewToolResultText(fmt.Sprintf("%v", v)), nil
	}
}
