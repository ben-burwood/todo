<template>
    <div class="join">
        <input type="text" class="input w-full join-item" placeholder="New todo..." v-model="newTodo" @keyup.enter="addTodo" />
        <button class="btn btn-accent join-item" @click="addTodo" :disabled="newTodo.trim() === ''">Add</button>
    </div>
</template>

<script setup lang="ts">
import { ref } from "vue";

const emit = defineEmits(["add"]);

const newTodo = ref("");

function addTodo() {
    const content = newTodo.value.trim();
    if (!content) return;

    emit("add", content);
    todo.value = "";
}
</script>
