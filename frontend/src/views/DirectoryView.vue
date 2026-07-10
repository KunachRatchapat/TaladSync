<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'

import AppHeader from '@/components/AppHeader.vue'
import DirectoryRow from '@/components/DirectoryRow.vue'
import SearchBar from '@/components/SearchBar.vue'
import StateBlock from '@/components/StateBlock.vue'
import { useDirectory } from '@/composables/useDirectory'
import { todayLocal } from '@/utils/date'

const date = todayLocal()
const query = ref('')
const { state, results, errorMessage, search } = useDirectory()

// debounce การพิมพ์ ~300ms กันยิง API ทุกตัวอักษร
let timer: ReturnType<typeof setTimeout> | undefined
watch(query, (q) => {
  clearTimeout(timer)
  timer = setTimeout(() => search(q.trim(), date), 300)
})

onMounted(() => search('', date))
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <AppHeader :date="date" variant="public" />

    <main class="mx-auto flex max-w-app flex-col gap-4 px-4 py-4">
      <h1 class="text-xl font-bold tracking-tight text-gray-900">ค้นหาร้านในตลาด</h1>

      <SearchBar v-model="query" />

      <ul v-if="state === 'ready'" class="flex flex-col gap-3">
        <li v-for="entry in results" :key="entry.stallCode">
          <DirectoryRow :entry="entry" />
        </li>
      </ul>

      <StateBlock
        v-else
        :state="state === 'error' ? 'error' : state === 'empty' ? 'empty' : 'loading'"
        :message="
          state === 'error'
            ? errorMessage
            : state === 'empty'
              ? query
                ? 'ไม่พบร้านที่ค้นหา'
                : 'วันนี้ยังไม่มีร้านค้า'
              : ''
        "
        @retry="search(query.trim(), date)"
      />
    </main>
  </div>
</template>
