// UI types — อิง shape จาก backend (domain/models.go) เพื่อสลับไป fetch จริงได้ทีเดียว

export type Role = 'REGULAR' | 'WALKIN' | 'STAFF'

export interface User {
  id: number
  name: string
  phone: string
  role: Role
}

export type BookingStatus = 'EXPECTED' | 'AVAILABLE' | 'BOOKED' | 'ABSENT'

export interface Stall {
  id: number
  code: string
  zone: string
  posX: number
  posY: number
}

// ตรงกับ 1 element ใน { stalls: [...] } จาก GET /api/stalls/available
export interface AvailableStall {
  id: number
  stallId: number
  marketDate: string
  status: BookingStatus
  shopName?: string
  productType?: string
  stall: Stall
}

export interface BookingForm {
  stallId: number
  shopName: string
  productType: string
  date: string
}

// สถานะ async รวมของหน้าจอ
export type ViewState = 'idle' | 'loading' | 'ready' | 'empty' | 'error'
