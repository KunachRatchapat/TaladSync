import { defineStore } from 'pinia'
import { ref } from 'vue'

import { mockLogin } from '@/api/auth'
import { clearStoredToken, getStoredToken, setStoredToken } from '@/api/client'
import type { User } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(getStoredToken())
  const user = ref<User | null>(null)
  const loading = ref(false)
  const error = ref('')

  async function ensureWalkin(): Promise<void> {
    if (token.value) return

    loading.value = true
    error.value = ''
    try {
      const res = await mockLogin('WALKIN')
      token.value = res.token
      user.value = res.user
      setStoredToken(res.token)
    } catch {
      error.value = 'เข้าสู่ระบบไม่สำเร็จ ลองใหม่อีกครั้ง'
      throw new Error(error.value)
    } finally {
      loading.value = false
    }
  }

  function logout() {
    token.value = null
    user.value = null
    clearStoredToken()
  }

  return { token, user, loading, error, ensureWalkin, logout }
})
