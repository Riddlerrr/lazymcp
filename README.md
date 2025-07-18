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
Perform basic arithmetic operations.

**Parameters:**
- `operation` (required): The operation to perform (add, subtract, multiply, divide)
- `x` (required): First number
- `y` (required): Second number

**Example:**
```json
{
  "name": "calculate",
  "arguments": {
    "operation": "add",
    "x": 5,
    "y": 3
  }
}
```

**Returns:** The result of the arithmetic operation as a formatted number.

## Development

To run the server in development mode:
```bash
go run main.go
```

## License

MIT
