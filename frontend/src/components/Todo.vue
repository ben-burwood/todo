<template>
    <div class="card bg-base-100 hover:bg-primary/10 shadow-sm">
        <div class="p-2">
            <div class="flex flex-row gap-5">
                <input
                    type="checkbox"
                    class="checkbox checkbox-primary"
                    :checked="completed"
                    @click="emit('completed', !completed)"
                    :class="{ 'checkbox-primary': !completed, 'checkbox-primary/50': completed }"
                />
                <label :style="{ textDecoration: completed ? 'line-through' : 'none' }" :class="{ 'text-gray-500': completed }">{{ title }}</label>
            </div>
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
