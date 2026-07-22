<script setup lang="ts">
import type { AvailableStall } from '@/types'

// 📥 props = ข้อมูลที่ "แม่" ส่งลงมาให้ component นี้
//    item     = ข้อมูลแผง (มี stall.code, stall.zone)
//    selected = แผงนี้ถูกเลือกอยู่ไหม
//    disabled = แผงนี้กดไม่ได้ (ถูกจองแล้ว)
withDefaults(
  defineProps<{
    item: AvailableStall
    selected?: boolean
    disabled?: boolean
  }>(),
  { selected: false, disabled: false },
)

// 📤 emit = ส่งสัญญาณกลับขึ้นไปให้ "แม่" ว่าการ์ดนี้ถูกกด
const emit = defineEmits<{ (e: 'select', item: AvailableStall): void }>()
</script>

<template>
  <button
    type="button"
    :disabled="disabled"
    class="w-full rounded-3xl py-3.5 text-center text-lg font-bold transition"
    :class="
      disabled
        ? 'cursor-not-allowed bg-white/60 text-slate-300'
        : selected
          ? 'bg-brand-50 text-brand-700 ring-2 ring-brand-500 shadow-sm'
          : 'bg-white text-slate-800 shadow-sm active:bg-brand-50'
    "
    @click="emit('select', item)"
  >
    {{ item.stall.code }}
  </button>
</template>
