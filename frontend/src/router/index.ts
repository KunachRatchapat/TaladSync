import { createRouter, createWebHistory } from 'vue-router'

import WalkinView from '@/views/WalkinView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'walkin',
      component: WalkinView,
    },
  ],
})

export default router
