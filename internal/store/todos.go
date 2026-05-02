package store

import (
	"errors"
	"time"
	"todo/internal/todo"
)

var ErrNotFound = errors.New("todo not found")

func scanTodo(scanner interface{ Scan(dest ...any) error }) (todo.Todo, error) {
	var t todo.Todo
	var createdAt string
	if err := scanner.Scan(&t.UUID, &t.ToDo, &t.Completed, &createdAt); err != nil {
		return todo.Todo{}, err
	}
	parsed, err := time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return todo.Todo{}, err
	}
	t.CreatedAt = parsed
	return t, nil
}

func execAffectingOne(query string, args ...any) error {
	res, err := db.Exec(query, args...)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// List returns all todos.
func List() ([]todo.Todo, error) {
	rows, err := db.Query(`SELECT uuid, todo, completed, created_at FROM todos`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []todo.Todo{}
	for rows.Next() {
		t, err := scanTodo(rows)
		if err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}

// Create inserts a todo, or returns the existing one if a row with the same
// UUID already exists, so retried syncs from offline clients are idempotent.
func Create(newTodo todo.Todo) (*todo.Todo, error) {
	row := db.QueryRow(
		`INSERT INTO todos (uuid, todo, completed, created_at) VALUES (?, ?, ?, ?)
         ON CONFLICT(uuid) DO UPDATE SET uuid = excluded.uuid
         RETURNING uuid, todo, completed, created_at`,
		newTodo.UUID,
		newTodo.ToDo,
		newTodo.Completed,
		newTodo.CreatedAt.Format(time.RFC3339Nano),
	)
	t, err := scanTodo(row)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func Update(uuid todo.TodoUUID, updatedTodo string) error {
	return execAffectingOne(`UPDATE todos SET todo = ? WHERE uuid = ?`, updatedTodo, uuid)
}

func ToggleComplete(uuid todo.TodoUUID) error {
	return execAffectingOne(`UPDATE todos SET completed = 1 - completed WHERE uuid = ?`, uuid)
}

func Delete(uuid todo.TodoUUID) error {
	return execAffectingOne(`DELETE FROM todos WHERE uuid = ?`, uuid)
}

func ClearCompleted() error {
	_, err := db.Exec(`DELETE FROM todos WHERE completed = 1`)
	return err
}
