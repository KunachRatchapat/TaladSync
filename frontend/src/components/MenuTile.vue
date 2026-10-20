<script setup lang="ts">
import SvgIcon from '@jamescoyle/vue-icon'

withDefaults(
  defineProps<{
    icon: string
    label: string
    disabled?: boolean
  }>(),
  { disabled: false },
)

const emit = defineEmits<{ (e: 'click'): void }>()
</script>

<template>
  <button
    type="button"
    :disabled="disabled"
    class="flex w-full flex-col items-center justify-center gap-2 rounded-3xl py-3 transition"
    :class="
      disabled
        ? 'cursor-not-allowed bg-slate-200 text-slate-400'
        : 'bg-white text-slate-700 active:bg-brand-50'
    "
    @click="emit('click')"
  >
    <SvgIcon
      type="mdi"
      :path="icon"
      :size="28"
      :class="disabled ? 'text-slate-400' : 'text-emerald-500'"
    />
    <span class="text-sm font-semibold">{{ label }}</span>
  </button>
</template>
