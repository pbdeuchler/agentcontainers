package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type MCPServer struct {
	baseURL string
}

type MCPRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type MCPResponse struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      any       `json:"id"`
	Result  any       `json:"result,omitempty"`
	Error   *MCPError `json:"error,omitempty"`
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type InitializeParams struct {
	ProtocolVersion string         `json:"protocolVersion"`
	Capabilities    map[string]any `json:"capabilities"`
	ClientInfo      ClientInfo     `json:"clientInfo"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResult struct {
	ProtocolVersion string         `json:"protocolVersion"`
	Capabilities    map[string]any `json:"capabilities"`
	ServerInfo      ServerInfo     `json:"serverInfo"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
}

type CallToolParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments,omitempty"`
}

type CallToolResult struct {
	Content []Content `json:"content"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Todo struct {
	UID            string     `json:"uid"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Data           string     `json:"data"`
	Priority       int        `json:"priority"`
	DueDate        *time.Time `json:"due_date"`
	RecursOn       string     `json:"recurs_on"`
	MarkedComplete *time.Time `json:"marked_complete"`
	ExternalURL    string     `json:"external_url"`
	CreatedBy      string     `json:"created_by"`
	CompletedBy    string     `json:"completed_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type Background struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Preferences struct {
	Key       string    `json:"key"`
	Specifier string    `json:"specifier"`
	Data      string    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Notes struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	RelevantUser string    `json:"relevant_user"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewMCPServer(baseURL string) *MCPServer {
	return &MCPServer{baseURL: baseURL}
}

func main() {
	baseURL := os.Getenv("ASSISTANT_API_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	server := NewMCPServer(baseURL)
	server.Run()
}

func (s *MCPServer) Run() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var request MCPRequest
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			s.sendError(request.ID, -32700, "Parse error")
			continue
		}

		switch request.Method {
		case "initialize":
			s.handleInitialize(request)
		case "tools/list":
			s.handleToolsList(request)
		case "tools/call":
			s.handleToolsCall(request)
		default:
			s.sendError(request.ID, -32601, "Method not found")
		}
	}
}

func (s *MCPServer) handleInitialize(request MCPRequest) {
	result := InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: map[string]any{
			"tools": map[string]any{},
		},
		ServerInfo: ServerInfo{
			Name:    "assistant-mcp-server",
			Version: "1.0.0",
		},
	}

	s.sendResponse(request.ID, result)
}

