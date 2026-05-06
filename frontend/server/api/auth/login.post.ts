import type { LoginCredentials, LoginResponse } from '~/models/Auth/login.model'
import { LoginResponseToSession } from '~/models/Auth/login.model'
import { sign, serialize } from '~/utils/cookie'

export default defineEventHandler(async (event) => {
  const body = await readBody<LoginCredentials>(event)

  const { user, password } = body
  if (!user || !password) {
    throw createError({
      statusCode: 400,
      message: 'Faltan el nombre de usuario o la contraseña',
    })
  }
  const config = useRuntimeConfig()

  try {
    const data: LoginResponse = await $fetch(config.public.authUrl + '/login', {
      method: 'POST',
      body: JSON.stringify(body),
    })

    const sessionData = await LoginResponseToSession(data)

    const session = serialize(sessionData)
    const signedSession = sign(session, config.sessionPassword)

    setCookie(event, config.cookieName, signedSession, {
      httpOnly: true,
      path: '/',
      sameSite: 'strict',
      secure: process.env.NODE_ENV === 'production',
      expires: new Date(Date.now() + config.cookieExpires*1000), // cookieExpires está en segundos
    })

    return sessionData
  }
  catch (error: any) {
    throw createError({
      statusCode: error.response?.status || error.statusCode || 500,
      message: error.data?.message || error.message || 'Error desconocido en el servicio de autenticación',
    })
  }
})
