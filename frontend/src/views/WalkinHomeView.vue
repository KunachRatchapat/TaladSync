<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  mdiAccountOutline,
  mdiCalendarMonthOutline,
  mdiCancel,
  mdiClockOutline,
  mdiHistory,
  mdiStorefrontOutline,
} from '@mdi/js'

import MenuTile from '@/components/MenuTile.vue'
import PageHeader from '@/components/PageHeader.vue'
import StateBlock from '@/components/StateBlock.vue'
import WalkinStatusCard from '@/components/WalkinStatusCard.vue'
import { fetchMyBookingToday } from '@/api/booking'
import { useAvailableStalls } from '@/composables/useAvailableStalls'
import { useAuthStore } from '@/stores/auth'
import type { AvailableStall } from '@/types'
import { formatThaiDate, todayLocal } from '@/utils/date'

const date = todayLocal()
const router = useRouter()
const auth = useAuthStore()
const myBooking = ref<AvailableStall | null>(null)
const { state, stalls, load } = useAvailableStalls()

const shopName = computed(() => myBooking.value?.shopName || auth.user?.name || '')
const loading = computed(() => auth.loading || state.value === 'loading')

const menu = computed(() => {
  const booked = myBooking.value !== null
  return [
    { label: 'จองแผง', icon: mdiCalendarMonthOutline, disabled: booked, to: 'walkin-book' },
    { label: 'ยกเลิกการจอง', icon: mdiCancel, disabled: true },
    { label: 'แจ้งเปิดแผงสาย', icon: mdiClockOutline, disabled: true },
    { label: 'ข้อมูลร้านค้า', icon: mdiStorefrontOutline, disabled: false, to: 'walkin-shop' },
    { label: 'ประวัติการจอง', icon: mdiHistory, disabled: false },
    { label: 'โปรไฟล์', icon: mdiAccountOutline, disabled: false },
  ]
})

function open(to?: string) {
  if (to) router.push({ name: to })
}

async function reload() {
  try {
    await auth.ensureRole('WALKIN')
  } catch {
    return
  }
  myBooking.value = await fetchMyBookingToday(date)
  if (!myBooking.value) {
    await load(date)
  }
}

onMounted(reload)
</script>

<template>
  <div class="min-h-screen bg-brand-50">
    <PageHeader title="TaladSync · ตลาดเทศบาล" :show-back="false" />

    <main class="mx-auto max-w-app px-4 py-5">
      <p class="mb-4 text-base text-slate-600">{{ formatThaiDate(date) }}</p>

      <StateBlock
        v-if="loading || auth.error || state === 'error'"
        :state="loading ? 'loading' : 'error'"
        :message="auth.error || 'โหลดข้อมูลไม่สำเร็จ ลองใหม่อีกครั้ง'"
        @retry="reload"
      />

      <template v-else>
        <WalkinStatusCard
          :shop-name="shopName"
          :booking="myBooking"
          :available-count="stalls.length"
        />

        <h2 class="mb-3 mt-5 text-base font-semibold text-slate-600">เมนู</h2>

        <div class="grid grid-cols-3 gap-x-3 gap-y-4">
          <MenuTile
            v-for="item in menu"
            :key="item.label"
            :icon="item.icon"
            :label="item.label"
            :disabled="item.disabled"
            @click="open(item.to)"
          />
        </div>
      </template>
    </main>
  </div>
</template>