func (s *MCPServer) handleToolsList(request MCPRequest) {
	tools := []Tool{
		{
			Name:        "create_todo",
			Description: "Create a new todo item",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"title":        map[string]any{"type": "string", "description": "Todo title"},
					"description":  map[string]any{"type": "string", "description": "Todo description"},
					"data":         map[string]any{"type": "string", "description": "Additional data"},
					"priority":     map[string]any{"type": "integer", "description": "Priority (1=low, 2=medium, 3=high, 4=critical)"},
					"due_date":     map[string]any{"type": "string", "description": "Due date in RFC3339 format"},
					"recurs_on":    map[string]any{"type": "string", "description": "Recurrence pattern"},
					"external_url": map[string]any{"type": "string", "description": "External URL"},
					"created_by":   map[string]any{"type": "string", "description": "Creator identifier"},
				},
				"required": []string{"title"},
			},
		},
		{
			Name:        "get_todo",
			Description: "Get a todo item by UID",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"uid": map[string]any{"type": "string", "description": "Todo UID"},
				},
				"required": []string{"uid"},
			},
		},
		{
			Name:        "list_todos",
			Description: "List todos with optional filters and sorting",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"limit":      map[string]any{"type": "integer", "description": "Number of items to return"},
					"offset":     map[string]any{"type": "integer", "description": "Number of items to skip"},
					"sort_by":    map[string]any{"type": "string", "description": "Field to sort by"},
					"sort_dir":   map[string]any{"type": "string", "description": "Sort direction (asc/desc)"},
					"title":      map[string]any{"type": "string", "description": "Filter by title"},
					"priority":   map[string]any{"type": "string", "description": "Filter by priority"},
					"created_by": map[string]any{"type": "string", "description": "Filter by creator"},
				},
			},
		},
		{
			Name:        "complete_todo",
			Description: "Complete a todo item",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"uid":          map[string]any{"type": "string", "description": "Todo UID"},
					"completed_by": map[string]any{"type": "string", "description": "Completer identifier"},
				},
				"required": []string{"uid"},
			},
		},
		{
			Name:        "update_todo",
			Description: "Update a todo item",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"uid":             map[string]any{"type": "string", "description": "Todo UID"},
					"title":           map[string]any{"type": "string", "description": "Todo title"},
					"description":     map[string]any{"type": "string", "description": "Todo description"},
					"data":            map[string]any{"type": "string", "description": "Additional data"},
					"priority":        map[string]any{"type": "integer", "description": "Priority (1=low, 2=medium, 3=high, 4=critical)"},
					"due_date":        map[string]any{"type": "string", "description": "Due date in RFC3339 format"},
					"recurs_on":       map[string]any{"type": "string", "description": "Recurrence pattern"},
					"marked_complete": map[string]any{"type": "string", "description": "Completion timestamp in RFC3339 format"},
					"external_url":    map[string]any{"type": "string", "description": "External URL"},
					"completed_by":    map[string]any{"type": "string", "description": "Completer identifier"},
				},
				"required": []string{"uid"},
			},
		},
		{
			Name:        "delete_todo",
			Description: "Delete a todo item",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"uid": map[string]any{"type": "string", "description": "Todo UID"},
				},
				"required": []string{"uid"},
			},
		},
		{
			Name:        "create_background",
			Description: "Create a new background/wallpaper entry",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"key":   map[string]any{"type": "string", "description": "Background key"},
					"value": map[string]any{"type": "string", "description": "Background value (URL or data)"},
				},
				"required": []string{"key", "value"},
			},
		},
		{
			Name:        "get_background",
			Description: "Get a background by key",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"key": map[string]any{"type": "string", "description": "Background key"},
				},
				"required": []string{"key"},
			},
		},
		{
			Name:        "list_backgrounds",
			Description: "List backgrounds",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"limit":  map[string]any{"type": "integer", "description": "Number of items to return"},
					"offset": map[string]any{"type": "integer", "description": "Number of items to skip"},
				},
			},
		},
		{
			Name:        "update_background",
			Description: "Update a background entry",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"key":   map[string]any{"type": "string", "description": "Background key"},
					"value": map[string]any{"type": "string", "description": "Background value (URL or data)"},
				},
				"required": []string{"key", "value"},
			},
		},
		{
			Name:        "delete_background",
			Description: "Delete a background entry",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"key": map[string]any{"type": "string", "description": "Background key"},
				},
				"required": []string{"key"},
			},
		},
		{
			Name:        "create_preferences",
			Description: "Create new preferences",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"key":       map[string]any{"type": "string", "description": "Preference key"},
					"specifier": map[string]any{"type": "string", "description": "Preference specifier"},
					"data":      map[string]any{"type": "string", "description": "Preference data (JSON)"},
				},
				"required": []string{"key", "specifier", "data"},
			},
		},
		{
			Name:        "get_preferences",
			Description: "Get preferences by key and specifier",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"key":       map[string]any{"type": "string", "description": "Preference key"},
					"specifier": map[string]any{"type": "string", "description": "Preference specifier"},
				},
				"required": []string{"key", "specifier"},
			},
		},
		{
			Name:        "list_preferences",
			Description: "List preferences",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"limit":  map[string]any{"type": "integer", "description": "Number of items to return"},
					"offset": map[string]any{"type": "integer", "description": "Number of items to skip"},
				},
			},
		},
		{
			Name:        "update_preferences",
			Description: "Update preferences",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"key":       map[string]any{"type": "string", "description": "Preference key"},
					"specifier": map[string]any{"type": "string", "description": "Preference specifier"},
					"data":      map[string]any{"type": "string", "description": "Preference data (JSON)"},
				},
				"required": []string{"key", "specifier", "data"},
			},
		},
		{
			Name:        "delete_preferences",
			Description: "Delete preferences",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"key":       map[string]any{"type": "string", "description": "Preference key"},
					"specifier": map[string]any{"type": "string", "description": "Preference specifier"},
				},
				"required": []string{"key", "specifier"},
			},
		},
		{
			Name:        "create_note",
			Description: "Create a new note",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"title":         map[string]any{"type": "string", "description": "Note title"},
					"relevant_user": map[string]any{"type": "string", "description": "Relevant user"},
					"content":       map[string]any{"type": "string", "description": "Note content"},
				},
				"required": []string{"title", "content"},
			},
		},
		{
			Name:        "get_note",
			Description: "Get a note by ID",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"id": map[string]any{"type": "string", "description": "Note ID"},
				},
				"required": []string{"id"},
			},
		},
		{
			Name:        "list_notes",
			Description: "List notes with optional filters",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"limit":         map[string]any{"type": "integer", "description": "Number of items to return"},
					"offset":        map[string]any{"type": "integer", "description": "Number of items to skip"},
					"sort_by":       map[string]any{"type": "string", "description": "Field to sort by"},
					"sort_dir":      map[string]any{"type": "string", "description": "Sort direction (asc/desc)"},
					"title":         map[string]any{"type": "string", "description": "Filter by title"},
					"relevant_user": map[string]any{"type": "string", "description": "Filter by relevant user"},
				},
			},
		},
		{
			Name:        "update_note",
			Description: "Update a note",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"id":            map[string]any{"type": "string", "description": "Note ID"},
					"title":         map[string]any{"type": "string", "description": "Note title"},
					"relevant_user": map[string]any{"type": "string", "description": "Relevant user"},
					"content":       map[string]any{"type": "string", "description": "Note content"},
				},
				"required": []string{"id"},
			},
		},
		{
			Name:        "delete_note",
			Description: "Delete a note",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"id": map[string]any{"type": "string", "description": "Note ID"},
				},
				"required": []string{"id"},
			},
		},
	}

	result := map[string]any{
		"tools": tools,
	}

	s.sendResponse(request.ID, result)
}

