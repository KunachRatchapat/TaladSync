<script setup lang="ts">
import { computed } from 'vue'

import BaseButton from './BaseButton.vue'
import type { AvailableStall } from '@/types'

const props = defineProps<{ item: AvailableStall; busy?: boolean }>()

const emit = defineEmits<{ (e: 'leave', item: AvailableStall): void }>()

// EXPECTED = ยังจะมาขาย, AVAILABLE = แจ้งหยุดแล้ว (ปล่อยพูล), BOOKED = ขาจรจองต่อแล้ว
const isOpen = computed(() => props.item.status === 'EXPECTED')

const statusText = computed(() => {
  if (props.item.status === 'EXPECTED') return 'เปิดขายวันนี้'
  if (props.item.status === 'BOOKED') return 'แจ้งหยุดแล้ว · ขาจรจองต่อแล้ว'
  return 'แจ้งหยุดแล้ว · ปล่อยให้ขาจร'
})
const statusClass = computed(() =>
  isOpen.value ? 'bg-brand-50 text-brand-700' : 'bg-gray-100 text-gray-500',
)
</script>

<template>
  <div
    class="rounded-2xl border border-l-4 bg-white p-4 shadow-sm transition-colors"
    :class="isOpen ? 'border-gray-200 border-l-brand-500' : 'border-gray-200 border-l-gray-200'"
  >
    <div class="flex items-center gap-2">
      <span class="rounded-md bg-gray-900 px-2 py-0.5 font-mono text-sm font-bold text-white">
        {{ item.stall.code }}
      </span>
      <span class="text-xs text-gray-400">{{ item.stall.zone }}</span>
      <span class="ml-auto rounded-full px-2.5 py-0.5 text-xs font-medium" :class="statusClass">
        {{ statusText }}
      </span>
    </div>

    <p v-if="item.shopName" class="mt-2 truncate font-medium text-gray-900">{{ item.shopName }}</p>

    <BaseButton
      v-if="isOpen"
      variant="secondary"
      :loading="busy"
      class="mt-3"
      @click="emit('leave', item)"
    >
      แจ้งหยุดวันนี้
    </BaseButton>
  </div>
</template>
