import type { User } from '../User/user.model'

export interface LoginCredentials {
  user: string
  password: string
}

export interface LoginResponse {
  expire: string
  token: string
  user: User

}

export interface LoginResponseError {
  status: string | undefined
  message: string | undefined
}

export interface Session {
  username: string
  userId: string
  token: string
  role: string
  expire: string
}

export async function LoginResponseToSession(data: LoginResponse): Promise<Session> {
  const session: Session = {
    username: data.user.username,
    userId: data.user._id,
    token: data.token,
    role: data.user.role,
    expire: data.expire,
  }
  return session
}
