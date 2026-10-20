import { createRouter, createWebHistory } from 'vue-router'

import DirectoryView from '@/views/DirectoryView.vue'
import RegularView from '@/views/RegularView.vue'
import StaffView from '@/views/StaffView.vue'
import WalkinHomeView from '@/views/WalkinHomeView.vue'
import WalkinView from '@/views/WalkinView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'walkin',
      component: WalkinHomeView,
    },
    {
      path: '/walkin/book',
      name: 'walkin-book',
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
    {
      path: '/directory',
      name: 'directory',
      component: DirectoryView,
    },
  ],
})

export default router
