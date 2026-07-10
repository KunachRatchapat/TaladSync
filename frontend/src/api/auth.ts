import { apiFetch } from '@/api/client'
import type { User } from '@/types'

interface MockLoginResponse {
  token: string
  user: User
}

export async function mockLogin(role: 'WALKIN'): Promise<MockLoginResponse> {
  return apiFetch<MockLoginResponse>('/api/auth/mock-login', {
    method: 'POST',
    body: JSON.stringify({ role }),
  })
}
