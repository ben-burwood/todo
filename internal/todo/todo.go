package todo

import (
	"time"

	"github.com/google/uuid"
)

type TodoUUID string

type Todo struct {
	UUID      TodoUUID  `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	ToDo      string    `json:"todo"`
	Completed bool      `json:"completed"`
}

// NewTodo - Constructor to create a new Todo with a unique UUID and current time as created_at.
func NewTodo(todo string) *Todo {
	return &Todo{
		UUID:      TodoUUID(uuid.NewString()),
		CreatedAt: time.Now(),
		ToDo:      todo,
		Completed: false,
	}
}
