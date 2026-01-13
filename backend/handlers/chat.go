package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/abhisurya1c/Generic_Chat_Application/backend/db"
	"github.com/abhisurya1c/Generic_Chat_Application/backend/middleware"
	"github.com/abhisurya1c/Generic_Chat_Application/backend/services"
)

const SystemPrompt = `You are a helpful, harmless, and honest AI assistant. 
Your goal is to assist the user with their requests to the best of your ability.
You can answer questions, help with coding, writing, analysis, and general conversation.
Be concise and clear in your responses.
`

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(int)

	var req struct {
		Prompt string `json:"prompt"`
		Model  string `json:"model"`
		ChatID *int   `json:"chat_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Model == "" {
		req.Model = "llama3"
	}

	// 1. Create or Get Chat
	var chatID int
	if req.ChatID == nil || *req.ChatID == 0 {
		// Create new chat
		title := req.Prompt
		if len(title) > 30 {
			title = title[:30] + "..."
		}
		err := db.DB.QueryRow("INSERT INTO chats (user_id, title) VALUES ($1, $2) RETURNING id", userID, title).Scan(&chatID)
		if err != nil {
			http.Error(w, "Failed to create chat", http.StatusInternalServerError)
			return
		}
	} else {
		chatID = *req.ChatID
		// Verify ownership
		var ownerID int
		err := db.DB.QueryRow("SELECT user_id FROM chats WHERE id = $1", chatID).Scan(&ownerID)
		if err != nil || ownerID != userID {
			http.Error(w, "Chat not found", http.StatusNotFound)
			return
		}
	}

	// 2. Save User Message
	_, err := db.DB.Exec("INSERT INTO messages (chat_id, role, content) VALUES ($1, 'user', $2)", chatID, req.Prompt)
	if err != nil {
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	fullPrompt := fmt.Sprintf("%s\n\nUser: %s\n\nAssistant:", SystemPrompt, req.Prompt)

	// 3. Call AI
	response, err := services.CallOllama(fullPrompt, req.Model, false, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Backend error: %v", err), http.StatusInternalServerError)
		return
	}

	// 4. Save AI Response
	_, err = db.DB.Exec("INSERT INTO messages (chat_id, role, content) VALUES ($1, 'ai', $2)", chatID, response)
	if err != nil {
		fmt.Printf("Failed to save AI response: %v\n", err)
	}

	json.NewEncoder(w).Encode(struct {
		Response string `json:"response"`
		ChatID   int    `json:"chat_id"`
	}{
		Response: response,
		ChatID:   chatID,
	})
}

func StreamChatHandler(w http.ResponseWriter, r *http.Request) {
	prompt := r.URL.Query().Get("prompt")
	model := r.URL.Query().Get("model")
	chatIDStr := r.URL.Query().Get("chat_id")
	userID := r.Context().Value(middleware.UserIDKey).(int)

	if prompt == "" {
		http.Error(w, "Prompt is required", http.StatusBadRequest)
		return
	}
	if model == "" {
		model = "llama3"
	}

	// 1. Create or Get Chat
	var chatID int
	if chatIDStr == "" || chatIDStr == "0" {
		title := prompt
		if len(title) > 30 {
			title = title[:30] + "..."
		}
		err := db.DB.QueryRow("INSERT INTO chats (user_id, title) VALUES ($1, $2) RETURNING id", userID, title).Scan(&chatID)
		if err != nil {
			http.Error(w, "Failed to create chat", http.StatusInternalServerError)
			return
		}
	} else {
		var err error
		chatID, err = strconv.Atoi(chatIDStr)
		if err != nil {
			http.Error(w, "Invalid chat ID", http.StatusBadRequest)
			return
		}
		// Verify ownership
		var ownerID int
		err = db.DB.QueryRow("SELECT user_id FROM chats WHERE id = $1", chatID).Scan(&ownerID)
		if err != nil || ownerID != userID {
			http.Error(w, "Chat not found", http.StatusNotFound)
			return
		}
	}

	// 2. Save User Message
	_, err := db.DB.Exec("INSERT INTO messages (chat_id, role, content) VALUES ($1, 'user', $2)", chatID, prompt)
	if err != nil {
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	fullPrompt := fmt.Sprintf("%s\n\nUser: %s\n\nAssistant:", SystemPrompt, prompt)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Send ID first so frontend knows context
	idData, _ := json.Marshal(map[string]int{"chat_id": chatID})
	fmt.Fprintf(w, "data: %s\n\n", idData)
	flusher.Flush()

	fullResponse, err := services.CallOllama(fullPrompt, model, true, func(chunk string) error {
		data, _ := json.Marshal(map[string]string{"chunk": chunk})
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
		return nil
	})

	if err != nil {
		fmt.Fprintf(w, "event: error\ndata: %v\n\n", err)
		return
	}

	// 3. Save AI Response
	if fullResponse != "" {
		db.DB.Exec("INSERT INTO messages (chat_id, role, content) VALUES ($1, 'ai', $2)", chatID, fullResponse)
	}
}
