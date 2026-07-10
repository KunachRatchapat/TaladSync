export type BookErrorCode =
  | 'STALL_NOT_AVAILABLE'
  | 'ALREADY_HAS_STALL'
  | 'NOT_OWN_STALL'
  | 'LEAVE_NOT_ALLOWED'
  | 'UNKNOWN'


export class ApiError extends Error {
  constructor(public code: BookErrorCode) {
    super(code)
    this.name = 'ApiError'
  }
}
