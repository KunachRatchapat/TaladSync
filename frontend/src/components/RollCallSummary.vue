<script setup lang="ts">
import { computed } from 'vue'
import SvgIcon from '@jamescoyle/vue-icon'
import { mdiCheckCircle } from '@mdi/js'

import type { AttendanceSummary } from '@/types'

const props = defineProps<{ summary: AttendanceSummary }>()

const checked = computed(() => props.summary.present + props.summary.noShow)
const pct = computed(() =>
  props.summary.total === 0 ? 0 : Math.round((checked.value / props.summary.total) * 100),
)
const done = computed(() => props.summary.total > 0 && props.summary.pending === 0)
</script>

<template>
  <section class="rounded-2xl border border-gray-200 bg-white p-4 shadow-sm">
    <div class="mb-2 flex items-baseline justify-between">
      <p v-if="done" class="flex items-center gap-1.5 text-sm font-bold text-brand-700">
        <SvgIcon type="mdi" :path="mdiCheckCircle" :size="18" aria-hidden="true" />
        เช็คครบแล้ว
      </p>
      <p v-else class="text-sm font-medium text-gray-700">
        เช็คแล้ว
        <span class="font-mono text-base font-bold text-gray-900">{{ checked }}</span>
        <span class="text-gray-400">/{{ summary.total }}</span>
      </p>
      <span class="font-mono text-xs font-medium text-gray-400">{{ pct }}%</span>
    </div>

    <!-- progress meter = signature -->
    <div class="h-2 overflow-hidden rounded-full bg-gray-100" role="progressbar">
      <div
        class="h-full rounded-full transition-all duration-300"
        :class="done ? 'bg-brand-500' : 'bg-brand-600'"
        :style="{ width: `${pct}%` }"
      />
    </div>

    <div class="mt-3 flex items-center gap-4 text-xs font-medium">
      <span class="flex items-center gap-1.5 text-brand-700">
        <span class="h-2.5 w-2.5 rounded-full bg-brand-600" aria-hidden="true" />
        มา {{ summary.present }}
      </span>
      <span class="flex items-center gap-1.5 text-danger-700">
        <span class="h-2.5 w-2.5 rounded-full bg-danger-600" aria-hidden="true" />
        ไม่มา {{ summary.noShow }}
      </span>
      <span class="flex items-center gap-1.5 text-gray-500">
        <span class="h-2.5 w-2.5 rounded-full bg-gray-300" aria-hidden="true" />
        เหลือ {{ summary.pending }}
      </span>
    </div>
  </section>
</template>
