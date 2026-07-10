import { apiFetch } from '@/api/client'
import type { ChecklistItem } from '@/types'

interface ChecklistResponse {
  items: ChecklistItem[]
}

export async function fetchChecklist(date: string): Promise<ChecklistItem[]> {
  const q = date ? `?date=${encodeURIComponent(date)}` : ''
  const data = await apiFetch<ChecklistResponse>(`/api/staff/checklist${q}`, { auth: true })
  return data.items.filter((it) => it.stall != null)
}

// :id = daily_bookings.id (ไม่ใช่ stallId) — backend idempotent, สลับ มา/ไม่มา ได้
export async function markAttendance(
  bookingId: number,
  present: boolean,
): Promise<ChecklistItem> {
  return apiFetch<ChecklistItem>(`/api/bookings/${bookingId}/attendance`, {
    method: 'POST',
    auth: true,
    body: JSON.stringify({ present }),
  })
}
