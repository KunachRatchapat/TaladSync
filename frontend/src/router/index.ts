import { createRouter, createWebHistory } from 'vue-router'

import RegularView from '@/views/RegularView.vue'
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
    {
      path: '/regular',
      name: 'regular',
      component: RegularView,
    },
  ],
})

export default router
