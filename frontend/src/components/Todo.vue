<template>
    <div class="border border-base-300/50 shadow-sm rounded-lg bg-base-100 hover:bg-primary/10">
        <div class="flex flex-row gap-4 p-2 items-center">
            <input
                type="checkbox"
                class="checkbox checkbox-primary"
                :checked="completed"
                @click="emit('completed', !completed)"
                :class="{ 'checkbox-primary': !completed, 'checkbox-primary/50': completed }"
            />
            <span class="overflow-x-auto" :style="{ textDecoration: completed ? 'line-through' : 'none' }" :class="{ 'text-gray-500': completed }">{{
                title
            }}</span>
        </div>
    </div>
</template>

<script setup lang="ts">
const props = defineProps<{
    title: string;
    completed: boolean;
}>();

const emit = defineEmits(["completed"]);
</script>
