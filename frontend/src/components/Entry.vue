<template>
    <div class="join">
        <input type="text" class="input w-full join-item" placeholder="New todo..." v-model="todo" @keyup.enter="addTodo" />
        <button class="btn btn-accent join-item" @click="addTodo" :disabled="todo.trim() === ''">Add</button>
    </div>
</template>

<script setup lang="ts">
import { ref } from "vue";

const emit = defineEmits(["add"]);

const todo = ref("");

function addTodo() {
    if (todo.value.trim() === "") return;
    emit("add", todo.value);
    todo.value = "";
}
</script>
