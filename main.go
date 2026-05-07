package main

import (
	"log"
	"net/http"
	"os"
	"todo/internal/api"
	"todo/internal/store"
)

func main() {
	if err := store.Init(); err != nil {
		log.Fatalf("store init: %v", err)
	}
	defer store.Close()

	webMux := http.NewServeMux()
	webMux.HandleFunc("GET /todos", api.ListTodos)
	webMux.HandleFunc("POST /todos/create", api.CreateTodo)
	webMux.HandleFunc("PUT /todos/{uuid}", api.UpdateTodo)
	webMux.HandleFunc("PUT /todos/{uuid}/complete", api.ToggleComplete)
	webMux.HandleFunc("DELETE /todos/{uuid}", api.DeleteTodo)
	webMux.HandleFunc("DELETE /todos/clear", api.ClearCompletedTodos)

	// API Backend
	apiKey := os.Getenv("TODO_API_KEY")
	webMux.Handle("POST /api/todos", api.RequireBearerToken(apiKey, http.HandlerFunc(api.CreateTodo)))

	// Serve Static Frontend
	webMux.Handle("/", http.FileServer(http.Dir("./frontend/dist")))

	// Start web server on 8080
	log.Fatal(http.ListenAndServe("[::]:8080", webMux))
}
