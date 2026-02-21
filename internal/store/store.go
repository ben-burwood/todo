package store

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"todo/internal/todo"
)

const TodoStoreFile = "store/todos.json"

var mu sync.Mutex

// loadTodos reads the todos from the JSON file
func loadTodos() ([]todo.Todo, error) {
	file, err := os.Open(TodoStoreFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []todo.Todo{}, nil // treat as empty list if file doesn't exist
		}
		return nil, err
	}
	defer file.Close()

	var todos []todo.Todo
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&todos); err != nil && err.Error() != "EOF" {
		return nil, err
	}
	return todos, nil
}

// saveTodos writes the todos to the JSON file
func saveTodos(todos []todo.Todo) error {
	if err := os.MkdirAll("store", os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(TodoStoreFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(todos)
}

// List returns all todos.
func List() ([]todo.Todo, error) {
	mu.Lock()
	defer mu.Unlock()
	return loadTodos()
}

// Create adds a new todo with the given title.
func Create(newTodo todo.Todo) (*todo.Todo, error) {
	mu.Lock()
	defer mu.Unlock()

	todos, err := loadTodos()
	if err != nil {
		return nil, err
	}

	todos = append(todos, newTodo)

	if err := saveTodos(todos); err != nil {
		return nil, err
	}
	return &newTodo, nil
}

// ToggleComplete toggles the completion status of the todo with the given UUID.
func ToggleComplete(uuid todo.TodoUUID) error {
	mu.Lock()
	defer mu.Unlock()

	todos, err := loadTodos()
	if err != nil {
		return err
	}

	found := false
	for i, t := range todos {
		if t.UUID == uuid {
			todos[i].Completed = !todos[i].Completed
			found = true
			break
		}
	}
	if !found {
		return errors.New("todo not found")
	}
	return saveTodos(todos)
}

// Update updates the title of the todo with the given UUID.
func Update(uuid todo.TodoUUID, newTitle string) error {
	mu.Lock()
	defer mu.Unlock()

	todos, err := loadTodos()
	if err != nil {
		return err
	}

	found := false
	for i, t := range todos {
		if t.UUID == uuid {
			todos[i].Title = newTitle
			found = true
			break
		}
	}
	if !found {
		return errors.New("todo not found")
	}
	return saveTodos(todos)
}

// Delete removes a todo with the given UUID.
func Delete(uuid todo.TodoUUID) error {
	mu.Lock()
	defer mu.Unlock()

	todos, err := loadTodos()
	if err != nil {
		return err
	}

	found := false
	newTodos := make([]todo.Todo, 0, len(todos))
	for _, t := range todos {
		if t.UUID == uuid {
			found = true
			continue
		}
		newTodos = append(newTodos, t)
	}
	if !found {
		return errors.New("todo not found")
	}
	return saveTodos(todos)
}

// ClearCompleted deletes all completed todos.
func ClearCompleted() error {
	mu.Lock()
	defer mu.Unlock()

	todos, err := loadTodos()
	if err != nil {
		return err
	}

	newTodos := make([]todo.Todo, 0, len(todos))
	for _, t := range todos {
		if !t.Completed {
			newTodos = append(newTodos, t)
		}
	}
	return saveTodos(newTodos)
}
