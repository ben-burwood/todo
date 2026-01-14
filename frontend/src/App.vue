<template>
    <div class="min-h-screen min-w-screen bg-base-200">
        <ThemeSwitcher class="fixed top-4 right-4" />

        <div class="flex flex-col items-center justify-center p-5 w-lg mx-auto">
            <h1 class="text-4xl font-bold">Todo</h1>

            <Entry @add="addTodo" class="m-5 w-full" />

            <div class="m-5 flex flex-col gap-2">
                <span v-if="outstandingTodos.length === 0" class="text-4xl">🎉</span>
                <Todo
                    class="w-md"
                    v-else
                    v-for="todo in outstandingTodos"
                    :key="todo.uuid"
                    :title="todo.title"
                    :completed="todo.completed"
                    @completed="toggleCompleted(todo.uuid)"
                />
            </div>
            <div v-if="completedTodos.length > 0">
                <div class="divider"></div>
                <div class="flex flex-row justify-between">
                    <span class="collapse-title text-gray-500">Completed ({{ completedTodos.length }})</span>
                    <button class="btn btn-sm btn-outline btn-error" @click="clearCompleted">Clear</button>
                </div>
                <div class="flex flex-col gap-2">
                    <Todo
                        class="w-md"
                        v-for="todo in completedTodos"
                        :key="todo.uuid"
                        :title="todo.title"
                        :completed="todo.completed"
                        @completed="toggleCompleted(todo.uuid)"
                    />
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import Entry from "@/components/Entry.vue";
import Todo from "@/components/Todo.vue";
import ThemeSwitcher from "@/components/ThemeSwitcher.vue";
import { SERVER_URL } from "@/main";

const todos = ref<{ uuid: string; title: string; completed: boolean }[]>([]);
const outstandingTodos = computed(() => todos.value.filter((todo) => !todo.completed));
const completedTodos = computed(() => todos.value.filter((todo) => todo.completed));

async function fetchTodos() {
    const res = await fetch(`${SERVER_URL}/todos`);
    todos.value = await res.json();
}
onMounted(fetchTodos);

async function addTodo(title: string) {
    const res = await fetch(`${SERVER_URL}/todos/create`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title }),
    });
    if (res.ok) await fetchTodos();
}

async function toggleCompleted(uuid: string) {
    const res = await fetch(`${SERVER_URL}/todos/${uuid}/complete`, { method: "PUT" });
    if (res.ok) await fetchTodos();
}

async function clearCompleted() {
    const res = await fetch(`${SERVER_URL}/todos/clear`, { method: "DELETE" });
    if (res.ok) await fetchTodos();
}
</script>
