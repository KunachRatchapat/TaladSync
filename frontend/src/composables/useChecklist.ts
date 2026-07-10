import { computed, ref } from 'vue'

import { fetchChecklist, markAttendance } from '@/api/staff'
import type { AttendanceStatus, AttendanceSummary, ChecklistItem, ViewState } from '@/types'

export function useChecklist() {
  const state = ref<ViewState>('idle')
  const items = ref<ChecklistItem[]>([])
  const errorMessage = ref('')
  const marking = ref(new Set<number>()) // id ที่กำลังยิง API อยู่ (ไว้ disable ปุ่ม)

  // นับสรุปจาก items ตรงๆ → อัปเดตทันทีหลังเช็ค ไม่ต้อง refetch
  const summary = computed<AttendanceSummary>(() => {
    const s: AttendanceSummary = { present: 0, noShow: 0, pending: 0, total: items.value.length }
    for (const it of items.value) {
      if (it.attendance === 'PRESENT') s.present++
      else if (it.attendance === 'NOSHOW') s.noShow++
      else s.pending++
    }
    return s
  })

  async function load(date: string) {
    state.value = 'loading'
    errorMessage.value = ''
    try {
      const data = await fetchChecklist(date)
      items.value = data
      state.value = data.length > 0 ? 'ready' : 'empty'
    } catch {
      state.value = 'error'
      errorMessage.value = 'โหลดรายการเช็คชื่อไม่สำเร็จ ลองใหม่อีกครั้ง'
    }
  }

  // optimistic: เปลี่ยน UI ทันที แล้วค่อย sync กับ response — error ค่อย rollback
  async function mark(id: number, present: boolean) {
    const item = items.value.find((it) => it.id === id)
    if (!item) return

    const target: AttendanceStatus = present ? 'PRESENT' : 'NOSHOW'
    if (item.attendance === target || marking.value.has(id)) return

    const prev = item.attendance
    item.attendance = target
    marking.value.add(id)
    errorMessage.value = ''
    try {
      const updated = await markAttendance(id, present)
      item.attendance = updated.attendance
      item.checkedAt = updated.checkedAt
    } catch {
      item.attendance = prev
      errorMessage.value = 'บันทึกไม่สำเร็จ ลองใหม่อีกครั้ง'
    } finally {
      marking.value.delete(id)
    }
  }

  return { state, items, errorMessage, marking, summary, load, mark }
}
