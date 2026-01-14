package todo

import (
	"time"

	"github.com/google/uuid"
)

type TodoUUID string

type Todo struct {
	UUID      TodoUUID  `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
}

// NewTodo - Constructor to create a new Todo with a unique UUID and current time as created_at.
func NewTodo(title string) *Todo {
	return &Todo{
		UUID:      TodoUUID(uuid.NewString()),
		CreatedAt: time.Now(),
		Title:     title,
		Completed: false,
	}
}
