export type BookErrorCode = 'STALL_NOT_AVAILABLE' | 'ALREADY_HAS_STALL' | 'UNKNOWN'


export class ApiError extends Error {
  constructor(public code: BookErrorCode) {
    super(code)
    this.name = 'ApiError'
  }
}
