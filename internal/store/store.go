package store

import (
	"database/sql"
	"errors"
	"os"
	"todo/internal/todo"

	_ "modernc.org/sqlite"
)

const TodoStoreFile = "store/todos.db"

var db *sql.DB

// Initialize sets up the database connection and creates tables
func Initialize() error {
	if err := os.MkdirAll("store", os.ModePerm); err != nil {
		return err
	}

	var err error
	db, err = sql.Open("sqlite", TodoStoreFile)
	if err != nil {
		return err
	}

	// Create todos table if it doesn't exist
	schema := `
	CREATE TABLE IF NOT EXISTS todos (
		uuid TEXT PRIMARY KEY,
		created_at DATETIME NOT NULL,
		todo TEXT NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_completed ON todos(completed);
	`
	_, err = db.Exec(schema)
	return err
}

// Close closes the database connection
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// List returns all todos.
func List() ([]todo.Todo, error) {
	rows, err := db.Query("SELECT uuid, created_at, todo, completed FROM todos ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []todo.Todo
	for rows.Next() {
		var t todo.Todo
		if err := rows.Scan(&t.UUID, &t.CreatedAt, &t.ToDo, &t.Completed); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, rows.Err()
}

// Create adds a new todo.
func Create(newTodo todo.Todo) (*todo.Todo, error) {
	_, err := db.Exec(
		"INSERT INTO todos (uuid, created_at, todo, completed) VALUES (?, ?, ?, ?)",
		newTodo.UUID, newTodo.CreatedAt, newTodo.ToDo, newTodo.Completed,
	)
	if err != nil {
		return nil, err
	}
	return &newTodo, nil
}

// ToggleComplete toggles the completion status of the todo with the given UUID.
func ToggleComplete(uuid todo.TodoUUID) error {
	result, err := db.Exec(
		"UPDATE todos SET completed = NOT completed WHERE uuid = ?",
		uuid,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("todo not found")
	}
	return nil
}

// Update updates the todo with the given UUID.
func Update(uuid todo.TodoUUID, updatedTodo string) error {
	result, err := db.Exec(
		"UPDATE todos SET todo = ? WHERE uuid = ?",
		updatedTodo, uuid,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("todo not found")
	}
	return nil
}

// Delete removes a todo with the given UUID.
func Delete(uuid todo.TodoUUID) error {
	result, err := db.Exec("DELETE FROM todos WHERE uuid = ?", uuid)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("todo not found")
	}
	return nil
}

// ClearCompleted deletes all completed todos.
func ClearCompleted() error {
	_, err := db.Exec("DELETE FROM todos WHERE completed = 1")
	return err
}
