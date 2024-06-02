import { computed } from 'vue';
<template>
    <button class="px-4 py-4 sm:px-6 w-full" @click="handleClick">
        <div class="w-full flex items-center justify-between">
            <p class="text-sm font-semibold text-indigo-600 truncate">{{ email.Subject }}</p>
            <div class="ml-2 flex-shrink-0 flex">
                <p class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-blue-100 text-blue-700">
                    {{ emailDate.date }}
                </p>
            </div>
        </div>
        <div class="mt-2 flex flex-col">
            
            <div class="flex gap-1 items-center text-sm text-gray-500">
                <UserIcon class="size-5"/>
                {{ email.From }}
            </div>
            
            <div class="flex gap- items-center text-sm text-gray-500">
                <EnvelopeIcon class="h-5 w-5" />
                <div class="truncate">
                    {{ email.To?.join(', ') }}
                </div>
            </div>
        </div>
    </button>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import dayjs from 'dayjs';
import type { Email } from '@/types';
import EnvelopeIcon from './icons/Envelope.vue';
import UserIcon from './icons/User.vue';

interface Props {
    email: Email
};

const props = defineProps<Props>();

// const emit = defineEmits(['selected']);
const emit = defineEmits<{
    selected: [value: Email]
}>();

const handleClick = () => {
    emit('selected', props.email);
}

const emailDate = computed(() => {
    const date = dayjs(props.email.Date);
    return {
        date: date.format("YYYY-MMM-DD"),
        time: date.format("HH:mm")
    };
});

</script>