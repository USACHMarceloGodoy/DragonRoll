import type { Session } from '~/models/Auth/login.model'

export const useAuthUser = () => {
  return useState<Session | null>('session', () => null)
}
