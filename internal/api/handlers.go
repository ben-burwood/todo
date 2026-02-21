package api

import (
	"encoding/json"
	"net/http"
	"todo/internal/store"
	"todo/internal/todo"
)

// ListTodos handles GET /todos
func ListTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := store.List()
	if err != nil {
		http.Error(w, "Failed to load todos", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// CreateTodo handles POST /todos
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	newTodo, err := store.Create(*todo.NewTodo(req.Title))
	if err != nil {
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}

// UpdateTodo handles PUT /todos/{uuid}
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	uuidString := r.PathValue("uuid")
	if uuidString == "" {
		http.Error(w, "Missing UUID", http.StatusBadRequest)
		return
	}
	var req struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := store.Update(todo.TodoUUID(uuidString), req.Title); err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ToggleComplete handles PUT /todos/{uuid}/complete
func ToggleComplete(w http.ResponseWriter, r *http.Request) {
	uuidString := r.PathValue("uuid")
	if uuidString == "" {
		http.Error(w, "Missing UUID", http.StatusBadRequest)
		return
	}
	if err := store.ToggleComplete(todo.TodoUUID(uuidString)); err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ClearCompletedTodos handles DELETE /todos/completed
func ClearCompletedTodos(w http.ResponseWriter, r *http.Request) {
	if err := store.ClearCompleted(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
