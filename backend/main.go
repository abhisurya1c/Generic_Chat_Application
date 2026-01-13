package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/abhisurya1c/Generic_Chat_Application/backend/db"
	"github.com/abhisurya1c/Generic_Chat_Application/backend/handlers"
	"github.com/abhisurya1c/Generic_Chat_Application/backend/middleware"
)

func main() {
	// 1. Initialize DB
	connStr := "postgres://user:password@localhost:5432/sqlchat?sslmode=disable"
	if err := db.Connect(connStr); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	// 2. Auth Routes (Public)
	mux.HandleFunc("/api/register", handlers.RegisterHandler)
	mux.HandleFunc("/api/login", handlers.LoginHandler)

	// 3. Protected Routes
	// Helper to wrap handlers with AuthMiddleware
	protected := func(h http.HandlerFunc) http.HandlerFunc {
		return middleware.AuthMiddleware(h)
	}

	mux.HandleFunc("/api/chat", protected(handlers.ChatHandler))
	mux.HandleFunc("/api/chat/stream", protected(handlers.StreamChatHandler))

	mux.HandleFunc("/api/history/chats", protected(handlers.GetChatsHandler))
	mux.HandleFunc("/api/history/messages", protected(handlers.GetMessagesHandler))
	mux.HandleFunc("/api/history/delete", protected(handlers.DeleteChatHandler))

	// 4. Global Middleware (CORS)
	handler := middleware.EnableCORS(mux)

	port := 8080
	fmt.Printf("Server starting on port %d...\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), handler); err != nil {
		log.Fatal(err)
	}
}
