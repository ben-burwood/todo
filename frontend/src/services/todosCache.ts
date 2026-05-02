import type { Todo } from "@/types/todo";
import { readJSONArray, writeJSON } from "@/services/storage";

const KEY = "cached_todos";

export function save(todos: Todo[]): void {
    writeJSON(KEY, todos);
}

export function load(): Todo[] {
    return readJSONArray<Todo>(KEY);
}
