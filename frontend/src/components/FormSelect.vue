<script setup lang="ts">
import SvgIcon from '@jamescoyle/vue-icon'
import { mdiChevronDown } from '@mdi/js'

const model = defineModel<string>({ required: true })

defineProps<{
  label: string
  options: readonly string[]
  placeholder: string
}>()
</script>

<template>
  <label class="flex flex-col gap-2">
    <span class="text-sm font-semibold text-slate-600">{{ label }}</span>
    <div class="relative">
      <select
        v-model="model"
        class="h-11 w-full appearance-none rounded-xl border border-slate-200 bg-white pl-4 pr-12 text-base focus:border-emerald-500 focus:outline-none focus:ring-1 focus:ring-emerald-500"
        :class="model ? 'text-slate-800' : 'text-slate-400'"
      >
        <option value="" disabled>{{ placeholder }}</option>
        <option v-for="option in options" :key="option" :value="option" class="text-slate-800">
          {{ option }}
        </option>
      </select>
      <span class="pointer-events-none absolute inset-y-0 right-4 flex items-center text-slate-400">
        <SvgIcon type="mdi" :path="mdiChevronDown" :size="24" />
      </span>
    </div>
  </label>
</template>
