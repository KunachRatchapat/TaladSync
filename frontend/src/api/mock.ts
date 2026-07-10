import { ApiError } from '@/api/errors'
import type { AvailableStall, BookingForm } from '@/types'

// mock สำหรับ UI-first — สลับเป็น fetch จริงภายหลังได้ที่ composable โดยไม่แตะ UI

export type MockScenario = 'normal' | 'empty' | 'error'

const poolStalls: AvailableStall[] = Array.from({ length: 8 }, (_, i) => {
  const n = i + 1
  return {
    id: 100 + n,
    stallId: 12 + n,
    marketDate: '',
    status: 'AVAILABLE',
    stall: {
      id: 12 + n,
      code: `P-${String(n).padStart(2, '0')}`,
      zone: 'Pool',
      posX: n,
      posY: 2,
    },
  }
})

const FAKE_DELAY_MS = 800

export function fetchAvailableStallsMock(
  date: string,
  scenario: MockScenario = 'normal',
): Promise<AvailableStall[]> {
  return new Promise((resolve, reject) => {
    setTimeout(() => {
      if (scenario === 'error') {
        reject(new Error('mock: โหลดข้อมูลไม่สำเร็จ'))
        return
      }
      if (scenario === 'empty') {
        resolve([])
        return
      }
      resolve(poolStalls.map((s) => ({ ...s, marketDate: date })))
    }, FAKE_DELAY_MS)
  })
}

// จำลองผลการจอง — success (201) / taken, already (409)
export type BookMockResult = 'success' | 'taken' | 'already'

export function bookStallMock(
  _form: BookingForm,
  result: BookMockResult = 'success',
): Promise<void> {
  return new Promise((resolve, reject) => {
    setTimeout(() => {
      if (result === 'taken') {
        reject(new ApiError('STALL_NOT_AVAILABLE'))
        return
      }
      if (result === 'already') {
        reject(new ApiError('ALREADY_HAS_STALL'))
        return
      }
      resolve()
    }, 700)
  })
}
