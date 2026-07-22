<script setup lang="ts">
import SvgIcon from '@jamescoyle/vue-icon'
import { mdiChevronLeft } from '@mdi/js'

withDefaults(
  defineProps<{
    title: string
    showBack?: boolean
  }>(),
  { showBack: true },
)

const emit = defineEmits<{ (e: 'back'): void }>()
</script>

<template>
  <header class="sticky top-0 z-10 bg-white">
    <!-- relative + absolute = ทำให้ title อยู่กลางจอเป๊ะ ไม่ถูกปุ่มย้อนกลับดันให้เบี้ยว -->
    <div class="relative mx-auto flex h-14 max-w-app items-center justify-center px-2">
      <button
        v-if="showBack"
        type="button"
        class="absolute left-1 rounded-full p-2 text-brand-600 active:bg-brand-50"
        aria-label="ย้อนกลับ"
        @click="emit('back')"
      >
        <SvgIcon type="mdi" :path="mdiChevronLeft" :size="28" />
      </button>

      <h1 class="text-base font-bold text-slate-800">{{ title }}</h1>
    </div>
  </header>
</template>
