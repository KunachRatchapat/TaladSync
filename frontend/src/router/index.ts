import { createRouter, createWebHistory } from 'vue-router'

import StaffView from '@/views/StaffView.vue'
import WalkinView from '@/views/WalkinView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'walkin',
      component: WalkinView,
    },
    {
      path: '/staff',
      name: 'staff',
      component: StaffView,
    },
  ],
})

export default router
