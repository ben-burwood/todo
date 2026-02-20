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
    // Escape HTML entities using a string-based approach
    const escapeHtml = (str: string): string => {
        return str
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/"/g, '&quot;')
            .replace(/'/g, '&#039;');
    };
    
    // Regular expression to match URLs
    const urlRegex = /(https?:\/\/[^\s]+)/g;
    
    // Find all URLs and their positions before escaping
    const urlMatches: { url: string; start: number; end: number }[] = [];
    let match;
    while ((match = urlRegex.exec(text)) !== null) {
        urlMatches.push({
            url: match[0],
            start: match.index,
            end: match.index + match[0].length
        });
    }
    
    // Build the result by processing text segments
    let result = '';
    let lastIndex = 0;
    
    for (const urlMatch of urlMatches) {
        // Add escaped text before the URL
        result += escapeHtml(text.substring(lastIndex, urlMatch.start));
        
        // Process the URL
        let url = urlMatch.url;
        let cleanUrl = url;
        const trailingPunctuation = /[.,;!?)\}\]]+$/;
        const punctMatch = url.match(trailingPunctuation);
        let trailing = '';
        
        if (punctMatch) {
            trailing = punctMatch[0];
            cleanUrl = url.slice(0, -trailing.length);
        }
        
        // Validate that the URL starts with http:// or https://
        if (!cleanUrl.match(/^https?:\/\//i)) {
            result += escapeHtml(url);
        } else {
            // Encode the URL for safe insertion in href attribute
            const encodedUrl = encodeURI(cleanUrl);
            result += `<a href="${encodedUrl}" target="_blank" rel="noopener noreferrer" class="text-primary hover:underline" onclick="event.stopPropagation()">${escapeHtml(cleanUrl)}</a>${escapeHtml(trailing)}`;
        }
        
        lastIndex = urlMatch.end;
    }
    
    // Add any remaining text after the last URL
    result += escapeHtml(text.substring(lastIndex));
    
    return result;
}

const linkifiedTitle = computed(() => linkifyText(props.title));
</script>
