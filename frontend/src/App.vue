<template>
    <div class="min-h-screen min-w-screen bg-base-200">
        <div class="absolute top-4 right-4 hidden md:flex items-center gap-2">
            <span v-if="!online" class="badge badge-warning">Offline</span>
            <span v-else-if="pendingCount > 0" class="badge badge-info">Syncing {{ pendingCount }}…</span>
            <ThemeSwitcher />
        </div>

        <div class="flex flex-col items-center justify-center p-5 w-full max-w-lg mx-auto">
            <h1 class="text-4xl font-bold hidden md:block">ToDo</h1>

            <Entry @add="addTodo" class="mt-5 w-full" />
            <div class="divider"></div>

            <div class="flex flex-col gap-2 w-full">
                <span v-if="outstandingTodos.length === 0" class="text-4xl text-center">🎉</span>
                <div v-else v-for="todo in outstandingTodos" :key="todo.uuid" class="flex flex-col">
                    <Todo
                        :todo="todo.todo"
                        :completed="todo.completed"
                        :class="{ 'opacity-60': todo.pending }"
                        @completed="toggleCompleted(todo.uuid)"
                        @edit="(updatedTodo) => updateTodo(todo.uuid, updatedTodo)"
                        @delete="deleteTodo(todo.uuid)"
                    />
                    <span v-if="todo.pending" class="text-xs text-gray-500 ml-2 mt-1">queued — will sync when online</span>
                </div>
                <div class="mt-10" v-if="completedTodos.length > 0">
                    <div class="flex flex-row justify-between items-center m-2 mb-4">
                        <span class="text-gray-500">Completed ({{ completedTodos.length }})</span>
                        <button class="btn btn-sm btn-outline btn-error" @click="clearCompleted">Clear</button>
                    </div>
                    <div class="flex flex-col gap-2">
                        <Todo
                            v-for="todo in completedTodos"
                            :key="todo.uuid"
                            :todo="todo.todo"
                            :completed="todo.completed"
                            @completed="toggleCompleted(todo.uuid)"
                        />
                    </div>
                </div>
            </div>

            <div v-if="errorMessage" role="alert" class="alert alert-error fixed bottom-10">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 shrink-0 stroke-current" fill="none" viewBox="0 0 24 24">
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                </svg>
                {{ errorMessage }}
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import Entry from "@/components/Entry.vue";
import Todo from "@/components/Todo.vue";
import ThemeSwitcher from "@/components/ThemeSwitcher.vue";
import { SERVER_URL } from "@/main";
import * as offlineQueue from "@/services/offlineQueue";
import * as todosCache from "@/services/todosCache";
import { useOnline } from "@/composables/useOnline";
import type { Todo as TodoType, TodoItem } from "@/types/todo";

const errorMessage = ref("");
watch(errorMessage, (newError) => {
    if (newError) {
        setTimeout(() => {
            errorMessage.value = "";
        }, 5000);
    }
});

const { online } = useOnline();
const todos = ref<TodoItem[]>([]);

const outstandingTodos = computed(() => todos.value.filter((todo) => !todo.completed));

const completedTodos = computed(() => todos.value.filter((todo) => todo.completed));

const pendingCount = computed(() => todos.value.filter((t) => t.pending).length);

function mergePending(serverList: TodoType[]): TodoItem[] {
    const queued = offlineQueue.list();
    const seen = new Set(serverList.map((t) => t.uuid));
    const stillPending = queued
        .filter((q) => !seen.has(q.uuid))
        .map<TodoItem>((q) => ({ uuid: q.uuid, todo: q.todo, completed: false, created_at: q.created_at, pending: true }));
    // Pending items are the freshest entries (just typed by the user). Server list
    // arrives sorted by created_at DESC, so prepending pending preserves DESC order.
    return [...stillPending, ...serverList];
}

async function fetchTodos() {
    try {
        const res = await fetch(`${SERVER_URL}/todos`);
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        const fresh: TodoType[] = await res.json();
        todosCache.save(fresh);
        todos.value = mergePending(fresh);
    } catch (error: any) {
        if (!online.value) {
            // Offline — cached list (already in todos.value) is the best we have.
            return;
        }
        errorMessage.value = `Error: Fetching Todos : ${error.message}`;
    }
}

async function flushQueue() {
    if (!online.value) return;
    const outcome = await offlineQueue.flush(SERVER_URL);
    if (outcome.rejected.length > 0) {
        errorMessage.value = `Sync rejected ${outcome.rejected.length} todo(s)`;
    }
}

async function syncWithServer() {
    await flushQueue();
    await fetchTodos();
}

async function addTodo(content: string) {
    const newTodo: TodoItem = {
        uuid: crypto.randomUUID(),
        todo: content,
        completed: false,
        created_at: new Date().toISOString(),
        pending: true,
    };
    todos.value.unshift(newTodo);
    offlineQueue.enqueue({ uuid: newTodo.uuid, todo: newTodo.todo, created_at: newTodo.created_at });

    if (online.value) {
        await syncWithServer();
    }
}

async function toggleCompleted(uuid: string) {
    try {
        const res = await fetch(`${SERVER_URL}/todos/${uuid}/complete`, { method: "PUT" });
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        await fetchTodos();
    } catch (error: any) {
        errorMessage.value = `Error: Toggling Complete : ${error.message}`;
    }
}

async function updateTodo(uuid: string, updatedTodo: string) {
    try {
        const res = await fetch(`${SERVER_URL}/todos/${uuid}`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ todo: updatedTodo }),
        });
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        await fetchTodos();
    } catch (error: any) {
        errorMessage.value = `Error: Updating Todo : ${error.message}`;
    }
}

async function deleteTodo(uuid: string) {
    if (!confirm("Are you sure you want to delete this todo?")) {
        return;
    }
    try {
        const res = await fetch(`${SERVER_URL}/todos/${uuid}`, { method: "DELETE" });
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        await fetchTodos();
    } catch (error: any) {
        errorMessage.value = `Error: Deleting Todo : ${error.message}`;
    }
}

async function clearCompleted() {
    try {
        const res = await fetch(`${SERVER_URL}/todos/clear`, { method: "DELETE" });
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        await fetchTodos();
    } catch (error: any) {
        errorMessage.value = `Error: Clearing Completed : ${error.message}`;
    }
}

watch(online, (isOnline) => {
    if (isOnline) syncWithServer();
});

onMounted(async () => {
    todos.value = mergePending(todosCache.load());
    if (online.value) {
        await syncWithServer();
    }
});
</script>
