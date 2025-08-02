# AgentContainers

A containerized Claude Code service that provides secure, isolated execution of AI assistant requests through Docker containers and AWS Lambda integration.

## Overview

Run Claude Code (or, in the future, maybe other agentic tooling) via a Lambda Function or HTTP Request. This project was started to help build an open source LLM personal assistant but can be used as a baseline for how to abstract agentic tooling to run via stateless container/HTTP request.

## Quick Start

### Using Docker Compose

```bash
# Build containers
docker-compose build shim
docker-compose build assistant
```

## Configuration

### Environment Variables

- `ANTHROPIC_API_KEY`: Required for Claude API access
- `LAMBDA`: Set to "true" for Lambda execution mode
- `MAX_TURNS`: Limit conversation turns
- `MODEL`: Specify Claude model (e.g., "claude-sonnet-4-20250514")
- `SYSTEM_PROMPT`: Override default system prompt

### Request Format

#### Detailed Request Example

```json
{
  "prompt": "Your request to Claude",
  "append_system_prompt": "Additional system context",
  "allowed_tools": ["Read", "Write", "Bash"],
  "disallowed_tools": ["WebFetch"],
  "resume_session_id": "session-123",
  "env": {
    "CUSTOM_VAR": "value"
  }
}
```

#### Local Testing

> `docker-compose run --rm assistant '{"prompt":"who are you"}'`

```
<nil>
{200 map[] map[] {"type":"result","subtype":"success","is_error":false,"duration_ms":8245,"duration_api_ms":7691,"num_turns":1,"result":"Hello! I'm Clarice, your personal assistant. I'm here to help you and your partner manage your daily lives more efficiently and smoothly.\n\nMy role is to support you both with a wide range of tasks, including:\n\n- **Task and calendar management** - keeping track of your todos, scheduling, and making sure nothing falls through the cracks\n- **Information curation** - providing relevant news briefings and summarizing important information\n- **Home management** - helping with meal planning, grocery shopping, home improvement projects, and general household organization\n- **Personal support** - offering advice on everything from finances to childcare, and even helping with relationship matters when needed\n- **Vacation and event planning** - coordinating schedules and making arrangements\n\nI'm designed to be proactive, anticipating your needs and suggesting optimizations to make your lives easier. I maintain appropriate boundaries and privacy between you and your partner while facilitating great communication and coordination.\n\nWhat can I help you with today? Whether it's a specific task, planning something for the week ahead, or just getting organized, I'm here to assist!","session_id":"cd67c0e6-7de1-4325-8710-ac95707523c9","total_cost_usd":0.05214525,"usage":{"input_tokens":3,"cache_creation_input_tokens":12927,"cache_read_input_tokens":0,"output_tokens":244,"server_tool_use":{"web_search_requests":0},"service_tier":"standard"}}
 false}
```

## Development

### Prerequisites

- Go 1.23.3+
- Docker and Docker Compose
- AWS CLI (for Lambda deployment)

### Container Development

The assistant container includes:

- Ubuntu 22.04 base with development tools
- Claude Code CLI
- Custom "Clarice" assistant persona
- Git, fzf, zsh with oh-my-zsh
- Development utilities (jq, nano, vim, curl)

## License

This project is licensed under the GNU GPLv3 License with the [Commons Clause License Condition v1.0](https://commonsclause.com/).
