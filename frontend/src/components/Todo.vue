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
            <span 
                class="overflow-x-auto" 
                :style="{ textDecoration: completed ? 'line-through' : 'none' }" 
                :class="{ 'text-gray-500': completed }"
                v-html="linkifiedTitle"
            ></span>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed } from "vue";

const props = defineProps<{
    title: string;
    completed: boolean;
}>();

const emit = defineEmits(["completed"]);

// Function to convert URLs in text to clickable links
function linkifyText(text: string): string {
    // Escape HTML entities in the input text first
    const escapeHtml = (str: string): string => {
        const div = document.createElement('div');
        div.textContent = str;
        return div.innerHTML;
    };
    
    // Regular expression to match URLs
    const urlRegex = /(https?:\/\/[^\s]+)/g;
    
    // Escape the entire text first
    const escapedText = escapeHtml(text);
    
    // Replace URLs with anchor tags
    return escapedText.replace(urlRegex, (url) => {
        // Remove trailing punctuation that might not be part of the URL
        let cleanUrl = url;
        const trailingPunctuation = /[.,;!?)}\]]+$/;
        const match = url.match(trailingPunctuation);
        let trailing = '';
        
        if (match) {
            trailing = match[0];
            cleanUrl = url.slice(0, -trailing.length);
        }
        
        // Validate that the URL starts with http:// or https://
        // This prevents javascript: or data: URLs
        if (!cleanUrl.match(/^https?:\/\//i)) {
            return url; // Return original text if not a valid http(s) URL
        }
        
        return `<a href="${cleanUrl}" target="_blank" rel="noopener noreferrer" class="text-primary hover:underline" onclick="event.stopPropagation()">${cleanUrl}</a>${trailing}`;
    });
}

const linkifiedTitle = computed(() => linkifyText(props.title));
</script>
