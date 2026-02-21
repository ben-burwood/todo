<template>
    <div class="border border-base-300/50 shadow-sm rounded-lg bg-base-100 hover:bg-primary/10">
        <div class="flex flex-row gap-4 p-2 items-center">
            <input
                type="checkbox"
                class="checkbox checkbox-primary"
                :checked="completed"
                @click="emit('completed', !completed)"
                :class="{ 'checkbox-primary': !completed, 'checkbox-primary/50': completed }"
                :disabled="isEditing"
            />
            <span
                v-if="!isEditing"
                class="overflow-x-auto flex-1"
                :style="{ textDecoration: completed ? 'line-through' : 'none' }"
                :class="{ 'text-gray-500': completed }"
                v-html="linkifiedTitle"
            ></span>
            <input
                v-else
                type="text"
                class="input input-sm input-bordered flex-1"
                v-model="editedTitle"
                @keyup.enter="saveEdit"
                @keyup.esc="cancelEdit"
                ref="editInput"
            />
            <div class="flex gap-2">
                <button v-if="!isEditing" class="btn btn-sm btn-ghost" @click="startEdit" :disabled="completed">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                        />
                    </svg>
                </button>
                <template v-else>
                    <button class="btn btn-sm btn-success" @click="saveEdit" :disabled="editedTitle.trim() === ''">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                        </svg>
                    </button>
                    <button class="btn btn-sm btn-error" @click="cancelEdit">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </template>
                <button @click="emit('delete')" class="btn btn-sm btn-ghost btn-circle text-error hover:bg-error/10" title="Delete todo">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                </button>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, ref, nextTick } from "vue";

const props = defineProps<{
    title: string;
    completed: boolean;
}>();

const emit = defineEmits(["completed", "edit", "delete"]);

const isEditing = ref(false);
const editedTitle = ref("");
const editInput = ref<HTMLInputElement | null>(null);

function startEdit() {
    editedTitle.value = props.title;
    isEditing.value = true;
    nextTick(() => {
        editInput.value?.focus();
    });
}

function saveEdit() {
    if (editedTitle.value.trim() === "") return;
    emit("edit", editedTitle.value);
    isEditing.value = false;
}

function cancelEdit() {
    isEditing.value = false;
    editedTitle.value = "";
}

// Convert URLs to clickable links (anchor tags)
function linkifyText(text: string): string {
    const escapeHtml = (str: string): string => {
        return str.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;").replace(/'/g, "&#039;");
    };

    const urlRegex = /(https?:\/\/[^\s<>"]+)/g;

    const urlMatches: { url: string; start: number; end: number }[] = [];
    let match;
    while ((match = urlRegex.exec(text)) !== null) {
        urlMatches.push({
            url: match[0],
            start: match.index,
            end: match.index + match[0].length,
        });
    }

    let result = "";
    let lastIndex = 0;

    for (const urlMatch of urlMatches) {
        result += escapeHtml(text.substring(lastIndex, urlMatch.start));

        let url = urlMatch.url;
        let cleanUrl = url;
        let trailing = "";

        const trailingPunctuation = /[.,;!?)\}\]]+$/;
        const punctMatch = url.match(trailingPunctuation);
        if (punctMatch) {
            trailing = punctMatch[0];
            cleanUrl = url.slice(0, -trailing.length);
        }

        if (!cleanUrl.match(/^https?:\/\//i)) {
            result += escapeHtml(url);
        } else {
            const encodedUrl = encodeURI(cleanUrl);
            result += `<a href="${encodedUrl}" target="_blank" rel="noopener noreferrer" class="text-primary hover:underline" onclick="event.stopPropagation()">${escapeHtml(cleanUrl)}</a>${escapeHtml(trailing)}`;
        }

        lastIndex = urlMatch.end;
    }

    result += escapeHtml(text.substring(lastIndex));
    return result;
}

const linkifiedTitle = computed(() => linkifyText(props.title));
</script>
