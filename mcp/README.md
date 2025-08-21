# MCP HTTP Proxy Server

A generic MCP (Model Context Protocol) server that provides HTTP proxy functionality to other hosted MCP servers.

## Usage

```bash
./mcp-proxy -name <server-name>
```

The proxy server uses the `name` parameter to lookup an environment variable `${name}_HOST` that specifies the target MCP server to proxy to.

## Example

```bash
# Set the target MCP server
export myserver_HOST="http://localhost:3000"

# Start the proxy
./mcp-proxy -name myserver
```

## Features

- **Capability-Aware Proxying**: Only exposes and registers capabilities that the origin server actually supports
- **Dynamic Discovery**: Automatically discovers and registers tools, resources, and prompts from the target server
- **Full MCP Compatibility**: Supports all standard MCP operations including:
  - Tools (list and call)
  - Resources (list and read) 
  - Prompts (list and get)
- **HTTP Proxy**: Transparently proxies requests to target MCP servers over HTTP
- **Error Handling**: Comprehensive error handling with detailed logging

## Architecture

The proxy server:
1. Takes a command line argument `name` 
2. Looks up environment variable `${name}_HOST` for the target server URL
3. Initializes connection to the target MCP server and discovers its capabilities
4. Creates a proxy server with only the capabilities that the origin server supports
5. Discovers and registers only the features (tools, resources, prompts) that the origin server provides
6. Registers handlers that proxy requests to the target server
7. Runs as a standard MCP server over stdio

This ensures that clients connecting to the proxy only see the capabilities and features that are actually available from the origin server, following the MCP specification for capability negotiation.

## Dependencies

Built using the [mcp-go](https://github.com/mark3labs/mcp-go) library for MCP protocol implementation.