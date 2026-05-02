import type { PendingTodo } from "@/types/todo";
import { readJSONArray, writeJSON } from "@/services/storage";

export type FlushOutcome = {
    synced: PendingTodo[];
    rejected: { todo: PendingTodo; status: number }[];
    remaining: number;
};

const KEY = "pending_todos";

export function list(): PendingTodo[] {
    return readJSONArray<PendingTodo>(KEY);
}

export function enqueue(todo: PendingTodo): void {
    const items = list();
    if (items.some((t) => t.uuid === todo.uuid)) return;
    items.push(todo);
    writeJSON(KEY, items);
}

export function remove(uuid: string): void {
    writeJSON(KEY, list().filter((t) => t.uuid !== uuid));
}

export async function flush(serverUrl: string): Promise<FlushOutcome> {
    const synced: PendingTodo[] = [];
    const rejected: FlushOutcome["rejected"] = [];

    for (const item of list()) {
        let res: Response;
        try {
            res = await fetch(`${serverUrl}/todos/create`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ todo: item.todo, uuid: item.uuid, created_at: item.created_at }),
            });
        } catch {
            // Network error — stop early, leave this and any later items queued.
            break;
        }

        if (res.ok) {
            remove(item.uuid);
            synced.push(item);
            continue;
        }

        if (res.status >= 400 && res.status < 500) {
            // Client error — won't succeed on retry. Drop it.
            remove(item.uuid);
            rejected.push({ todo: item, status: res.status });
            continue;
        }

        // 5xx — server is alive but failing. Stop, retry later.
        break;
    }

    return { synced, rejected, remaining: list().length };
}
