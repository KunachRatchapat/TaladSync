<script setup lang="ts">
import { computed } from 'vue'

import { formatThaiDate } from '@/utils/date'

const props = withDefaults(
  defineProps<{
    date: string
    variant?: 'walkin' | 'staff' | 'regular' | 'public'
  }>(),
  { variant: 'walkin' },
)

const thaiDate = computed(() => formatThaiDate(props.date))
// แต่ละ role มี badge ต่างสี ให้รู้ทันทีว่าอยู่โหมดไหน
const tag = computed(() => {
  if (props.variant === 'staff') return 'เจ้าหน้าที่ · เช็คชื่อ'
  if (props.variant === 'regular') return 'แม่ค้าประจำ · แผงของฉัน'
  if (props.variant === 'public') return 'ค้นหาร้านค้า'
  return 'จองแผงวันนี้ได้เลยจ้า \u{1F603}'
})
const tagClass = computed(() => {
  if (props.variant === 'staff') return 'bg-gray-900 text-white'
  if (props.variant === 'regular') return 'bg-warning-100 text-warning-700'
  if (props.variant === 'public') return 'bg-sky-100 text-sky-700'
  return 'bg-brand-50 text-brand-700'
})
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
