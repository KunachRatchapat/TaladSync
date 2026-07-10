<script setup lang="ts">
import type { AvailableStall } from '@/types'

defineProps<{ item: AvailableStall }>()

const emit = defineEmits<{ (e: 'select', item: AvailableStall): void }>()
</script>

<template>
  <button
    type="button"
    class="flex w-full items-center justify-between rounded-2xl border border-gray-200 bg-white p-4 text-left shadow-sm transition active:border-brand-500 active:bg-brand-50"
    @click="emit('select', item)"
  >
    <div class="flex flex-col gap-1">
      <span class="text-lg font-bold text-gray-900">{{ item.stall.code }}</span>
      <span class="text-sm text-gray-500">{{ item.stall.zone }}</span>
    </div>

    <span
      class="rounded-full bg-brand-50 px-3 py-1 text-xs font-medium text-brand-700"
    >
      ว่าง
    </span>
  </button>
</template>
