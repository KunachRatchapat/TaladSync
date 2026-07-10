import { apiFetch } from '@/api/client'
import type { AvailableStall } from '@/types'

interface MyStallsResponse {
  stalls: AvailableStall[]
}

// แผงที่ vendor เป็นเจ้าของ + สถานะวันนี้ (EXPECTED/AVAILABLE/BOOKED)
export async function fetchMyStalls(date: string): Promise<AvailableStall[]> {
  const q = date ? `?date=${encodeURIComponent(date)}` : ''
  const data = await apiFetch<MyStallsResponse>(`/api/vendor/mine${q}`, { auth: true })
  return data.stalls.filter((s) => s.stall != null)
}

// แจ้งหยุด — ใช้ stallId (stalls.id) ไม่ใช่ daily_bookings.id
export async function leaveStall(stallId: number, date: string): Promise<void> {
  await apiFetch<{ status: string }>('/api/vendor/leave', {
    method: 'POST',
    auth: true,
    body: JSON.stringify({ stallId, date: date || undefined }),
  })
}
