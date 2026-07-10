<script setup lang="ts">
import SvgIcon from '@jamescoyle/vue-icon'
import { mdiCheckCircleOutline } from '@mdi/js'

import type { AvailableStall } from '@/types'

defineProps<{ booking: AvailableStall }>()
</script>

<template>
  <div
    class="flex flex-col items-center gap-4 rounded-2xl border border-brand-100 bg-white/90 px-6 py-8 text-center shadow-sm backdrop-blur"
  >
    <span
      class="flex h-14 w-14 items-center justify-center rounded-full bg-brand-50 text-brand-600"
      aria-hidden="true"
    >
      <SvgIcon type="mdi" :path="mdiCheckCircleOutline" :size="36" />
    </span>
    <div>
      <p class="text-lg font-semibold text-gray-900">คุณจองแผงวันนี้แล้ว</p>
      <p class="mt-2 text-2xl font-bold text-brand-700">{{ booking.stall.code }}</p>
      <p class="mt-1 text-sm text-gray-500">{{ booking.stall.zone }}</p>
      <p v-if="booking.shopName" class="mt-3 text-sm text-gray-700">
        {{ booking.shopName }} · {{ booking.productType }}
      </p>
      <p class="mt-3 text-xs text-gray-400">เริ่มขายได้ตั้งแต่ 16:00 เป็นต้นไป</p>
    </div>
  </div>
</template>
