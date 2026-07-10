import { apiFetch } from '@/api/client'
import type { DirectoryEntry } from '@/types'

interface DirectoryResponse {
  results: DirectoryEntry[]
}

// public — ไม่ต้อง auth · q ว่าง = ทุกร้านของวันนั้น
export async function searchDirectory(q: string, date: string): Promise<DirectoryEntry[]> {
  const params = new URLSearchParams()
  if (q) params.set('q', q)
  if (date) params.set('date', date)
  const qs = params.toString()
  const data = await apiFetch<DirectoryResponse>(`/api/directory${qs ? `?${qs}` : ''}`)
  return data.results
}
