package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	Prompt             json.RawMessage   `json:"prompt"`
	AppendSystemPrompt *string           `json:"append_system_prompt"`
	AllowedTools       []string          `json:"allowed_tools"`
	DisallowedTools    []string          `json:"disallowed_tools"`
	ResumeSessionID    *string           `json:"resume_session_id"`
	Env                map[string]string `json:"env"`
	User               json.RawMessage   `json:"user,omitempty"` // Optional user field
}

func buildArgs(r Request) []string {
	if len(r.Prompt) == 0 {
		log.Fatal("prompt is required")
	}

	args := []string{
		"--output-format=json",
		"--dangerously-skip-permissions",
		"-p", fmt.Sprintf(`"%s"`, string(r.Prompt)),
	}

	if r.AppendSystemPrompt != nil {
		args = append(args, "--append-system-prompt", *r.AppendSystemPrompt)
	}
	if len(r.AllowedTools) > 0 {
		args = append(args, "--allowed-tools", strings.Join(r.AllowedTools, ","))
	}
	if len(r.DisallowedTools) > 0 {
		args = append(args, "--disallowed-tools", strings.Join(r.DisallowedTools, ","))
	}
	if r.ResumeSessionID != nil {
		args = append(args, "--resume", *r.ResumeSessionID)
	}
	if maxTurns := os.Getenv("MAX_TURNS"); maxTurns != "" {
		args = append(args, "--max-turns", maxTurns)
	}
	if model := os.Getenv("MODEL"); model != "" {
		args = append(args, "--model", model)
	}
	if sysPrompt := os.Getenv("SYSTEM_PROMPT"); sysPrompt != "" {
		args = append(args, "--system-prompt", sysPrompt)
	}
	return args
}

func runClaude(r Request) (string, string, error) {
	args := buildArgs(r)

	cmd := exec.Command("claude", args...)

	env := os.Environ()
	if apiKey := os.Getenv("ANTHROPIC_API_KEY"); apiKey != "" {
		env = append(env, "ANTHROPIC_API_KEY="+apiKey)
	}
	for k, v := range r.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = env

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return out.String(), stderr.String(), err
	}
	return out.String(), stderr.String(), nil
}

func handler(r Request) (events.APIGatewayProxyResponse, error) {
	stdOut, stdErr, err := runClaude(r)
	statusCode := 200
	body := stdOut
	if err != nil {
		statusCode = 400
		if stdErr != "" {
			body = stdErr
		}
	}
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       body,
	}, nil
}

func fsShim() {
	// check if /mnt/state exists
	if _, err := os.Stat("/mnt/state"); os.IsNotExist(err) {
		log.Println("fs-shim: /mnt/state does not exist, please mount it")
		return
	}
	// check if /mnt/state/projects exists
	if _, err := os.Stat("/mnt/state/projects"); os.IsNotExist(err) {
		// create /mnt/state/projects
		if err := os.MkdirAll("/mnt/state/projects", 0o755); err != nil {
			log.Fatalf("failed to create /mnt/state/projects: %v", err)
		}
	}
	// symlink /mnt/state/projects to /root/.claude/projects
	if err := os.Symlink("/mnt/state/projects", "/root/.claude/projects"); err != nil {
		log.Fatalf("failed to create symlink: %v", err)
	}
	// check if /mnt/state/__store.db exists
	if _, err := os.Stat("/mnt/state/__store.db"); os.IsNotExist(err) {
		// create an empty sqlite db at __store.db
		if _, err := os.Create("/mnt/state/__store.db"); err != nil {
			log.Fatalf("failed to create __store.db: %v", err)
		}
	}
	// symlink /mnt/state/__store.db to /root/.claude/__store.db
	if err := os.Symlink("/mnt/state/__store.db", "/root/.claude/__store.db"); err != nil {
		log.Fatalf("failed to create symlink: %v", err)
	}
}

func main() {
	if os.Getenv("LAMBDA") == "true" {
		lambda.Start(handler)
		return
	}

	if os.Getenv("FS_SHIM") != "" {
		fsShim()
	}

	if len(os.Args) < 2 {
		log.Fatal("usage: program '{\"prompt\":\"...\"}'")
	}
	var req Request
	if err := json.Unmarshal([]byte(os.Args[1]), &req); err != nil {
		log.Fatalf("invalid json: %v", err)
	}
	resp, err := handler(req)
	fmt.Println(err)
	fmt.Println(resp)
}
