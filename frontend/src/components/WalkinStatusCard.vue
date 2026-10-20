<script setup lang="ts">
import SvgIcon from '@jamescoyle/vue-icon'
import { mdiStorefrontOutline } from '@mdi/js'

import type { AvailableStall } from '@/types'

defineProps<{
  shopName: string
  booking: AvailableStall | null
  availableCount: number
}>()
</script>

<template>
  <section class="flex items-center justify-between gap-4 rounded-3xl bg-emerald-600 p-4 text-white">
    <div class="flex min-w-0 flex-col items-start">
      <span class="rounded-full bg-indigo-50 px-3 py-1 text-sm font-medium text-indigo-500">
        สถานะ : {{ booking ? 'จองแล้ว' : 'ไม่มีที่ประจำ' }}
      </span>

      <p class="mt-3 truncate text-lg font-bold">{{ shopName }}</p>

      <template v-if="booking">
        <p class="mt-2 text-sm">แผงของคุณ</p>
        <p class="text-xl font-bold">{{ booking.stall.code }}</p>
      </template>
      <template v-else>
        <p class="mt-2 text-sm">เหลือล็อคว่าง</p>
        <p class="text-xl font-bold">{{ availableCount }} ล็อค</p>
      </template>
    </div>

    <SvgIcon type="mdi" :path="mdiStorefrontOutline" :size="88" class="shrink-0" />
  </section>
</template>
