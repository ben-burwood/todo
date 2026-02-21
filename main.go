package main

import (
	"net/http"
	"todo/internal/api"
	"todo/internal/store"
)

func main() {
	// Initialize the database
	if err := store.Initialize(); err != nil {
		panic(err)
	}
	defer store.Close()

	webMux := http.NewServeMux()
	webMux.HandleFunc("GET /todos", api.ListTodos)
	webMux.HandleFunc("POST /todos/create", api.CreateTodo)
	webMux.HandleFunc("PUT /todos/{uuid}", api.UpdateTodo)
	webMux.HandleFunc("PUT /todos/{uuid}/complete", api.ToggleComplete)
	webMux.HandleFunc("DELETE /todos/{uuid}", api.DeleteTodo)
	webMux.HandleFunc("DELETE /todos/clear", api.ClearCompletedTodos)
	// Serve Static Frontend
	webMux.Handle("/", http.FileServer(http.Dir("./frontend/dist")))

	// Start web server on 8080
	http.ListenAndServe("[::]:8080", api.CORSMiddleware(webMux))
}
