import { apiFetch } from '@/api/client'
import type { Role, User } from '@/types'

interface MockLoginResponse {
  token: string
  user: User
}

export async function fetchMe(): Promise<User> {
  const data = await apiFetch<{ user: User }>('/api/auth/me', { auth: true })
  return data.user
}

export async function mockLogin(role: Role): Promise<MockLoginResponse> {
  return apiFetch<MockLoginResponse>('/api/auth/mock-login', {
    method: 'POST',
    body: JSON.stringify({ role }),
  })
}
