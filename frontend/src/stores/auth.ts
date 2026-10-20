import { defineStore } from 'pinia'
import { ref } from 'vue'

import { fetchMe, mockLogin } from '@/api/auth'
import { clearStoredAuth, getStoredRole, getStoredToken, setStoredAuth } from '@/api/client'
import type { Role, User } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(getStoredToken())
  const role = ref<Role | null>(getStoredRole())
  const user = ref<User | null>(null)
  const loading = ref(false)
  const error = ref('')

  
  async function ensureRole(target: Role): Promise<void> {
    if (token.value && role.value === target) {
      if (user.value) return
      try {
        user.value = await fetchMe()
        return
      } catch {
        logout()
      }
    }

    loading.value = true
    error.value = ''
    try {
      const res = await mockLogin(target)
      token.value = res.token
      role.value = res.user.role
      user.value = res.user
      setStoredAuth(res.token, res.user.role)
    } catch {
      error.value = 'เข้าสู่ระบบไม่สำเร็จ ลองใหม่อีกครั้ง'
      throw new Error(error.value)
    } finally {
      loading.value = false
    }
  }

  function logout() {
    token.value = null
    role.value = null
    user.value = null
    clearStoredAuth()
  }

  return { token, role, user, loading, error, ensureRole, logout }
})
