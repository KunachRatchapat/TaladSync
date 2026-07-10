import { apiFetch } from '@/api/client'
import type { AvailableStall, BookingForm } from '@/types'

interface AvailableStallsResponse {
  stalls: AvailableStall[]
}

interface MyBookingResponse {
  booking: AvailableStall | null
}

export async function fetchMyBookingToday(date: string): Promise<AvailableStall | null> {
  try {
    const q = date ? `?date=${encodeURIComponent(date)}` : ''
    const data = await apiFetch<MyBookingResponse>(`/api/bookings/mine${q}`, { auth: true })
    if (!data.booking?.stall) return null
    return data.booking
  } catch {
    // API เก่ายังไม่มี endpoint นี้ (404) หรือ error อื่น → ข้ามไปโหลดแผงว่าง
    return null
  }
}

export async function fetchAvailableStalls(date: string): Promise<AvailableStall[]> {
  const q = date ? `?date=${encodeURIComponent(date)}` : ''
  const data = await apiFetch<AvailableStallsResponse>(`/api/stalls/available${q}`, { auth: true })

  return data.stalls.filter((s) => s.stall != null)
}

export async function bookStall(form: BookingForm): Promise<void> {
  await apiFetch<{ status: string }>('/api/bookings', {
    method: 'POST',
    auth: true,
    body: JSON.stringify({
      stallId: form.stallId,
      shopName: form.shopName,
      productType: form.productType,
      date: form.date || undefined,
    }),
  })
}
