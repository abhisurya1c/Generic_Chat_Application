package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"backend/services"
)

const systemPrompt = `You are an expert SQL assistant. 
Your sole purpose is to generate valid, optimized, and secure SQL queries for PostgreSQL, MySQL, and SQL Server.

STRICT OPERATIONAL RULES:
1. OUTPUT FORMAT: Return ONLY valid SQL code. Do NOT include explanations, introductory text, or concluding remarks unless the user explicitly asks for an explanation.
2. SAFETY GUARDRAILS: You are strictly prohibited from generating queries that contain destructive or administrative commands. 
   - FORBIDDEN KEYWORDS: DROP, DELETE, TRUNCATE, ALTER, GRANT, REVOKE, SHUTDOWN.
   - If a user requests a destructive operation, respond with: "-- Error: Destructive SQL commands (DROP, DELETE, TRUNCATE) are not permitted for safety reasons."
3. CODE QUALITY: Ensure all table and column names are handled as described by the user. If schema information is provided, adhere to it strictly.
4. SYSTEM CONTEXT: You are an expert SQL assistant. Return ONLY valid SQL queries. Do NOT explain unless explicitly asked.`

var forbidden = []string{"DROP", "DELETE", "TRUNCATE", "ALTER", "GRANT", "REVOKE", "SHUTDOWN"}

type chatRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type chatResponse struct {
	Text string `json:"text"`
}

func containsForbidden(s string) bool {
	up := strings.ToUpper(s)
	for _, f := range forbidden {
		if strings.Contains(up, f) {
			return true
		}
	}
	return false
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	// guardrails
	if containsForbidden(req.Prompt) {
		resp := chatResponse{Text: "-- Error: Destructive SQL commands (DROP, DELETE, TRUNCATE) are not permitted for safety reasons."}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	fullPrompt := systemPrompt + "\n\nUser: " + req.Prompt

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()
	out, err := services.Generate(ctx, req.Model, fullPrompt)
	if err != nil {
		http.Error(w, "upstream error: "+err.Error(), http.StatusBadGateway)
		return
	}
	// If upstream produced forbidden keywords, override with error
	if containsForbidden(out) {
		out = "-- Error: Destructive SQL commands (DROP, DELETE, TRUNCATE) are not permitted for safety reasons."
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(chatResponse{Text: out})
}

func StreamChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	q := r.URL.Query()
	model := q.Get("model")
	prompt := q.Get("prompt")
	if model == "" {
		model = "sqlcoder"
	}
	if prompt == "" {
		http.Error(w, "missing prompt", http.StatusBadRequest)
		return
	}
	// guardrails
	if containsForbidden(prompt) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = io.WriteString(w, "data: -- Error: Destructive SQL commands (DROP, DELETE, TRUNCATE) are not permitted for safety reasons.\n\n")
		return
	}
	fullPrompt := systemPrompt + "\n\nUser: " + prompt

	ctx := r.Context()
	// SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	// stream from model API and forward as SSE events
	err := services.StreamToWriter(ctx, model, fullPrompt, w)
	if err != nil {
		_, _ = io.WriteString(w, "data: [STREAM ERROR]\n\n")
	}
	flusher.Flush()
}
