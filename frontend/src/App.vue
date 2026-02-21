<template>
    <div class="min-h-screen min-w-screen bg-base-200">
        <ThemeSwitcher class="absolute top-4 right-4 hidden md:block" />

        <div class="flex flex-col items-center justify-center p-5 w-full max-w-lg mx-auto">
            <h1 class="text-4xl font-bold hidden md:block">ToDo</h1>

            <Entry @add="addTodo" class="mt-5 w-full" />
            <div class="divider"></div>

            <div class="flex flex-col gap-2 w-full">
                <span v-if="outstandingTodos.length === 0" class="text-4xl text-center">🎉</span>
                <Todo
                    v-else
                    v-for="todo in outstandingTodos"
                    :key="todo.uuid"
                    :title="todo.title"
                    :completed="todo.completed"
                    @completed="toggleCompleted(todo.uuid)"
                    @edit="(newTitle) => updateTodo(todo.uuid, newTitle)"
                    @delete="deleteTodo(todo.uuid)"
                />
                <div class="mt-10" v-if="completedTodos.length > 0">
                    <div class="flex flex-row justify-between items-center m-2 mb-4">
                        <span class="text-gray-500">Completed ({{ completedTodos.length }})</span>
                        <button class="btn btn-sm btn-outline btn-error" @click="clearCompleted">Clear</button>
                    </div>
                    <div class="flex flex-col gap-2">
                        <Todo
                            v-for="todo in completedTodos"
                            :key="todo.uuid"
                            :title="todo.title"
                            :completed="todo.completed"
                            @completed="toggleCompleted(todo.uuid)"
                            @edit="(newTitle) => updateTodo(todo.uuid, newTitle)"
                            @delete="deleteTodo(todo.uuid)"
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

const errorMessage = ref("");
// clear error after 5 seconds
watch(errorMessage, (newError) => {
    if (newError) {
        setTimeout(() => {
            errorMessage.value = "";
        }, 5000);
    }
});

const todos = ref<{ uuid: string; title: string; completed: boolean; created_at: string }[]>([]);
const outstandingTodos = computed(() =>
    todos.value.filter((todo) => !todo.completed).sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()),
);

const completedTodos = computed(() =>
    todos.value.filter((todo) => todo.completed).sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()),
);

async function fetchTodos() {
    try {
        const res = await fetch(`${SERVER_URL}/todos`);
        todos.value = await res.json();
    } catch (error: any) {
        errorMessage.value = `Error: Fetching Todos : ${error.message}`;
    }
}
onMounted(fetchTodos);

async function addTodo(title: string) {
    try {
        const res = await fetch(`${SERVER_URL}/todos/create`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ title }),
        });
        await fetchTodos();
    } catch (error: any) {
        errorMessage.value = `Error: Adding Todo : ${error.message}`;
    }
}

async function toggleCompleted(uuid: string) {
    try {
        const res = await fetch(`${SERVER_URL}/todos/${uuid}/complete`, { method: "PUT" });
        await fetchTodos();
    } catch (error: any) {
        errorMessage.value = `Error: Toggling Complete : ${error.message}`;
    }
}

async function updateTodo(uuid: string, newTitle: string) {
    try {
        const res = await fetch(`${SERVER_URL}/todos/${uuid}`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ title: newTitle }),
        });
        await fetchTodos();
    } catch (error) {
        errorMessage.value = `Error: Updating Todo : ${error.message}`;
    }
}

async function deleteTodo(uuid: string) {
    if (!confirm("Are you sure you want to delete this todo?")) {
        return;
    }
    try {
        const res = await fetch(`${SERVER_URL}/todos/${uuid}`, { method: "DELETE" });
        await fetchTodos();
    } catch (error: any) {
        errorMessage.value = `Error: Deleting Todo : ${error.message}`;
    }
}

async function clearCompleted() {
    try {
        const res = await fetch(`${SERVER_URL}/todos/clear`, { method: "DELETE" });
        await fetchTodos();
    } catch (error: any) {
        errorMessage.value = `Error: Clearing Completed : ${error.message}`;
    }
}
</script>
