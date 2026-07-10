<script setup lang="ts">
import { computed } from 'vue'

import type { ChecklistItem } from '@/types'

const props = defineProps<{ item: ChecklistItem; busy?: boolean }>()

const emit = defineEmits<{ (e: 'mark', id: number, present: boolean): void }>()

const isPresent = computed(() => props.item.attendance === 'PRESENT')
const isNoShow = computed(() => props.item.attendance === 'NOSHOW')

// เส้นซ้ายบอกสถานะแบบเหลือบเห็นเร็ว
const accent = computed(() => {
  if (isPresent.value) return 'border-l-brand-500'
  if (isNoShow.value) return 'border-l-danger-500'
  return 'border-l-gray-200'
})
</script>

<template>
  <div
    class="rounded-2xl border border-l-4 border-gray-200 bg-white p-4 shadow-sm transition-colors"
    :class="accent"
  >
    <div class="mb-3 flex items-start justify-between gap-3">
      <div class="min-w-0">
        <div class="flex items-center gap-2">
          <span class="rounded-md bg-gray-900 px-2 py-0.5 font-mono text-sm font-bold text-white">
            {{ item.stall.code }}
          </span>
          <span class="truncate text-xs text-gray-400">{{ item.stall.zone }}</span>
        </div>
        <p class="mt-1.5 truncate font-medium text-gray-900">
          {{ item.shopName || '—' }}
        </p>
        <p v-if="item.productType" class="truncate text-sm text-gray-500">
          {{ item.productType }}
        </p>
      </div>
    </div>

    <div class="flex gap-2">
      <button
        type="button"
        :disabled="busy"
        :aria-pressed="isPresent"
        class="flex-1 rounded-xl border py-3 text-sm font-semibold transition disabled:opacity-50"
        :class="
          isPresent
            ? 'border-brand-600 bg-brand-600 text-white'
            : 'border-gray-300 bg-white text-gray-700 active:bg-gray-50'
        "
        @click="emit('mark', item.id, true)"
      >
        มา
      </button>
      <button
        type="button"
        :disabled="busy"
        :aria-pressed="isNoShow"
        class="flex-1 rounded-xl border py-3 text-sm font-semibold transition disabled:opacity-50"
        :class="
          isNoShow
            ? 'border-danger-600 bg-danger-600 text-white'
            : 'border-gray-300 bg-white text-gray-700 active:bg-gray-50'
        "
        @click="emit('mark', item.id, false)"
      >
        ไม่มา
      </button>
    </div>
  </div>
</template>
