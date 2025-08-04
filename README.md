# LazyMCP - MCP Time Server

A general-purpose MCP (Model Context Protocol) server written in Go that provides real-time data, starting with basic calculator functionality.

## Installation

1. Make sure you have Go installed (version 1.24 or later)
2. Clone this repository:
   ```bash
   git clone https://github.com/Riddlerrr/lazymcp.git
   cd lazymcp
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```

## Building

Build the server:
```bash
bin/build
```

## Usage

### Running the Server

The server uses Streamable HTTP transport for communication:

```bash
bin/run
```

### Available Tools

#### `calculate`
Evaluate mathematical expressions using natural syntax.

**Parameters:**
- `expression` (required): Mathematical expression to evaluate. Supports +, -, *, /, ^, sqrt(), sin(), cos(), tan(), asin(), acos(), atan(), log(), ln(), abs(), ceil(), floor(), round(), pi, e

**Examples:**
```json
{
  "name": "calculate",
  "arguments": {
    "expression": "2 + 3 * 4"
  }
}
```

```json
{
  "name": "calculate",
  "arguments": {
    "expression": "sin(pi/4)"
  }
}
```

```json
{
  "name": "calculate",
  "arguments": {
    "expression": "sqrt(16)"
  }
}
```

```json
{
  "name": "calculate",
  "arguments": {
    "expression": "pow(2, 8)"
  }
}
```

**Returns:** The result of the expression evaluation as a formatted number.

## Development

To run the server in development mode:
```bash
go run main.go
```

## License

MIT
