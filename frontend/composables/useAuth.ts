import { useAuthUser } from './useAuthUser'
import type { Session } from '~/models/Auth/login.model'

export const useAuth = () => {
  const authUser = useAuthUser()

  const setUser = (session: Session | null) => {
    authUser.value = session
  }

  const login = async (LoginCredentials) => {
    const data = await $fetch<Session>('/api/auth/login', {
      method: 'POST',
      body: LoginCredentials,
    })

    setUser(data)

    return authUser
  }

  const logout = async () => {
    const data = await $fetch<Session>('/api/auth/logout', {
      method: 'POST',
    })

    setUser(null)
    useRouter().push('/login')
  }

  const status = computed(() => {
    if (!authUser.value) {
      return 'unauthenticated'
    }
    return 'authenticated'
  })

  const me = async () => {
    if (!authUser.value) {
      try {
        const data = await $fetch<Session>('/api/auth/me', {
          headers: useRequestHeaders(['cookie']) as HeadersInit,
        })
        if (!data) {
          setUser(null)
        }
        setUser(data)
      }
      catch {
        setUser(null)
      }
    }

    return authUser
  }

  const data = computed(() => authUser.value)
  const token = computed(() => authUser.value?.token)

  return {
    login,
    logout,
    me,
    data,
    token,
    status,
  }
}
