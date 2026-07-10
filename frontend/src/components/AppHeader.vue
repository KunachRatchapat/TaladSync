<script setup lang="ts">
import { computed } from 'vue'

import { formatThaiDate } from '@/utils/date'

const props = withDefaults(
  defineProps<{
    date: string
    variant?: 'walkin' | 'staff'
  }>(),
  { variant: 'walkin' },
)

const thaiDate = computed(() => formatThaiDate(props.date))
const tag = computed(() =>
  props.variant === 'staff' ? 'เจ้าหน้าที่ · เช็คชื่อ' : 'จองแผงวันนี้ได้เลยจ้า \u{1F603}',
)
// โหมด staff = badge สีเข้ม (ink) ให้รู้สึก "โหมดทำงาน" ต่างจากขาจรสีเขียว
const tagClass = computed(() =>
  props.variant === 'staff' ? 'bg-gray-900 text-white' : 'bg-brand-50 text-brand-700',
)
</script>

<template>
  <header class="sticky top-0 z-10 border-b border-gray-200 bg-white/95 backdrop-blur">
    <div class="mx-auto flex max-w-app flex-col gap-2 px-4 py-3">
      <div class="flex items-center gap-2">
        <span class="text-lg font-semibold text-brand-700">TaladSync</span>
        <span class="rounded-full px-2 py-0.5 text-xs font-medium" :class="tagClass">
          {{ tag }}
        </span>
      </div>

      <p class="text-sm text-gray-500">
        <span class="font-medium text-gray-700">วันนี้</span> · {{ thaiDate }}
      </p>
    </div>
  </header>
</template>
