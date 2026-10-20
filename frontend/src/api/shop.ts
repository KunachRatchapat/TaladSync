import { apiFetch } from '@/api/client'
import type { ShopProfileForm } from '@/types'

export const PRODUCT_TYPES = ['ผัก', 'ผลไม้', 'ของกินเล่น', 'อาหารสด', 'ของหวาน'] as const

export async function fetchZones(): Promise<string[]> {
  const data = await apiFetch<{ zones: string[] }>('/api/zones', { auth: true })
  return data.zones
}

export async function saveShopProfile(form: ShopProfileForm): Promise<void> {
  await apiFetch<unknown>('/api/me/shop', {
    method: 'PUT',
    auth: true,
    body: JSON.stringify(form),
  })
}
