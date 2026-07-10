<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

import AppHeader from '@/components/AppHeader.vue'
import BookedTodayCard from '@/components/BookedTodayCard.vue'
import BookingModal from '@/components/BookingModal.vue'
import StallCard from '@/components/StallCard.vue'
import StateBlock from '@/components/StateBlock.vue'
import { fetchMyBookingToday } from '@/api/booking'
import { useAvailableStalls } from '@/composables/useAvailableStalls'
import { useAuthStore } from '@/stores/auth'
import type { AvailableStall } from '@/types'
import { todayLocal } from '@/utils/date'

const date = todayLocal()
const auth = useAuthStore()
const myBooking = ref<AvailableStall | null>(null)
const { state, stalls, errorMessage, load } = useAvailableStalls()

const stallCount = computed(() => stalls.value.length)

async function reload() {
  try {
    await auth.ensureWalkin()
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
  <div class="min-h-screen bg-gray-50">
    <AppHeader :date="date" />

    <main class="mx-auto max-w-app px-4 py-4">
      <div class="mb-5 flex flex-col items-center gap-2">
        <div class="flex items-center gap-2">
          <h1 class="text-xl font-bold tracking-tight text-brand-700">
            {{ myBooking ? 'การจองของคุณวันนี้' : 'แผงว่างวันนี้' }}
          </h1>
          <span
            v-if="!myBooking && state === 'ready'"
            class="rounded-full bg-brand-600 px-2.5 py-0.5 text-xs font-bold text-white"
          >
            {{ stallCount }}
          </span>
        </div>
        <span
          v-if="!myBooking && state === 'ready'"
          class="rounded-full bg-brand-50 px-3 py-1 text-xs font-medium text-brand-700"
        >
          {{ stallCount }} แผงว่าง
        </span>
      </div>

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

        <ul v-else-if="state === 'ready'" class="flex flex-col gap-3">
          <li v-for="item in stalls" :key="item.id">
            <StallCard :item="item" @select="handleSelect" />
          </li>
        </ul>
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
