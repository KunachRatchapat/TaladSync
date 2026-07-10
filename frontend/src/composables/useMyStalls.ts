import { ref } from 'vue'

import { ApiError } from '@/api/errors'
import { fetchMyStalls, leaveStall } from '@/api/vendor'
import type { AvailableStall, ViewState } from '@/types'
import { leaveErrorMessageTH } from '@/utils/errors'

export function useMyStalls() {
  const state = ref<ViewState>('idle')
  const stalls = ref<AvailableStall[]>([])
  const errorMessage = ref('')
  const leaving = ref(new Set<number>()) // stallId ที่กำลังยิงแจ้งหยุด (ไว้ disable ปุ่ม)

  async function load(date: string) {
    state.value = 'loading'
    errorMessage.value = ''
    try {
      const data = await fetchMyStalls(date)
      stalls.value = data
      state.value = data.length > 0 ? 'ready' : 'empty'
    } catch {
      state.value = 'error'
      errorMessage.value = 'โหลดแผงของคุณไม่สำเร็จ ลองใหม่อีกครั้ง'
    }
  }

  // optimistic: EXPECTED → AVAILABLE ทันที, error ค่อย rollback
  async function leave(stallId: number, date: string) {
    const item = stalls.value.find((s) => s.stallId === stallId)
    if (!item || item.status !== 'EXPECTED' || leaving.value.has(stallId)) return

    const prev = item.status
    item.status = 'AVAILABLE'
    leaving.value.add(stallId)
    errorMessage.value = ''
    try {
      await leaveStall(stallId, date)
    } catch (err) {
      item.status = prev
      errorMessage.value =
        err instanceof ApiError ? leaveErrorMessageTH(err.code) : 'แจ้งหยุดไม่สำเร็จ ลองใหม่อีกครั้ง'
    } finally {
      leaving.value.delete(stallId)
    }
  }

  return { state, stalls, errorMessage, leaving, load, leave }
}
