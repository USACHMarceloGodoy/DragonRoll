import { unsign, deserialize } from '~/utils/cookie'
import type { Session } from '~/models/Auth/login.model'

export default defineEventHandler(async (event) => {
  const session = await getSessionFromContext(event)
  if (session) event.context.session = session
})

declare module 'h3' {
  interface H3EventContext {
    session?: Session
  }
}

export async function getSessionFromContext(event: H3Event) {
  const config = useRuntimeConfig()

  const cookie = getCookie(event, config.cookieName)
  if (!cookie) return null

  const unsignedSession = unsign(cookie, config.sessionPassword)
  if (!unsignedSession) return null

  const session = deserialize(unsignedSession)

  return session
}
