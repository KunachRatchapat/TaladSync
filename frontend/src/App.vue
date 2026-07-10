<script setup lang="ts">
import DevNav from './components/DevNav.vue'

const isDev = import.meta.env.DEV
</script>

<template>
  <RouterView />
  <DevNav v-if="isDev" />
</template>
