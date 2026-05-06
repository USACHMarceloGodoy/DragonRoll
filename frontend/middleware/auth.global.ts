export default defineNuxtRouteMiddleware(async (to, from) => {
  const { me, status } = useAuth()
  const session = await me()
  if (!session) {
    return navigateTo('/login')
  }
  const isAuthenticated = status.value === 'authenticated'
  const toPageIsGuestOnly = to.meta?.middleware?.includes('guest-only') === true
  if (isAuthenticated && toPageIsGuestOnly) {
    return navigateTo({ name: 'index' })
  }
  const toPageRequiresAuth = to.meta?.auth === true
  if (!isAuthenticated && toPageRequiresAuth) {
    return navigateTo('/login')
  }
})