func (s *MCPServer) handleToolsCall(request MCPRequest) {
	var params CallToolParams
	if err := json.Unmarshal(request.Params, &params); err != nil {
		s.sendError(request.ID, -32602, "Invalid params")
		return
	}

	result, err := s.callTool(params.Name, params.Arguments)
	if err != nil {
		s.sendError(request.ID, -32603, err.Error())
		return
	}

	s.sendResponse(request.ID, result)
}

func (s *MCPServer) callTool(name string, arguments json.RawMessage) (CallToolResult, error) {
	switch name {
	case "create_todo":
		return s.createTodo(arguments)
	case "get_todo":
		return s.getTodo(arguments)
	case "list_todos":
		return s.listTodos(arguments)
	case "complete_todo":
		return s.completeTodo(arguments)
	case "update_todo":
		return s.updateTodo(arguments)
	case "delete_todo":
		return s.deleteTodo(arguments)
	case "create_background":
		return s.createBackground(arguments)
	case "get_background":
		return s.getBackground(arguments)
	case "list_backgrounds":
		return s.listBackgrounds(arguments)
	case "update_background":
		return s.updateBackground(arguments)
	case "delete_background":
		return s.deleteBackground(arguments)
	case "create_preferences":
		return s.createPreferences(arguments)
	case "get_preferences":
		return s.getPreferences(arguments)
	case "list_preferences":
		return s.listPreferences(arguments)
	case "update_preferences":
		return s.updatePreferences(arguments)
	case "delete_preferences":
		return s.deletePreferences(arguments)
	case "create_note":
		return s.createNote(arguments)
	case "get_note":
		return s.getNote(arguments)
	case "list_notes":
		return s.listNotes(arguments)
	case "update_note":
		return s.updateNote(arguments)
	case "delete_note":
		return s.deleteNote(arguments)
	default:
		return CallToolResult{}, fmt.Errorf("unknown tool: %s", name)
	}
}

func (s *MCPServer) sendResponse(id any, result any) {
	response := MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}

	data, _ := json.Marshal(response)
	fmt.Println(string(data))
}

func (s *MCPServer) sendError(id any, code int, message string) {
	response := MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &MCPError{
			Code:    code,
			Message: message,
		},
	}

	data, _ := json.Marshal(response)
	fmt.Println(string(data))
}

func parseTime(timeStr string) *time.Time {
	if timeStr == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil
	}
	return &t
}
