<script setup lang="ts">
import SvgIcon from '@jamescoyle/vue-icon'
import { mdiAlertOutline } from '@mdi/js'

import BaseButton from './BaseButton.vue'
import type { AvailableStall } from '@/types'

defineProps<{ item: AvailableStall | null; busy?: boolean }>()

const emit = defineEmits<{ (e: 'close'): void; (e: 'confirm'): void }>()
</script>

<template>
  <Transition name="sheet">
    <div
      v-if="item"
      class="fixed inset-0 z-20 flex items-end justify-center bg-black/40"
      role="dialog"
      aria-modal="true"
      aria-label="ยืนยันแจ้งหยุด"
      @click.self="emit('close')"
    >
      <div class="w-full max-w-app rounded-t-2xl bg-white p-5 pb-8 shadow-xl">
        <div class="mx-auto mb-4 h-1.5 w-10 rounded-full bg-gray-200" aria-hidden="true" />

        <div class="flex flex-col items-center gap-3 text-center">
          <span
            class="flex h-14 w-14 items-center justify-center rounded-full bg-warning-100 text-warning-600"
            aria-hidden="true"
          >
            <SvgIcon type="mdi" :path="mdiAlertOutline" :size="32" />
          </span>
          <p class="text-lg font-bold text-gray-900">แจ้งหยุดแผง {{ item.stall.code }}?</p>
          <p class="text-sm text-gray-500">
            แผงจะถูกปล่อยให้ขาจรจองแทนทันที และ<span class="font-semibold text-danger-700"
              >ยกเลิกไม่ได้</span
            >
          </p>
        </div>

        <div class="mt-5 flex flex-col gap-2">
          <BaseButton variant="danger" :loading="busy" @click="emit('confirm')">
            ยืนยันแจ้งหยุด
          </BaseButton>
          <BaseButton variant="secondary" :disabled="busy" @click="emit('close')">
            ยกเลิก
          </BaseButton>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.sheet-enter-active,
.sheet-leave-active {
  transition: opacity 0.2s ease;
}
.sheet-enter-from,
.sheet-leave-to {
  opacity: 0;
}
</style>
