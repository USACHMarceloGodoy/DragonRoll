import type { LoginResponse, LoginCredentials, LoginResponseError } from '~/models/Auth/login.model'

// Define a type for the login response, adjust fields as needed

export class AuthService {
  async login(loginCredentials: LoginCredentials): Promise<LoginResponse> {
    try {
      const { $authFetch } = useNuxtApp()
      const data: LoginResponse | LoginResponseError = await $authFetch('/login', {
        method: 'POST',
        body: JSON.stringify(loginCredentials),
      })
      if (!('token' in data) || data.error || !data.user || !data.user._id) {
        if (data.message == 'missing Username or Password') {
          throw new Error('Faltan el nombre de usuario o la contraseña')
        }
        throw new Error(data.message || 'Error desconocido en el login')
      }
      return data
    }
    catch (error: any) {
      const errorMessage = error.data?.message || error.message || 'Error desconocido en el login'
      throw new Error(errorMessage)
    }
  }

  async logout(): Promise<void> {
    try {
      const { $authFetch } = useNuxtApp()
      await $authFetch('/logout', {
        method: 'POST',
      })
    }
    catch (error: any) {
      const errorMessage = error.data?.message || error.message || 'Error al cerrar sesión'
      throw new Error(errorMessage)
    }
  }
}
