package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	var name string
	flag.StringVar(&name, "name", "", "Name of the MCP server to proxy to (required)")
	flag.Parse()

	if name == "" {
		log.Fatal("Error: -name argument is required")
	}

	// Look up the host from environment variable
	hostEnvVar := strings.ToUpper(name) + "_HOST"
	targetHost := os.Getenv(hostEnvVar)
	if targetHost == "" {
		log.Fatalf("Error: environment variable %s is not set", hostEnvVar)
	}

	log.Printf("Starting MCP proxy server for %s, proxying to %s", name, targetHost)

	// Create HTTP proxy client
	proxyClient := &HTTPProxyClient{targetHost: targetHost}

	// Initialize connection to target server and discover capabilities
	if err := proxyClient.Initialize(context.Background()); err != nil {
		log.Fatalf("Failed to initialize proxy client: %v", err)
	}

	// Create the MCP server with capabilities matching the origin server
	mcpServer := proxyClient.CreateMCPServerWithCapabilities()

	// Only discover and register features that the origin server supports
	if proxyClient.capabilities.Tools != nil {
		if err := proxyClient.RegisterToolsOnServer(context.Background(), mcpServer); err != nil {
			log.Printf("Warning: Failed to register tools: %v", err)
		}
	} else {
		log.Printf("Origin server does not support tools - skipping tool discovery")
	}

	if proxyClient.capabilities.Resources != nil {
		if err := proxyClient.RegisterResourcesOnServer(context.Background(), mcpServer); err != nil {
			log.Printf("Warning: Failed to register resources: %v", err)
		}
	} else {
		log.Printf("Origin server does not support resources - skipping resource discovery")
	}

	if proxyClient.capabilities.Prompts != nil {
		if err := proxyClient.RegisterPromptsOnServer(context.Background(), mcpServer); err != nil {
			log.Printf("Warning: Failed to register prompts: %v", err)
		}
	} else {
		log.Printf("Origin server does not support prompts - skipping prompt discovery")
	}

	// Create and run the stdio server
	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// HTTPProxyClient handles HTTP requests to the target MCP server
type HTTPProxyClient struct {
	targetHost   string
	client       *http.Client
	capabilities mcp.ServerCapabilities
}

// Initialize sets up the HTTP client and tests connectivity to target server
func (h *HTTPProxyClient) Initialize(ctx context.Context) error {
	h.client = &http.Client{}
	log.Printf("Initializing proxy connection to %s", h.targetHost)

	// Test connection with an initialize request
	initParams := mcp.InitializeParams{
		ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
		Capabilities: mcp.ClientCapabilities{
			Roots: &struct {
				ListChanged bool `json:"listChanged,omitempty"`
			}{ListChanged: false},
		},
		ClientInfo: mcp.Implementation{
			Name:    "mcp-proxy",
			Version: "1.0.0",
		},
	}

	result, err := h.proxyRequest(ctx, "initialize", initParams)
	if err != nil {
		return fmt.Errorf("failed to initialize target server: %w", err)
	}

	// Parse the initialize result to capture server capabilities
	var initResult mcp.InitializeResult
	if err := json.Unmarshal(result, &initResult); err != nil {
		return fmt.Errorf("failed to unmarshal initialize result: %w", err)
	}

	// Store the server capabilities for later use
	h.capabilities = initResult.Capabilities
	log.Printf("Discovered server capabilities: tools=%v, resources=%v, prompts=%v", 
		h.capabilities.Tools != nil, h.capabilities.Resources != nil, h.capabilities.Prompts != nil)

	log.Printf("Successfully initialized proxy to %s", h.targetHost)
	return nil
}

// RegisterToolsOnServer discovers tools from target server and registers them
func (h *HTTPProxyClient) RegisterToolsOnServer(ctx context.Context, mcpServer *server.MCPServer) error {
	log.Printf("Discovering tools from %s", h.targetHost)

	result, err := h.proxyRequest(ctx, "tools/list", mcp.PaginatedParams{})
	if err != nil {
		return fmt.Errorf("failed to list tools: %w", err)
	}

	var listResult mcp.ListToolsResult
	if err := json.Unmarshal(result, &listResult); err != nil {
		return fmt.Errorf("failed to unmarshal tools list: %w", err)
	}

	for _, tool := range listResult.Tools {
		toolName := tool.Name
		mcpServer.AddTool(tool, h.createToolHandler(toolName))
		log.Printf("Registered tool: %s", toolName)
	}

	return nil
}

// RegisterResourcesOnServer discovers resources from target server and registers them
func (h *HTTPProxyClient) RegisterResourcesOnServer(ctx context.Context, mcpServer *server.MCPServer) error {
	log.Printf("Discovering resources from %s", h.targetHost)

	result, err := h.proxyRequest(ctx, "resources/list", mcp.PaginatedParams{})
	if err != nil {
		return fmt.Errorf("failed to list resources: %w", err)
	}

	var listResult mcp.ListResourcesResult
	if err := json.Unmarshal(result, &listResult); err != nil {
		return fmt.Errorf("failed to unmarshal resources list: %w", err)
	}

	for _, resource := range listResult.Resources {
		resourceURI := resource.URI
		mcpServer.AddResource(resource, h.createResourceHandler(resourceURI))
		log.Printf("Registered resource: %s", resourceURI)
	}

	return nil
}

