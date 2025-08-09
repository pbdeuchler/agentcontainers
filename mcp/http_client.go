package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func (s *MCPServer) createTodo(arguments json.RawMessage) (CallToolResult, error) {
	var args struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Data        string `json:"data"`
		Priority    int    `json:"priority"`
		DueDate     string `json:"due_date"`
		RecursOn    string `json:"recurs_on"`
		ExternalURL string `json:"external_url"`
		CreatedBy   string `json:"created_by"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return CallToolResult{}, err
	}

	data, _ := json.Marshal(args)
	resp, err := http.Post(s.baseURL+"/todos", "application/json", bytes.NewReader(data))
	if err != nil {
		return CallToolResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CallToolResult{}, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	return CallToolResult{
		Content: []Content{{Type: "text", Text: string(body)}},
	}, nil
}

func (s *MCPServer) getTodo(arguments json.RawMessage) (CallToolResult, error) {
	var args struct {
		UID string `json:"uid"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return CallToolResult{}, err
	}

	resp, err := http.Get(s.baseURL + "/todos/" + args.UID)
	if err != nil {
		return CallToolResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CallToolResult{}, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	return CallToolResult{
		Content: []Content{{Type: "text", Text: string(body)}},
	}, nil
}

func (s *MCPServer) listTodos(arguments json.RawMessage) (CallToolResult, error) {
	var args struct {
		Limit     *int   `json:"limit"`
		Offset    *int   `json:"offset"`
		SortBy    string `json:"sort_by"`
		SortDir   string `json:"sort_dir"`
		Title     string `json:"title"`
		Priority  string `json:"priority"`
		CreatedBy string `json:"created_by"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return CallToolResult{}, err
	}

	params := url.Values{}
	if args.Limit != nil {
		params.Set("limit", strconv.Itoa(*args.Limit))
	}
	if args.Offset != nil {
		params.Set("offset", strconv.Itoa(*args.Offset))
	}
	if args.SortBy != "" {
		params.Set("sort_by", args.SortBy)
	}
	if args.SortDir != "" {
		params.Set("sort_dir", args.SortDir)
	}
	if args.Title != "" {
		params.Set("title", args.Title)
	}
	if args.Priority != "" {
		params.Set("priority", args.Priority)
	}
	if args.CreatedBy != "" {
		params.Set("created_by", args.CreatedBy)
	}

	url := s.baseURL + "/todos"
	if len(params) > 0 {
		url += "?" + params.Encode()
	}

	resp, err := http.Get(url)
	if err != nil {
		return CallToolResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CallToolResult{}, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	return CallToolResult{
		Content: []Content{{Type: "text", Text: string(body)}},
	}, nil
}

func (s *MCPServer) completeTodo(arguments json.RawMessage) (CallToolResult, error) {
	return s.updateTodo(arguments) // Reuse updateTodo for completion
}

func (s *MCPServer) updateTodo(arguments json.RawMessage) (CallToolResult, error) {
	var args struct {
		UID            string `json:"uid"`
		Title          string `json:"title"`
		Description    string `json:"description"`
		Data           string `json:"data"`
		Priority       int    `json:"priority"`
		DueDate        string `json:"due_date"`
		RecursOn       string `json:"recurs_on"`
		MarkedComplete string `json:"marked_complete"`
		ExternalURL    string `json:"external_url"`
		CompletedBy    string `json:"completed_by"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return CallToolResult{}, err
	}

	todo := Todo{
		Title:          args.Title,
		Description:    args.Description,
		Data:           args.Data,
		Priority:       args.Priority,
		DueDate:        parseTime(args.DueDate),
		RecursOn:       args.RecursOn,
		MarkedComplete: parseTime(args.MarkedComplete),
		ExternalURL:    args.ExternalURL,
		CompletedBy:    args.CompletedBy,
	}

	data, _ := json.Marshal(todo)
	req, _ := http.NewRequest(http.MethodPut, s.baseURL+"/todos/"+args.UID, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return CallToolResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CallToolResult{}, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	return CallToolResult{
		Content: []Content{{Type: "text", Text: string(body)}},
	}, nil
}

func (s *MCPServer) deleteTodo(arguments json.RawMessage) (CallToolResult, error) {
	var args struct {
		UID string `json:"uid"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return CallToolResult{}, err
	}

	req, _ := http.NewRequest(http.MethodDelete, s.baseURL+"/todos/"+args.UID, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return CallToolResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return CallToolResult{}, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return CallToolResult{
		Content: []Content{{Type: "text", Text: "Todo deleted successfully"}},
	}, nil
}

func (s *MCPServer) createBackground(arguments json.RawMessage) (CallToolResult, error) {
	var args struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return CallToolResult{}, err
	}

	background := Background{
		Key:   args.Key,
		Value: args.Value,
	}

	data, _ := json.Marshal(background)
	resp, err := http.Post(s.baseURL+"/backgrounds", "application/json", bytes.NewReader(data))
	if err != nil {
		return CallToolResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CallToolResult{}, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	return CallToolResult{
		Content: []Content{{Type: "text", Text: string(body)}},
	}, nil
}

func (s *MCPServer) getBackground(arguments json.RawMessage) (CallToolResult, error) {
	var args struct {
		Key string `json:"key"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return CallToolResult{}, err
	}

	resp, err := http.Get(s.baseURL + "/backgrounds/" + args.Key)
	if err != nil {
		return CallToolResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CallToolResult{}, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	return CallToolResult{
		Content: []Content{{Type: "text", Text: string(body)}},
	}, nil
}

func (s *MCPServer) listBackgrounds(arguments json.RawMessage) (CallToolResult, error) {
	var args struct {
		Limit  *int `json:"limit"`
		Offset *int `json:"offset"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return CallToolResult{}, err
	}

	params := url.Values{}
	if args.Limit != nil {
		params.Set("limit", strconv.Itoa(*args.Limit))
	}
	if args.Offset != nil {
		params.Set("offset", strconv.Itoa(*args.Offset))
	}

	url := s.baseURL + "/backgrounds"
	if len(params) > 0 {
		url += "?" + params.Encode()
	}

	resp, err := http.Get(url)
	if err != nil {
		return CallToolResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CallToolResult{}, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	return CallToolResult{
		Content: []Content{{Type: "text", Text: string(body)}},
	}, nil
}

func (s *MCPServer) updateBackground(arguments json.RawMessage) (CallToolResult, error) {
	var args struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return CallToolResult{}, err
	}

	background := Background{
		Key:   args.Key,
		Value: args.Value,
	}

	data, _ := json.Marshal(background)
	req, _ := http.NewRequest(http.MethodPut, s.baseURL+"/backgrounds/"+args.Key, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return CallToolResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CallToolResult{}, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	return CallToolResult{
		Content: []Content{{Type: "text", Text: string(body)}},
	}, nil
}

func (s *MCPServer) deleteBackground(arguments json.RawMessage) (CallToolResult, error) {
	var args struct {
		Key string `json:"key"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return CallToolResult{}, err
	}

	req, _ := http.NewRequest(http.MethodDelete, s.baseURL+"/backgrounds/"+args.Key, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return CallToolResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return CallToolResult{}, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return CallToolResult{
		Content: []Content{{Type: "text", Text: "Background deleted successfully"}},
	}, nil
}

func (s *MCPServer) createPreferences(arguments json.RawMessage) (CallToolResult, error) {
	var args struct {
		Key       string `json:"key"`
		Specifier string `json:"specifier"`
		Data      string `json:"data"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return CallToolResult{}, err
	}

	preferences := Preferences{
		Key:       args.Key,
		Specifier: args.Specifier,
		Data:      args.Data,
	}

	data, _ := json.Marshal(preferences)
	resp, err := http.Post(s.baseURL+"/preferences", "application/json", bytes.NewReader(data))
	if err != nil {
		return CallToolResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CallToolResult{}, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	return CallToolResult{
		Content: []Content{{Type: "text", Text: string(body)}},
	}, nil
}

func (s *MCPServer) getPreferences(arguments json.RawMessage) (CallToolResult, error) {
	var args struct {
		Key       string `json:"key"`
		Specifier string `json:"specifier"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return CallToolResult{}, err
	}

	params := url.Values{}
	params.Set("specifier", args.Specifier)

	url := s.baseURL + "/preferences/" + args.Key + "?" + params.Encode()

	resp, err := http.Get(url)
	if err != nil {
		return CallToolResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CallToolResult{}, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	return CallToolResult{
		Content: []Content{{Type: "text", Text: string(body)}},
	}, nil
}

func (s *MCPServer) listPreferences(arguments json.RawMessage) (CallToolResult, error) {
	var args struct {
		Limit  *int `json:"limit"`
		Offset *int `json:"offset"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return CallToolResult{}, err
	}

	params := url.Values{}
	if args.Limit != nil {
		params.Set("limit", strconv.Itoa(*args.Limit))
	}
	if args.Offset != nil {
		params.Set("offset", strconv.Itoa(*args.Offset))
	}

	url := s.baseURL + "/preferences"
	if len(params) > 0 {
		url += "?" + params.Encode()
	}

	resp, err := http.Get(url)
	if err != nil {
		return CallToolResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CallToolResult{}, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	return CallToolResult{
		Content: []Content{{Type: "text", Text: string(body)}},
	}, nil
}

func (s *MCPServer) updatePreferences(arguments json.RawMessage) (CallToolResult, error) {
	var args struct {
		Key       string `json:"key"`
		Specifier string `json:"specifier"`
		Data      string `json:"data"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return CallToolResult{}, err
	}

	preferences := Preferences{
		Key:       args.Key,
		Specifier: args.Specifier,
		Data:      args.Data,
	}

	data, _ := json.Marshal(preferences)
	req, _ := http.NewRequest(http.MethodPut, s.baseURL+"/preferences/"+args.Key+"?specifier="+args.Specifier, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return CallToolResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CallToolResult{}, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	return CallToolResult{
		Content: []Content{{Type: "text", Text: string(body)}},
	}, nil
}

func (s *MCPServer) deletePreferences(arguments json.RawMessage) (CallToolResult, error) {
	var args struct {
		Key       string `json:"key"`
		Specifier string `json:"specifier"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return CallToolResult{}, err
	}

	req, _ := http.NewRequest(http.MethodDelete, s.baseURL+"/preferences/"+args.Key+"?specifier="+args.Specifier, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return CallToolResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return CallToolResult{}, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return CallToolResult{
		Content: []Content{{Type: "text", Text: "Preferences deleted successfully"}},
	}, nil
}

func (s *MCPServer) createNote(arguments json.RawMessage) (CallToolResult, error) {
	var args struct {
		Title        string `json:"title"`
		RelevantUser string `json:"relevant_user"`
		Content      string `json:"content"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return CallToolResult{}, err
	}

	note := Notes{
		Title:        args.Title,
		RelevantUser: args.RelevantUser,
		Content:      args.Content,
	}

	data, _ := json.Marshal(note)
	resp, err := http.Post(s.baseURL+"/notes", "application/json", bytes.NewReader(data))
	if err != nil {
		return CallToolResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CallToolResult{}, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	return CallToolResult{
		Content: []Content{{Type: "text", Text: string(body)}},
	}, nil
}

func (s *MCPServer) getNote(arguments json.RawMessage) (CallToolResult, error) {
	var args struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return CallToolResult{}, err
	}

	resp, err := http.Get(s.baseURL + "/notes/" + args.ID)
	if err != nil {
		return CallToolResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CallToolResult{}, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	return CallToolResult{
		Content: []Content{{Type: "text", Text: string(body)}},
	}, nil
}

func (s *MCPServer) listNotes(arguments json.RawMessage) (CallToolResult, error) {
	var args struct {
		Limit        *int   `json:"limit"`
		Offset       *int   `json:"offset"`
		SortBy       string `json:"sort_by"`
		SortDir      string `json:"sort_dir"`
		Title        string `json:"title"`
		RelevantUser string `json:"relevant_user"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return CallToolResult{}, err
	}

	params := url.Values{}
	if args.Limit != nil {
		params.Set("limit", strconv.Itoa(*args.Limit))
	}
	if args.Offset != nil {
		params.Set("offset", strconv.Itoa(*args.Offset))
	}
	if args.SortBy != "" {
		params.Set("sort_by", args.SortBy)
	}
	if args.SortDir != "" {
		params.Set("sort_dir", args.SortDir)
	}
	if args.Title != "" {
		params.Set("title", args.Title)
	}
	if args.RelevantUser != "" {
		params.Set("relevant_user", args.RelevantUser)
	}

	url := s.baseURL + "/notes"
	if len(params) > 0 {
		url += "?" + params.Encode()
	}

	resp, err := http.Get(url)
	if err != nil {
		return CallToolResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CallToolResult{}, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	return CallToolResult{
		Content: []Content{{Type: "text", Text: string(body)}},
	}, nil
}

func (s *MCPServer) updateNote(arguments json.RawMessage) (CallToolResult, error) {
	var args struct {
		ID           string `json:"id"`
		Title        string `json:"title"`
		RelevantUser string `json:"relevant_user"`
		Content      string `json:"content"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return CallToolResult{}, err
	}

	note := Notes{
		Title:        args.Title,
		RelevantUser: args.RelevantUser,
		Content:      args.Content,
	}

	data, _ := json.Marshal(note)
	req, _ := http.NewRequest(http.MethodPut, s.baseURL+"/notes/"+args.ID, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return CallToolResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CallToolResult{}, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	return CallToolResult{
		Content: []Content{{Type: "text", Text: string(body)}},
	}, nil
}

func (s *MCPServer) deleteNote(arguments json.RawMessage) (CallToolResult, error) {
	var args struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return CallToolResult{}, err
	}

	req, _ := http.NewRequest(http.MethodDelete, s.baseURL+"/notes/"+args.ID, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return CallToolResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return CallToolResult{}, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return CallToolResult{
		Content: []Content{{Type: "text", Text: "Note deleted successfully"}},
	}, nil
}

