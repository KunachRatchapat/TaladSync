import { ref } from 'vue'

import { fetchAvailableStalls } from '@/api/booking'
import type { AvailableStall, ViewState } from '@/types'

export function useAvailableStalls() {
  const state = ref<ViewState>('idle')
  const stalls = ref<AvailableStall[]>([])
  const errorMessage = ref('')

  async function load(date: string) {
    state.value = 'loading'
    errorMessage.value = ''
    try {
      const data = await fetchAvailableStalls(date)
      stalls.value = data
      state.value = data.length > 0 ? 'ready' : 'empty'
    } catch {
      state.value = 'error'
      errorMessage.value = 'โหลดรายการแผงไม่สำเร็จ ลองใหม่อีกครั้ง'
    }
  }

  return { state, stalls, errorMessage, load }
}
