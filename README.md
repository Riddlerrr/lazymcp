# LazyMCP - An MCP server that can really help

A general-purpose MCP (Model Context Protocol) server written in Go that provides calculator, IP lookup, weather, and weather forecast functionality with real-time data access.

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

## Configuration

### Weather Tools Setup
The weather and weather forecast tools require an OpenWeatherMap API key to function:

1. Get a free API key from https://openweathermap.org/api
2. Copy `.env.example` to `.env`
3. Set your API key: `OPENWEATHER_API_KEY=your_api_key_here`

## Usage

### Running the Server

The server uses Streamable HTTP transport for communication on `http://localhost:3000/mcp`:

```bash
bin/run
```

### Available Tools

#### `calculate`
Evaluate mathematical expressions using natural syntax.

**Parameters:**
- `expression` (required): Mathematical expression to evaluate. Supports +, -, *, /, ^, sqrt(), sin(), cos(), tan(), asin(), acos(), atan(), log(), ln(), abs(), ceil(), floor(), round(), pow(), pi, e

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

#### `get_ip`
Get the IP address of the client making the request.

**Parameters:** None

**Example:**
```json
{
  "name": "get_ip",
  "arguments": {}
}
```

**Returns:** The client's IP address as a string.

#### `get_ip_data`
Get detailed information about an IP address including geolocation data.

**Parameters:**
- `ip` (optional): IP address to lookup. If not provided, uses the client's IP address.

**Examples:**
```json
{
  "name": "get_ip_data",
  "arguments": {}
}
```

```json
{
  "name": "get_ip_data",
  "arguments": {
    "ip": "8.8.8.8"
  }
}
```

**Returns:** Formatted markdown with location information including country, region, city, coordinates, timezone, ISP, and network details.

#### `get_weather`
Get current weather for a location. Uses client's IP location by default, or accepts a custom location parameter.

**Parameters:**
- `location` (optional): Location to get weather for. Can be city name (e.g., 'London' or 'New York,US') or coordinates (e.g., '40.7128,-74.0060'). Uses client IP location if not provided.

**Examples:**
```json
{
  "name": "get_weather",
  "arguments": {}
}
```

```json
{
  "name": "get_weather",
  "arguments": {
    "location": "London"
  }
}
```

```json
{
  "name": "get_weather",
  "arguments": {
    "location": "40.7128,-74.0060"
  }
}
```

**Returns:** Formatted markdown with current weather conditions, temperature, humidity, pressure, wind, visibility, and location details. Units are automatically determined (metric for most countries, imperial for US locations).

#### `get_weather_forecast`
Get 5-day weather forecast for a location. Uses client's IP location by default, or accepts a custom location parameter.

**Parameters:**
- `location` (optional): Location to get forecast for. Can be city name (e.g., 'London' or 'New York,US') or coordinates (e.g., '40.7128,-74.0060'). Uses client IP location if not provided.

**Examples:**
```json
{
  "name": "get_weather_forecast",
  "arguments": {}
}
```

```json
{
  "name": "get_weather_forecast",
  "arguments": {
    "location": "Valencia,ES"
  }
}
```

```json
{
  "name": "get_weather_forecast",
  "arguments": {
    "location": "39.4676,-0.3771"
  }
}
```

**Returns:** Formatted markdown with:
- **Next 24 Hours**: Weather forecast in 3-hour intervals with temperature, conditions, and precipitation probability
- **5-Day Forecast**: Daily summaries with temperature ranges, weather conditions, and precipitation chances
- **Location Details**: City, country, and coordinates
- **Automatic Units**: Metric for most countries, imperial for US locations

## Development

To run the server in development mode:
```bash
go run main.go
```

## License

MIT