// RegisterPromptsOnServer discovers prompts from target server and registers them
func (h *HTTPProxyClient) RegisterPromptsOnServer(ctx context.Context, mcpServer *server.MCPServer) error {
	log.Printf("Discovering prompts from %s", h.targetHost)

	result, err := h.proxyRequest(ctx, "prompts/list", mcp.PaginatedParams{})
	if err != nil {
		return fmt.Errorf("failed to list prompts: %w", err)
	}

	var listResult mcp.ListPromptsResult
	if err := json.Unmarshal(result, &listResult); err != nil {
		return fmt.Errorf("failed to unmarshal prompts list: %w", err)
	}

	for _, prompt := range listResult.Prompts {
		promptName := prompt.Name
		mcpServer.AddPrompt(prompt, h.createPromptHandler(promptName))
		log.Printf("Registered prompt: %s", promptName)
	}

	return nil
}

// createToolHandler creates a handler that proxies tool calls
func (h *HTTPProxyClient) createToolHandler(toolName string) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("Proxying call_tool request for tool '%s' to %s", toolName, h.targetHost)

		result, err := h.proxyRequest(ctx, "tools/call", request.Params)
		if err != nil {
			return nil, err
		}

		var callResult mcp.CallToolResult
		if err := json.Unmarshal(result, &callResult); err != nil {
			return nil, fmt.Errorf("failed to unmarshal call_tool result: %w", err)
		}

		return &callResult, nil
	}
}

// createResourceHandler creates a handler that proxies resource reads
func (h *HTTPProxyClient) createResourceHandler(resourceURI string) server.ResourceHandlerFunc {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		log.Printf("Proxying read_resource request for URI '%s' to %s", resourceURI, h.targetHost)

		result, err := h.proxyRequest(ctx, "resources/read", request.Params)
		if err != nil {
			return nil, err
		}

		var readResult mcp.ReadResourceResult
		if err := json.Unmarshal(result, &readResult); err != nil {
			return nil, fmt.Errorf("failed to unmarshal read_resource result: %w", err)
		}

		return readResult.Contents, nil
	}
}

// createPromptHandler creates a handler that proxies prompt requests
func (h *HTTPProxyClient) createPromptHandler(promptName string) server.PromptHandlerFunc {
	return func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		log.Printf("Proxying get_prompt request for prompt '%s' to %s", promptName, h.targetHost)

		result, err := h.proxyRequest(ctx, "prompts/get", request.Params)
		if err != nil {
			return nil, err
		}

		var getResult mcp.GetPromptResult
		if err := json.Unmarshal(result, &getResult); err != nil {
			return nil, fmt.Errorf("failed to unmarshal get_prompt result: %w", err)
		}

		return &getResult, nil
	}
}

// proxyRequest makes an HTTP request to the target MCP server
func (h *HTTPProxyClient) proxyRequest(ctx context.Context, method string, params any) ([]byte, error) {
	if h.client == nil {
		return nil, fmt.Errorf("HTTP client not initialized")
	}

	// Create JSON-RPC request
	requestBody := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  method,
		"params":  params,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", h.targetHost, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse JSON-RPC response
	var jsonRPCResp struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      any             `json:"id"`
		Result  json.RawMessage `json:"result,omitempty"`
		Error   *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Data    any    `json:"data,omitempty"`
		} `json:"error,omitempty"`
	}

	if err := json.Unmarshal(body, &jsonRPCResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON-RPC response: %w", err)
	}

	if jsonRPCResp.Error != nil {
		return nil, fmt.Errorf("JSON-RPC error %d: %s", jsonRPCResp.Error.Code, jsonRPCResp.Error.Message)
	}

	return jsonRPCResp.Result, nil
}

// CreateMCPServerWithCapabilities creates an MCP server with capabilities matching the origin server
func (h *HTTPProxyClient) CreateMCPServerWithCapabilities() *server.MCPServer {
	var options []server.ServerOption

	// Add capabilities based on what the origin server supports
	if h.capabilities.Tools != nil {
		options = append(options, server.WithToolCapabilities(h.capabilities.Tools.ListChanged))
	}

	if h.capabilities.Resources != nil {
		options = append(options, server.WithResourceCapabilities(
			h.capabilities.Resources.Subscribe, h.capabilities.Resources.ListChanged))
	}

	if h.capabilities.Prompts != nil {
		options = append(options, server.WithPromptCapabilities(h.capabilities.Prompts.ListChanged))
	}

	if h.capabilities.Logging != nil {
		options = append(options, server.WithLogging())
	}

	return server.NewMCPServer("mcp-proxy", "1.0.0", options...)
}

