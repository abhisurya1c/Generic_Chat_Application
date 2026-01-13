package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/abhisurya1c/Generic_Chat_Application/backend/db"
	"github.com/abhisurya1c/Generic_Chat_Application/backend/middleware"
)

type Chat struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt string    `json:"created_at"`
	Messages  []Message `json:"messages,omitempty"`
}

type Message struct {
	ID        int    `json:"id"`
	Role      string `json:"role"`
	Content   string `json:"text"` // Mapped to 'text' for frontend compatibility
	CreatedAt string `json:"created_at"`
}

func GetChatsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	rows, err := db.DB.Query("SELECT id, title, created_at FROM chats WHERE user_id = $1 ORDER BY created_at DESC", userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	chats := []Chat{}
	for rows.Next() {
		var c Chat
		if err := rows.Scan(&c.ID, &c.Title, &c.CreatedAt); err != nil {
			continue
		}
		chats = append(chats, c)
	}

	json.NewEncoder(w).Encode(chats)
}

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	chatIDStr := r.URL.Query().Get("chat_id")
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	// Verify ownership
	var OwnerID int
	err = db.DB.QueryRow("SELECT user_id FROM chats WHERE id = $1", chatID).Scan(&OwnerID)
	if err == sql.ErrNoRows {
		http.Error(w, "Chat not found", http.StatusNotFound)
		return
	}
	if OwnerID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	rows, err := db.DB.Query("SELECT id, role, content, created_at FROM messages WHERE chat_id = $1 ORDER BY id ASC", chatID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	messages := []Message{}
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
			continue
		}
		messages = append(messages, m)
	}

	json.NewEncoder(w).Encode(messages)
}

func DeleteChatHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	chatIDStr := r.URL.Query().Get("chat_id")
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec("DELETE FROM chats WHERE id = $1 AND user_id = $2", chatID, userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Chat not found or access denied", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
