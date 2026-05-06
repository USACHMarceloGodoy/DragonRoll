export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()

  deleteCookie(event, config.cookieName, {
    httpOnly: true,
    path: '/',
    sameSite: 'lax',
    secure: process.env.NODE_ENV === 'production',
  })

  const data = await $fetch(config.public.authUrl + '/logout', {
    method: 'POST',
  })
  if (data.error) {
    return createError({ statusCode: 500, statusMessage: 'Logout failed' })
  }

  return null
})
