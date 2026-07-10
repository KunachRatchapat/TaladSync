<script setup lang="ts">
import type { DirectoryEntry } from '@/types'

defineProps<{ entry: DirectoryEntry }>()
</script>

<template>
  <div class="flex items-center gap-3 rounded-2xl border border-gray-200 bg-white p-4 shadow-sm">
    <span
      class="shrink-0 rounded-md bg-gray-900 px-2 py-0.5 font-mono text-sm font-bold text-white"
    >
      {{ entry.stallCode }}
    </span>
    <div class="min-w-0">
      <p class="truncate font-medium text-gray-900">{{ entry.shopName }}</p>
      <p class="truncate text-sm text-gray-500">
        {{ entry.productType }}<span class="text-gray-300"> · </span>{{ entry.zone }}
      </p>
    </div>
  </div>
</template>
