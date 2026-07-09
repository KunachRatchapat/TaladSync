<script setup lang="ts">
const apiUrl = import.meta.env.VITE_API_URL ?? '(ยังไม่ได้ตั้ง VITE_API_URL)'
</script>

<template>
  <main class="mx-auto flex min-h-screen max-w-app flex-col items-center justify-center gap-4 px-4">
    <span class="rounded-full bg-brand-50 px-4 py-1 text-sm font-medium text-brand-700">
      TaladSync Frontend
    </span>

    <h1 class="text-xl font-semibold text-gray-900">ตั้งค่าโปรเจคสำเร็จ</h1>

    <p class="text-center text-sm text-gray-500">
      Vue 3 + Vite + TypeScript + Tailwind + Pinia + Vue Router พร้อมใช้งานแล้ว
    </p>

    <button
      class="w-full max-w-xs rounded-xl bg-brand-600 py-3 font-medium text-white transition active:bg-brand-700"
    >
      ปุ่มทดสอบ design token (สีเขียว brand-600)
    </button>

    <p class="text-xs text-gray-400">API: {{ apiUrl }}</p>
  </main>
</template>
