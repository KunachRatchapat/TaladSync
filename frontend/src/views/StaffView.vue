<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

import AppHeader from '@/components/AppHeader.vue'
import ChecklistRow from '@/components/ChecklistRow.vue'
import RollCallSummary from '@/components/RollCallSummary.vue'
import StateBlock from '@/components/StateBlock.vue'
import { useChecklist } from '@/composables/useChecklist'
import { useAuthStore } from '@/stores/auth'
import type { AttendanceStatus } from '@/types'
import { todayLocal } from '@/utils/date'

const date = todayLocal()
const auth = useAuthStore()
const { state, items, errorMessage, marking, summary, load, mark } = useChecklist()

type Filter = 'all' | 'pending' | 'present' | 'noshow'
const filter = ref<Filter>('all')

const filters: { key: Filter; label: string; match?: AttendanceStatus }[] = [
  { key: 'all', label: 'ทั้งหมด' },
  { key: 'pending', label: 'ยังไม่เช็ค', match: 'PENDING' },
  { key: 'present', label: 'มา', match: 'PRESENT' },
  { key: 'noshow', label: 'ไม่มา', match: 'NOSHOW' },
]

const visibleItems = computed(() => {
  if (filter.value === 'all') return items.value
  const f = filters.find((it) => it.key === filter.value)
  return items.value.filter((it) => it.attendance === f?.match)
})

async function reload() {
  try {
    await auth.ensureRole('STAFF')
  } catch {
    return
  }
  await load(date)
}

onMounted(reload)
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <AppHeader :date="date" variant="staff" />

    <main class="mx-auto flex max-w-app flex-col gap-4 px-4 py-4">
      <h1 class="text-xl font-bold tracking-tight text-gray-900">เช็คชื่อหน้างาน</h1>

      <template v-if="state === 'ready'">
        <RollCallSummary :summary="summary" />

        <!-- filter chips -->
        <div class="flex gap-2 overflow-x-auto pb-1">
          <button
            v-for="f in filters"
            :key="f.key"
            type="button"
            class="shrink-0 rounded-full border px-3.5 py-1.5 text-xs font-medium transition"
            :class="
              filter === f.key
                ? 'border-gray-900 bg-gray-900 text-white'
                : 'border-gray-300 bg-white text-gray-600 active:bg-gray-50'
            "
            @click="filter = f.key"
          >
            {{ f.label }}
          </button>
        </div>

        <!-- error จากการ mark (optimistic rollback) -->
        <p
          v-if="errorMessage"
          class="rounded-xl bg-danger-50 px-4 py-2.5 text-sm font-medium text-danger-700"
          role="alert"
        >
          {{ errorMessage }}
        </p>

        <ul class="flex flex-col gap-3">
          <li v-for="item in visibleItems" :key="item.id">
            <ChecklistRow :item="item" :busy="marking.has(item.id)" @mark="mark" />
          </li>
        </ul>

        <p
          v-if="visibleItems.length === 0"
          class="py-8 text-center text-sm text-gray-400"
        >
          ไม่มีแผงในตัวกรองนี้
        </p>
      </template>

      <StateBlock
        v-else
        :state="
          auth.loading || state === 'loading'
            ? 'loading'
            : auth.error || state === 'error'
              ? 'error'
              : 'empty'
        "
        :message="auth.error || errorMessage || (state === 'empty' ? 'วันนี้ยังไม่มีแผงให้เช็คชื่อ' : '')"
        @retry="reload"
      />
    </main>
  </div>
</template>
