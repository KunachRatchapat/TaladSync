<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import SvgIcon from '@jamescoyle/vue-icon'
import { mdiCheckCircleOutline } from '@mdi/js'

import BaseButton from './BaseButton.vue'
import { useBooking } from '@/composables/useBooking'
import type { AvailableStall } from '@/types'

const props = defineProps<{
  item: AvailableStall | null
  date: string
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'booked'): void
}>()

const shopName = ref('')
const productType = ref('')

const { submitting, success, errorMessage, errorCode, submit, reset } = useBooking()

const canSubmit = computed(
  () => shopName.value.trim() !== '' && productType.value.trim() !== '' && !submitting.value,
)

// เปิด modal ใหม่ทุกครั้งที่เลือกแผง = ล้างฟอร์ม
watch(
  () => props.item,
  (item) => {
    if (item) {
      shopName.value = ''
      productType.value = ''
      reset()
    }
  },
)

async function onSubmit() {
  if (!props.item) return
  await submit({
    stallId: props.item.stallId,
    shopName: shopName.value.trim(),
    productType: productType.value.trim(),
    date: props.date,
  })
  // จองสำเร็จ หรือ "จองครบแล้ววันนี้" (มีแผงอยู่แล้ว) → ให้หน้าหลักดึงการจองมาโชว์การ์ด
  if (success.value || errorCode.value === 'ALREADY_HAS_STALL') {
    emit('booked')
  }
}

function onClose() {
  emit('close')
}
</script>

<template>
  <Transition name="sheet">
    <div
      v-if="item"
      class="fixed inset-0 z-20 flex items-end justify-center bg-black/40"
      role="dialog"
      aria-modal="true"
      aria-label="ฟอร์มจองแผง"
      @click.self="onClose"
    >
      <div class="w-full max-w-app rounded-t-2xl bg-white p-5 pb-8 shadow-xl">
        <!-- แถบจับ -->
        <div class="mx-auto mb-4 h-1.5 w-10 rounded-full bg-gray-200" aria-hidden="true" />

        <!-- success -->
        <div v-if="success" class="flex flex-col items-center gap-4 py-6 text-center">
          <span
            class="flex h-16 w-16 items-center justify-center rounded-full bg-brand-50 text-brand-600"
            aria-hidden="true"
          >
            <SvgIcon type="mdi" :path="mdiCheckCircleOutline" :size="40" />
          </span>
          <div>
            <p class="text-lg font-semibold text-gray-900">จองสำเร็จ</p>
            <p class="mt-1 text-sm text-gray-500">
              จองสำเร็จ {{ item.stall.code }} เริ่มขายได้ตั้งแต่ 16:00 เป็นต้นไป
            </p>
          </div>
          <BaseButton @click="onClose">เสร็จสิ้น</BaseButton>
        </div>

        <!-- form -->
        <form v-else class="flex flex-col gap-4" @submit.prevent="onSubmit">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-lg font-bold text-gray-900">จองแผง {{ item.stall.code }}</p>
              <p class="text-sm text-gray-500">{{ item.stall.zone }} · {{ date }}</p>
            </div>
            <button
              type="button"
              class="rounded-full px-2 text-2xl leading-none text-gray-400 active:text-gray-600"
              aria-label="ปิด"
              @click="onClose"
            >
              ×
            </button>
          </div>

          <label class="flex flex-col gap-1">
            <span class="text-sm font-medium text-gray-700">ชื่อร้าน</span>
            <input
              v-model="shopName"
              maxlength="100"
              type="text"
              placeholder="เช่น ร้านป้าสมทรง"
              class="w-full rounded-xl border border-gray-300 px-4 py-3 text-gray-900 focus:border-brand-600 focus:outline-none focus:ring-1 focus:ring-brand-600"
            />
          </label>

          <label class="flex flex-col gap-1">
            <span class="text-sm font-medium text-gray-700">ประเภทสินค้า</span>
            <input
              v-model="productType"
              maxlength="100"
              type="text"
              placeholder="เช่น ผลไม้ / ของสด / เสื้อผ้า"
              class="w-full rounded-xl border border-gray-300 px-4 py-3 text-gray-900 focus:border-brand-600 focus:outline-none focus:ring-1 focus:ring-brand-600"
            />
          </label>

          <!-- error 409 ภาษาไทย -->
          <p
            v-if="errorMessage"
            class="rounded-xl bg-danger-50 px-4 py-3 text-sm font-medium text-danger-700"
            role="alert"
          >
            {{ errorMessage }}
          </p>

          <BaseButton type="submit" :disabled="!canSubmit" :loading="submitting">
            ยืนยันการจอง
          </BaseButton>
        </form>
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
