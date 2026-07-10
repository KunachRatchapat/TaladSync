<script setup lang="ts">
// เมนูสลับหน้า (dev only) — gate ด้วย import.meta.env.DEV ที่ App.vue
const links = [
  { to: '/', label: 'ขาจร' },
  { to: '/staff', label: 'เจ้าหน้าที่' },
  { to: '/regular', label: 'แม่ค้าประจำ' },
  { to: '/directory', label: 'ค้นหาร้าน' },
]
</script>

<template>
  <nav
    class="fixed bottom-3 left-1/2 z-30 -translate-x-1/2 rounded-full border border-gray-200 bg-white/95 px-1.5 py-1.5 shadow-lg backdrop-blur"
    aria-label="สลับหน้า (dev)"
  >
    <ul class="flex items-center gap-1">
      <li v-for="l in links" :key="l.to">
        <RouterLink
          :to="l.to"
          class="block rounded-full px-3 py-1.5 text-xs font-medium text-gray-500 transition active:bg-gray-100"
          exact-active-class="bg-gray-900 text-white"
        >
          {{ l.label }}
        </RouterLink>
      </li>
    </ul>
  </nav>
</template>
