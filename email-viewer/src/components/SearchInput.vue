<template>
    <div class="relative">
        <label for="search" class="sr-only">Search Emails</label>
        <Magnifier class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
        <input v-model="localQuery" @input="handleInput" id="search" type="text" name="search"
            class="pl-10 pr-3 py-4 w-[600px] border border-gray-300 rounded-md leading-5 bg-white placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
            placeholder="Search emails...">
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import Magnifier from './icons/Magnifier.vue';

const localQuery = ref('');
const emit = defineEmits(['query']);

let timeoutId: ReturnType<typeof setTimeout> | null = null;


const handleInput = () => {
    if (timeoutId) clearTimeout(timeoutId);
    timeoutId = setTimeout(() => {
        emit('query', localQuery.value);
    }, 300);
};
</script>

<style scoped>
 input::placeholder {
    @apply italic text-gray-400;
 }
</style>