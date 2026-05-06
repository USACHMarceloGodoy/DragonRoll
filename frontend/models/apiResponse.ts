export type ApiResponse<T> = {
  data?: T | null
  message: string
  status?: string
}

export type ApiError = {
  message: string
}
