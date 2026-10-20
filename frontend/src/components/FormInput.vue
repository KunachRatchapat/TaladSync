<script setup lang="ts">
const model = defineModel<string>({ required: true })

withDefaults(
  defineProps<{
    label: string
    placeholder?: string
    inputmode?: 'text' | 'numeric' | 'tel'
    maxlength?: number
    error?: string
  }>(),
  { placeholder: '', inputmode: 'text', maxlength: 100, error: '' },
)
</script>

<template>
  <label class="flex flex-col gap-2">
    <span class="text-sm font-semibold text-slate-600">{{ label }}</span>
    <input
      v-model="model"
      type="text"
      :inputmode="inputmode"
      :maxlength="maxlength"
      :placeholder="placeholder"
      class="h-11 w-full rounded-xl border bg-white px-4 text-base text-slate-800 placeholder:text-slate-400 focus:outline-none focus:ring-1"
      :class="
        error
          ? 'border-danger-600 focus:ring-danger-600'
          : 'border-slate-200 focus:border-emerald-500 focus:ring-emerald-500'
      "
    />
    <span v-if="error" class="text-sm text-danger-600">{{ error }}</span>
  </label>
</template>
