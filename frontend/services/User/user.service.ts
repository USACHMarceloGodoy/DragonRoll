import type { User } from '~/models/User/user.model'

export class UserService {
  /**
   * Fetches the current user from the API.
   * @returns {Promise<User | null>} The current user or null if not authenticated.
   */
  async currentUser(): Promise<User | null> {
    try {
      const { data, error } = await useAPI<User>('/auth/me', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      })

      if (error.value) {
        throw new Error(error.value.message || 'Error al obtener el usuario actual')
      }
      return data.value?.data ?? null
    }
    catch (error) {
      throw new Error(error.message || 'Error al procesar la solicitud')
    }
  }

  async currentUserError(): Promise<User | null> {
    try {
      const { data, error } = await useAPI<User>('/auth/meError', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      })

      if (error.value) {
        throw new Error(error.value.message || 'Error al obtener el usuario actual')
      }
      return data.value?.data ?? null
    }
    catch (error) {
      throw new Error(error.message || 'Error al procesar la solicitud')
    }
  }
}
