export default defineNuxtPlugin((nuxtApp) => {
  const { token } = useAuth()
  const baseURL = useRuntimeConfig().public.backendUrl.replace(/\/$/, '')
  const api = $fetch.create({
    baseURL: baseURL,
    onRequest({ options }) {
      if (token.value) {
        options.headers.set('Authorization', `Bearer ${token.value}`)
        options.headers.set('Cookie', `token=${token.value}`)
      }
    },
    onResponse({ response }) {
      const res = response._data
      if (response.status >= 400 && response.status < 600) {
        throw new Error(res.message || 'API returned an unknown error')
      }
      return res
    },
    onResponseError({ response }) {
      throw new Error(response._data.message || 'API returned an unknown error')
    },
  })

  return {
    provide: {
      api,
    },
  }
})
