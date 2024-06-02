<template>
    <div class="container mx-auto h-screen flex flex-col">
        
        <div class="flex gap-2 items-center justify-between bg-slate-100 p-4">
            <div class="flex items-center gap-2"> 
                <AtIcon class="bg-slate-400 text-white rounded-full size-8 " /> 
                <h3 class="text-xl">Elron Emails Viewer</h3>
            </div>
            <SearchInput @query="fetchEmails" />
        </div>
        <div class="bg-white flex shadow overflow-hidden sm:rounded-md divide-x">

            <ul class="divide-y divide-gray-200 max-w-lg overflow-y-auto">
                <li v-for="email in emails" :key="email.Id" 
                    :class="email.Id == selectedId ? 'bg-gray-200' : 'hover:bg-gray-100'" >
                    <EmailPreviewCard :email="email" @selected="setSelected" />
                </li>
            </ul>

            <div class="flex-1 overflow-y-auto">
                <EmailView :email="selectedEmail" :search="search" />
            </div>
        </div>
    
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import AtIcon from './components/icons/At.vue';
import SearchInput from './components/SearchInput.vue';
import EmailPreviewCard from './components/EmailPreviewCard.vue';
import EmailView from './components/EmailView.vue';
import { queryEmails } from './api/search';
import type { Email } from './types';

const emails = ref<Email[]>([]);
const selectedEmail = ref<Email | null>(null);
const selectedId = ref<string | null >(null);
const search = ref<string>('');

const setSelected = (email: Email) => {
    selectedEmail.value = email;
    selectedId.value = email.Id;
};

const fetchEmails = async (searchQuery: string) => {
    emails.value = await queryEmails(searchQuery);
    search.value = searchQuery;
    selectedEmail.value = null;
    selectedEmail.value = null;
};

</script>