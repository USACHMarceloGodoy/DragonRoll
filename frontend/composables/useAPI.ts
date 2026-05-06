import type { UseFetchOptions } from 'nuxt/app'
import type { FetchError } from 'ofetch'
import type { ApiResponse, ApiError } from '~/models/apiResponse'

export function useAPI<T>(url: string | (() => string), options?: UseFetchOptions<ApiResponse<T>>) {
  return useFetch<ApiResponse<T>, FetchError<ApiError>>(url, {
    ...options,
    $fetch: useNuxtApp().$api as typeof $fetch,
  })
}
