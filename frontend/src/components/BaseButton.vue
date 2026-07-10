<script setup lang="ts">
withDefaults(
  defineProps<{
    variant?: 'primary' | 'secondary' | 'danger'
    type?: 'button' | 'submit'
    disabled?: boolean
    loading?: boolean
    block?: boolean
  }>(),
  {
    variant: 'primary',
    type: 'button',
    disabled: false,
    loading: false,
    block: true,
  },
)

const emit = defineEmits<{ (e: 'click', ev: MouseEvent): void }>()

const variantClass: Record<string, string> = {
  primary: 'bg-brand-600 text-white active:bg-brand-700',
  secondary: 'border border-gray-300 bg-white text-gray-700 active:bg-gray-50',
  danger: 'bg-danger-600 text-white active:bg-danger-700',
}
</script>

<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    :class="[
      'inline-flex items-center justify-center gap-2 rounded-xl py-3 px-4 font-medium transition',
      'disabled:cursor-not-allowed disabled:opacity-50',
      variantClass[variant],
      block ? 'w-full' : '',
    ]"
    @click="emit('click', $event)"
  >
    <svg
      v-if="loading"
      class="h-5 w-5 animate-spin"
      viewBox="0 0 24 24"
      fill="none"
      aria-hidden="true"
    >
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.4 0 0 5.4 0 12h4z" />
    </svg>
    <slot />
  </button>
</template>
