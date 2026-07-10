import type { BookErrorCode } from '@/api/errors'


export function bookErrorMessageTH(code: BookErrorCode): string {
  switch (code) {
    case 'STALL_NOT_AVAILABLE':
      return 'แผงนี้เพิ่งถูกจองไปแล้ว ลองเลือกแผงอื่น'
    case 'ALREADY_HAS_STALL':
      return 'คุณจองครบแล้ววันนี้ (จำกัด 1 แผงต่อวัน) ไม่สามารถจองแผงเพิ่มได้'
    default:
      return 'จองไม่สำเร็จ ลองใหม่อีกครั้ง'
  }
}
