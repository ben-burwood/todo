export interface Todo {
    uuid: string;
    todo: string;
    completed: boolean;
    created_at: string;
}

export type PendingTodo = Pick<Todo, "uuid" | "todo" | "created_at">;

export type TodoItem = Todo & { pending?: boolean };
