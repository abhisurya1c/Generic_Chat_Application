package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"backend/handlers"
	"backend/middleware"
)

func main() {
	mux := http.NewServeMux()

	// API
	mux.HandleFunc("/api/chat", handlers.ChatHandler)
	mux.HandleFunc("/api/chat/stream", handlers.StreamChatHandler)

	// Static frontend
	fs := http.FileServer(http.Dir("./frontend"))
	mux.Handle("/", fs)

	handler := middleware.CORS(mux)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	go func() {
		log.Println("server listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Println("server stopped")
}
