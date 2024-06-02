import { computed } from 'vue';
<template>
    <div v-if="email" class="p-2">
        <div class="flex flex-col gap-3">
            <!-- header -->
            <div class="flex flex-col gap-1 bg-blue-100 rounded p-3">
                <div class="flex justify-between">
                    <div>
                        <h1 class="font-bold text-lg" v-html="emailSubjectHighlighted"/>
                        <div class="text-gray-500">
                            <span class="font-bold">
                                From:
                            </span>
                            {{ email.From }}
                        </div>
                    </div>
                    <div class="text-gray-500 text-right min-w-36">
                        <div> {{ emailDate.date }}</div>
                        <div> {{ emailDate.time }}</div>
                    </div>
                </div>
                <div class="text-gray-500">
                    <span class="font-bold">
                        To:
                    </span>
                    <span class="italic">
                        {{ email.To?.join(', ') }}
                    </span>
                </div>
                <div class="text-gray-500" v-if="email.Cc">
                    <span class="font-bold">
                        CC:
                    </span>
                    <span class="italic">
                        {{ email.Cc?.join(', ') }}
                    </span>
                </div>
                
            </div>
            <hr />
            <!-- email body -->
            <div class="text-gray-600">
                <div class="whitespace-pre-line p-2 rounded-md" v-html="emailBodyHighlighted"/>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import dayjs from 'dayjs';
import type { Email } from '../types';

interface Props {
    email: Email | null,
    search : string
};

const props = defineProps<Props>();
const emailDate = computed(() => {
    const date = dayjs(props.email?.Date);
    return {
        date: date.format("YYYY-MMM-DD"),
        time: date.format("HH:mm A")
    };
});

const emailBodyHighlighted = computed(() => {
    if (props.email?.Body){
        const regex = new RegExp(props.search, 'gi');
        return props.email?.Body.replace(regex, `<span class="highlight">${props.search}</span>`)
    }
    return '';
})

const emailSubjectHighlighted = computed(() => {
    if (props.email?.Body){
        const regex = new RegExp(props.search, 'gi');
        return props.email?.Subject.replace(regex, `<span class="highlight">${props.search}</span>`)
    }
    return '';
})

</script>

<style scoped>
    
</style>