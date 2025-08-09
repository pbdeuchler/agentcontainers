# Assistant MCP Server

This is an MCP (Model Context Protocol) server that wraps the Assistant HTTP API for LLM usage.

## Overview

The Assistant MCP Server exposes the Assistant HTTP API's functionality as MCP tools, allowing LLMs to interact with:

- **Todos** - Create, read, update, delete, and list todo items
- **Backgrounds** - Manage background/wallpaper entries
- **Preferences** - Store and retrieve user preferences
- **Notes** - Create and manage notes

## Installation

1. Build the server:
```bash
go build -o assistant-mcp-server
```

2. Set the Assistant API URL (optional, defaults to http://localhost:8080):
```bash
export ASSISTANT_API_URL=http://localhost:8080
```

## Usage

The MCP server communicates via JSON-RPC over stdin/stdout. It's designed to be used with MCP-compatible LLM clients.

### Configuration

Add to your MCP client configuration:

```json
{
  "mcpServers": {
    "assistant": {
      "command": "/path/to/assistant-mcp-server",
      "env": {
        "ASSISTANT_API_URL": "http://localhost:8080"
      }
    }
  }
}
```

## Available Tools

### Todo Tools
- `create_todo` - Create a new todo item
- `get_todo` - Get a todo by UID
- `list_todos` - List todos with filtering and sorting
- `update_todo` - Update an existing todo
- `delete_todo` - Delete a todo

### Background Tools
- `create_background` - Create a background entry
- `get_background` - Get a background by key
- `list_backgrounds` - List all backgrounds
- `update_background` - Update a background entry
- `delete_background` - Delete a background entry

### Preferences Tools
- `create_preferences` - Create user preferences
- `get_preferences` - Get preferences by key and specifier
- `list_preferences` - List all preferences
- `update_preferences` - Update preferences
- `delete_preferences` - Delete preferences

### Notes Tools
- `create_note` - Create a new note
- `get_note` - Get a note by ID
- `list_notes` - List notes with filtering
- `update_note` - Update an existing note
- `delete_note` - Delete a note

## Data Structures

### Todo
```json
{
  "uid": "string",
  "title": "string",
  "description": "string",
  "data": "string",
  "priority": 1-4,
  "due_date": "RFC3339 timestamp",
  "recurs_on": "string",
  "marked_complete": "RFC3339 timestamp",
  "external_url": "string",
  "created_by": "string",
  "completed_by": "string",
  "created_at": "RFC3339 timestamp",
  "updated_at": "RFC3339 timestamp"
}
```

### Background
```json
{
  "key": "string",
  "value": "string",
  "created_at": "RFC3339 timestamp",
  "updated_at": "RFC3339 timestamp"
}
```

### Preferences
```json
{
  "key": "string",
  "specifier": "string",
  "data": "string (JSON)",
  "created_at": "RFC3339 timestamp",
  "updated_at": "RFC3339 timestamp"
}
```

### Notes
```json
{
  "id": "string",
  "title": "string",
  "relevant_user": "string",
  "content": "string",
  "created_at": "RFC3339 timestamp",
  "updated_at": "RFC3339 timestamp"
}
```

## Priority Levels

- 1: Low
- 2: Medium  
- 3: High
- 4: Critical

## Requirements

- Go 1.24.3 or later
- Running Assistant HTTP server