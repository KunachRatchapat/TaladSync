<script setup lang="ts">
import SvgIcon from '@jamescoyle/vue-icon'
import { mdiClose, mdiMagnify } from '@mdi/js'

const model = defineModel<string>({ required: true })
</script>

<template>
  <div class="relative">
    <span
      class="pointer-events-none absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400"
      aria-hidden="true"
    >
      <SvgIcon type="mdi" :path="mdiMagnify" :size="20" />
    </span>
    <input
      v-model="model"
      type="text"
      inputmode="search"
      placeholder="ค้นหาร้าน / ประเภทสินค้า"
      aria-label="ค้นหาร้าน"
      class="w-full rounded-xl border border-gray-300 bg-white py-3 pl-11 pr-10 text-gray-900 focus:border-brand-600 focus:outline-none focus:ring-1 focus:ring-brand-600"
    />
    <button
      v-if="model"
      type="button"
      aria-label="ล้างคำค้นหา"
      class="absolute right-2.5 top-1/2 -translate-y-1/2 rounded-full p-1 text-gray-400 active:text-gray-600"
      @click="model = ''"
    >
      <SvgIcon type="mdi" :path="mdiClose" :size="18" />
    </button>
  </div>
</template>
