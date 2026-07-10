import { ApiError } from '@/api/errors'

const TOKEN_KEY = 'taladsync_token'

export function getStoredToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function setStoredToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token)
}

export function clearStoredToken(): void {
  localStorage.removeItem(TOKEN_KEY)
}

function apiBase(): string {
  return import.meta.env.VITE_API_URL ?? 'http://localhost:8080'
}

interface ApiErrorBody {
  error?: string
}

function mapHttpError(status: number, message: string): ApiError {
  if (status === 409) {
    if (message === 'stall is not available') return new ApiError('STALL_NOT_AVAILABLE')
    if (message === 'vendor already booked today') return new ApiError('ALREADY_HAS_STALL')
  }
  return new ApiError('UNKNOWN')
}

export async function apiFetch<T>(
  path: string,
  options: RequestInit & { auth?: boolean } = {},
): Promise<T> {
  const headers = new Headers(options.headers)
  if (options.body) headers.set('Content-Type', 'application/json')

  if (options.auth) {
    const token = getStoredToken()
    if (token) headers.set('Authorization', `Bearer ${token}`)
  }

  const res = await fetch(`${apiBase()}${path}`, { ...options, headers })

  if (!res.ok) {
    const body = (await res.json().catch(() => ({}))) as ApiErrorBody
    throw mapHttpError(res.status, body.error ?? '')
  }

  if (res.status === 204) return undefined as T
  return res.json() as Promise<T>
}
