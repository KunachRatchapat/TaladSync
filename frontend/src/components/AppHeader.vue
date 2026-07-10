<script setup lang="ts">
import { computed } from 'vue'

import { formatThaiDate } from '@/utils/date'

const props = defineProps<{ date: string }>()

const thaiDate = computed(() => formatThaiDate(props.date))
</script>

<template>
  <header class="sticky top-0 z-10 border-b border-gray-200 bg-white/95 backdrop-blur">
    <div class="mx-auto flex max-w-app flex-col gap-2 px-4 py-3">
      <div class="flex items-center gap-2">
        <span class="text-lg font-semibold text-brand-700">TaladSync</span>
        <span class="rounded-full bg-brand-50 px-2 py-0.5 text-xs font-medium text-brand-700">
          จองแผงวันนี้ได้เลยจ้า &#128515;
        </span>
      </div>

      <p class="text-sm text-gray-500">
        <span class="font-medium text-gray-700">วันนี้</span> · {{ thaiDate }}
      </p>
    </div>
  </header>
</template>
