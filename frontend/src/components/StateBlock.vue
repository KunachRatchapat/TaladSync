<script setup lang="ts">
import BaseButton from './BaseButton.vue'

withDefaults(
  defineProps<{
    state: 'loading' | 'empty' | 'error'
    message?: string
  }>(),
  { message: '' },
)

const emit = defineEmits<{ (e: 'retry'): void }>()
</script>

<template>
  <!-- loading = skeleton การ์ด -->
  <div v-if="state === 'loading'" class="flex flex-col gap-3" aria-busy="true">
    <div
      v-for="n in 4"
      :key="n"
      class="h-[76px] animate-pulse rounded-2xl border border-gray-200 bg-gray-100"
    />
  </div>

  <!-- empty -->
  <div
    v-else-if="state === 'empty'"
    class="flex flex-col items-center gap-3 rounded-2xl bg-warning-50 px-6 py-10 text-center"
  >
    <span class="text-3xl" aria-hidden="true">🗓️</span>
    <p class="text-sm font-medium text-warning-700">
      {{ message || 'วันนี้ไม่มีแผงว่างแล้ว' }}
    </p>
  </div>

  <!-- error -->
  <div
    v-else
    class="flex flex-col items-center gap-4 rounded-2xl bg-danger-50 px-6 py-10 text-center"
  >
    <span class="text-3xl" aria-hidden="true">⚠️</span>
    <p class="text-sm font-medium text-danger-700">
      {{ message || 'เกิดข้อผิดพลาด ลองใหม่อีกครั้ง' }}
    </p>
    <BaseButton variant="secondary" :block="false" @click="emit('retry')">
      ลองใหม่
    </BaseButton>
  </div>
</template>
