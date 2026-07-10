import { ref } from 'vue'

import { bookStall } from '@/api/booking'
import { ApiError, type BookErrorCode } from '@/api/errors'
import type { BookingForm } from '@/types'
import { bookErrorMessageTH } from '@/utils/errors'

export function useBooking() {
  const submitting = ref(false)
  const success = ref(false)
  const errorMessage = ref('')
  const errorCode = ref<BookErrorCode | null>(null)

  async function submit(form: BookingForm) {
    submitting.value = true
    success.value = false
    errorMessage.value = ''
    errorCode.value = null
    try {
      await bookStall(form)
      success.value = true
    } catch (err) {
      if (err instanceof ApiError) {
        errorCode.value = err.code
        errorMessage.value = bookErrorMessageTH(err.code)
      } else {
        errorMessage.value = 'จองไม่สำเร็จ ลองใหม่อีกครั้ง'
      }
    } finally {
      submitting.value = false
    }
  }

  function reset() {
    submitting.value = false
    success.value = false
    errorMessage.value = ''
    errorCode.value = null
  }

  return { submitting, success, errorMessage, errorCode, submit, reset }
}
