# LazyMCP - MCP Time Server

A general-purpose MCP (Model Context Protocol) server written in Go that provides real-time data, starting with current date and time functionality.

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

The server uses stdio transport for communication:

```bash
bin/run
```

### Available Tools

#### `get_current_time`

Get the current date and time in various formats.

Parameters:
- `format` (optional): Time format. Options include:
  - `unix` - Unix timestamp
  - `iso8601` - ISO 8601 format (default)
  - `rfc3339` - RFC 3339 format
  - `kitchen` - Kitchen format (e.g., "3:04PM")
  - Custom Go time format string
- `timezone` (optional): Timezone name (default: "UTC")
  - Examples: "America/New_York", "Europe/London", "Asia/Tokyo"

Example usage with Claude Desktop or other MCP clients:

```json
{
  "tool": "get_current_time",
  "arguments": {
    "format": "iso8601",
    "timezone": "America/New_York"
  }
}
```

### Integration with Claude Desktop

To use this server with Claude Desktop, add it to your configuration:

1. Build the server first
2. Add to your Claude Desktop configuration (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "lazymcp": {
      "command": "/path/to/lazymcp"
    }
  }
}
```

## Development

To run the server in development mode:
```bash
go run main.go
```

## License

MIT
