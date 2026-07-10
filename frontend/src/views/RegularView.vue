<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

import AppHeader from '@/components/AppHeader.vue'
import LeaveConfirmModal from '@/components/LeaveConfirmModal.vue'
import MyStallCard from '@/components/MyStallCard.vue'
import StateBlock from '@/components/StateBlock.vue'
import { useMyStalls } from '@/composables/useMyStalls'
import { useAuthStore } from '@/stores/auth'
import type { AvailableStall } from '@/types'
import { todayLocal } from '@/utils/date'

const date = todayLocal()
const auth = useAuthStore()
const { state, stalls, errorMessage, leaving, load, leave } = useMyStalls()

const openCount = computed(() => stalls.value.filter((s) => s.status === 'EXPECTED').length)

const selected = ref<AvailableStall | null>(null)
const confirming = computed(() => (selected.value ? leaving.value.has(selected.value.stallId) : false))

function askLeave(item: AvailableStall) {
  selected.value = item
}
async function confirmLeave() {
  if (!selected.value) return
  const stallId = selected.value.stallId
  await leave(stallId, date)
  // ปิด modal เฉพาะตอนสำเร็จ (status เปลี่ยนเป็น AVAILABLE) — error ค้างไว้ให้เห็น
  const item = stalls.value.find((s) => s.stallId === stallId)
  if (item && item.status !== 'EXPECTED') selected.value = null
}

async function reload() {
  try {
    await auth.ensureRole('REGULAR')
  } catch {
    return
  }
  await load(date)
}

onMounted(reload)
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <AppHeader :date="date" variant="regular" />

    <main class="mx-auto flex max-w-app flex-col gap-4 px-4 py-4">
      <div>
        <h1 class="text-xl font-bold tracking-tight text-gray-900">แผงของฉัน</h1>
        <p v-if="state === 'ready'" class="mt-1 text-sm text-gray-500">
          เปิดขายวันนี้
          <span class="font-mono font-bold text-brand-700">{{ openCount }}</span>
          <span class="text-gray-400">/{{ stalls.length }} แผง</span>
        </p>
      </div>

      <template v-if="state === 'ready'">
        <p
          v-if="errorMessage"
          class="rounded-xl bg-danger-50 px-4 py-2.5 text-sm font-medium text-danger-700"
          role="alert"
        >
          {{ errorMessage }}
        </p>

        <ul class="flex flex-col gap-3">
          <li v-for="item in stalls" :key="item.stallId">
            <MyStallCard :item="item" :busy="leaving.has(item.stallId)" @leave="askLeave" />
          </li>
        </ul>
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
        :message="
          auth.error || errorMessage || (state === 'empty' ? 'คุณยังไม่มีแผงประจำในระบบ' : '')
        "
        @retry="reload"
      />
    </main>

    <LeaveConfirmModal
      :item="selected"
      :busy="confirming"
      @close="selected = null"
      @confirm="confirmLeave"
    />
  </div>
</template>
