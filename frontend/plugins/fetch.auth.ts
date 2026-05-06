export default defineNuxtPlugin((nuxtApp) => {
  const config = useRuntimeConfig()

  const authFetch = async (endpoint: string, options: RequestInit = {}) => {
    const { token } = useAuth()
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(options.headers as Record<string, string> || {}),
    }

    if (token) {
      // Puedes enviar el token como header o como cookie, aquí ambos ejemplos:
      headers['Authorization'] = `Bearer ${token}`
      headers['Cookie'] = `token=${token}`
    }

    const response = await fetch(`${config.public.authUrl}${endpoint}`, {
      ...options,
      headers,
    })

    if (!response.ok) {
      const errorBody = await response.json()
      const errorMessage = errorBody?.message || 'Error desconocido'
      throw new Error(`${errorMessage}`)
    }

    return await response.json()
  }

  return {
    provide: {
      authFetch,
    },
  }
})
