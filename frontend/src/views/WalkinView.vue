<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import BookedTodayCard from '@/components/BookedTodayCard.vue'
import BookingModal from '@/components/BookingModal.vue'
import PageHeader from '@/components/PageHeader.vue'
import StallCard from '@/components/StallCard.vue'
import StateBlock from '@/components/StateBlock.vue'
import { fetchMyBookingToday } from '@/api/booking'
import { useAvailableStalls } from '@/composables/useAvailableStalls'
import { useAuthStore } from '@/stores/auth'
import type { AvailableStall } from '@/types'
import { todayLocal } from '@/utils/date'

const date = todayLocal()
const router = useRouter()
const auth = useAuthStore()
const myBooking = ref<AvailableStall | null>(null)
const { state, stalls, errorMessage, load } = useAvailableStalls()

// จัดกลุ่มแผงตามโซน + เรียงตามรหัสแผง → ได้ [{ zone: 'โซน A', items: [...] }, ...]
const stallsByZone = computed(() => {
  const groups = new Map<string, AvailableStall[]>()
  for (const item of stalls.value) {
    const zone = item.stall.zone
    if (!groups.has(zone)) groups.set(zone, [])
    groups.get(zone)!.push(item)
  }
  return [...groups.entries()]
    .map(([zone, items]) => ({
      zone,
      items: [...items].sort((a, b) => a.stall.code.localeCompare(b.stall.code)),
    }))
    .sort((a, b) => a.zone.localeCompare(b.zone))
})

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

const selected = ref<AvailableStall | null>(null)
function handleSelect(item: AvailableStall) {
  if (myBooking.value) return
  selected.value = item
}
function closeModal() {
  selected.value = null
}
async function onBooked() {
  myBooking.value = await fetchMyBookingToday(date)
}
</script>

<template>
  <div class="min-h-screen bg-brand-50">
    <PageHeader title="จองแผงค้า" @back="router.push({ name: 'walkin' })" />

    <main class="mx-auto max-w-app px-4 py-5">
      <h1 class="mb-5 text-center text-lg font-bold text-brand-700">
        {{ myBooking ? 'การจองของคุณวันนี้' : 'เลือกแผงค้าที่คุณต้องการจอง' }}
      </h1>

      <BookedTodayCard v-if="myBooking" :booking="myBooking" />

      <template v-else>
        <StateBlock
          v-if="auth.loading || state === 'loading' || state === 'empty' || state === 'error' || auth.error"
          :state="
            auth.loading || state === 'loading'
              ? 'loading'
              : auth.error || state === 'error'
                ? 'error'
                : 'empty'
          "
          :message="auth.error || errorMessage"
          @retry="reload"
        />

        <!-- วนทีละโซน → แต่ละโซนมีหัวข้อ + กริด 3 คอลัมน์ -->
        <div v-else-if="state === 'ready'" class="flex flex-col gap-6">
          <section v-for="group in stallsByZone" :key="group.zone">
            <h2 class="mb-3 text-base font-bold text-slate-800">{{ group.zone }}</h2>

            <div class="grid grid-cols-3 gap-3">
              <StallCard
                v-for="item in group.items"
                :key="item.id"
                :item="item"
                @select="handleSelect"
              />
            </div>
          </section>
        </div>
      </template>
    </main>

    <BookingModal
      v-if="!myBooking"
      :item="selected"
      :date="date"
      @close="closeModal"
      @booked="onBooked"
    />
  </div>
</template>
