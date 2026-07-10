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

export function leaveErrorMessageTH(code: BookErrorCode): string {
  switch (code) {
    case 'LEAVE_NOT_ALLOWED':
      return 'แผงนี้แจ้งหยุดไม่ได้แล้ว (อาจแจ้งไปแล้วหรือมีขาจรจองต่อ)'
    case 'NOT_OWN_STALL':
      return 'ไม่ใช่แผงของคุณ'
    default:
      return 'แจ้งหยุดไม่สำเร็จ ลองใหม่อีกครั้ง'
  }
}
