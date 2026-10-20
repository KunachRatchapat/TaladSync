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

export interface ShopProfileForm {
  shopName: string
  zone: string
  productType: string
  ownerName: string
  phone: string
  lineId: string
}

// --- STAFF roll-call (F2) ---

export type AttendanceStatus = 'PENDING' | 'PRESENT' | 'NOSHOW'

// 1 element ใน { items: [...] } จาก GET /api/staff/checklist (แผงที่มีคน + สถานะเช็คชื่อ)
export interface ChecklistItem {
  id: number // = daily_bookings.id ใช้ยิง /bookings/:id/attendance
  stallId: number
  status: BookingStatus
  shopName?: string
  productType?: string
  attendance: AttendanceStatus
  checkedAt?: string
  stall: Stall
}

// นับฝั่ง client จาก items (checklist = เซ็ตเดียวกับ summary) — camelCase ตาม UI
export interface AttendanceSummary {
  present: number
  noShow: number
  pending: number
  total: number
}

// --- Public directory (F4) ---

// 1 element ใน { results: [...] } จาก GET /api/directory — ไม่มี PII
export interface DirectoryEntry {
  stallCode: string
  zone: string
  shopName: string
  productType: string
  posX: number
  posY: number
}

// สถานะ async รวมของหน้าจอ
export type ViewState = 'idle' | 'loading' | 'ready' | 'empty' | 'error'
