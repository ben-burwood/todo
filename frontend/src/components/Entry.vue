<template>
    <div class="join">
        <input type="text" class="input w-full join-item" placeholder="New todo..." v-model="title" @keyup.enter="addTodo" />
        <button class="btn btn-accent join-item" @click="addTodo" :disabled="title.trim() === ''">Add</button>
    </div>
</template>

<script setup lang="ts">
import { ref } from "vue";

const emit = defineEmits(["add"]);

const title = ref("");

function addTodo() {
    if (title.value.trim() === "") return;
    emit("add", title.value);
    title.value = "";
}
</script>
