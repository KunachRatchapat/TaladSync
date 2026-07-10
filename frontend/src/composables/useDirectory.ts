import { ref } from 'vue'

import { searchDirectory } from '@/api/directory'
import type { DirectoryEntry, ViewState } from '@/types'

export function useDirectory() {
  const state = ref<ViewState>('idle')
  const results = ref<DirectoryEntry[]>([])
  const errorMessage = ref('')

  // reqId กัน response เก่ามาทับใหม่ (พิมพ์เร็วๆ)
  let reqId = 0

  async function search(q: string, date: string) {
    const current = ++reqId
    state.value = 'loading'
    errorMessage.value = ''
    try {
      const data = await searchDirectory(q, date)
      if (current !== reqId) return // มี request ใหม่กว่าแล้ว ทิ้งอันนี้
      results.value = data
      state.value = data.length > 0 ? 'ready' : 'empty'
    } catch {
      if (current !== reqId) return
      state.value = 'error'
      errorMessage.value = 'ค้นหาร้านไม่สำเร็จ ลองใหม่อีกครั้ง'
    }
  }

  return { state, results, errorMessage, search }
}
